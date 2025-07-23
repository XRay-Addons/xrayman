package grpcconn

import (
	"context"
	"fmt"
	"sync"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// grpc connection with retry
type GRPCConn struct {
	target string

	conn *grpc.ClientConn
	mu   sync.RWMutex

	log *zap.Logger
}

var _ grpc.ClientConnInterface = (*GRPCConn)(nil)

func New(target string, log *zap.Logger) (*GRPCConn, error) {

	if log == nil {
		return nil, fmt.Errorf("%w: grpconn connect: logger", errdefs.ErrNilArgPassed)
	}

	return &GRPCConn{
		target: target,
		log:    log,
	}, nil
}

func (c *GRPCConn) Connect(ctx context.Context) error {
	if c == nil {
		return fmt.Errorf("%w: grpconn", errdefs.ErrNilObjectCall)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil
	}

	conn, err := grpc.NewClient(c.target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("%w: connect: %v", errdefs.ErrGRPC, err)
	}

	// start connecting
	state := conn.GetState()
	if state == connectivity.Idle {
		conn.Connect()
	}

	// wait till connected or cancelled
	for {
		if state == connectivity.Ready {
			c.log.Info("grpc connected")
			c.conn = conn
			return nil
		}
		c.log.Warn("grpc connecting", zap.String("state", state.String()))
		if !conn.WaitForStateChange(ctx, state) {
			if err := conn.Close(); err != nil {
				c.log.Warn("grcp connection close", zap.Error(err))
			}
			return fmt.Errorf("grpc connect failed: %w", ctx.Err())
		}
		state = conn.GetState()
	}
}

func (c *GRPCConn) Disconnect(ctx context.Context) error {
	if c == nil {
		return fmt.Errorf("%w: grpconn", errdefs.ErrNilObjectCall)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return nil
	}
	defer func() { c.conn = nil }()

	state := c.conn.GetState()
	if state == connectivity.Shutdown {
		return nil
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("%w: disconnect: %w", errdefs.ErrGRPC, err)
	}
	return nil
}

func (c *GRPCConn) Close(ctx context.Context) error {
	if c == nil {
		return nil
	}
	return c.Disconnect(ctx)
}

func (c *GRPCConn) Invoke(
	ctx context.Context,
	method string,
	args any,
	reply any,
	opts ...grpc.CallOption,
) error {
	if c == nil || c.conn == nil {
		return fmt.Errorf("%w: grpc connection", errdefs.ErrNilObjectCall)
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.conn.Invoke(ctx, method, args, reply, opts...)
}

func (c *GRPCConn) NewStream(
	ctx context.Context,
	desc *grpc.StreamDesc,
	method string,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	if c == nil || c.conn == nil {
		return nil, fmt.Errorf("%w: grpc connection", errdefs.ErrNilObjectCall)
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.conn.NewStream(ctx, desc, method, opts...)
}
