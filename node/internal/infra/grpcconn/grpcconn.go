package grpcconn

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	healthst "google.golang.org/grpc/status"
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
		return nil, fmt.Errorf("%w: init grpc client: logger", errdefs.ErrNilObjectCall)
	}

	return &GRPCConn{
		target: target,
		log:    log,
	}, nil
}

func (conn *GRPCConn) Connect(ctx context.Context) error {
	if conn == nil {
		return fmt.Errorf("%w: grpconn", errdefs.ErrNilObjectCall)
	}

	conn.mu.Lock()
	defer conn.mu.Unlock()

	// run installing service loop
	return conn.tryConnect(ctx, conn.target, conn.log)
}

func (conn *GRPCConn) Disconnect(ctx context.Context) error {
	if conn == nil {
		return fmt.Errorf("%w: grpconn", errdefs.ErrNilObjectCall)
	}

	conn.mu.Lock()
	defer conn.mu.Unlock()

	if conn.conn == nil {
		return nil
	}

	err := conn.conn.Close()
	conn.conn = nil

	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: disconnect grpc conn: %v", errdefs.ErrGRPC, err)
}

func (conn *GRPCConn) Close(ctx context.Context) error {
	if conn == nil {
		return nil
	}
	return conn.Disconnect(ctx)
}

func (conn *GRPCConn) tryConnect(ctx context.Context, target string, log *zap.Logger) error {
	initFn := func(ctx context.Context) error {
		c, err := grpc.NewClient(target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			err = fmt.Errorf("%w: %v", errdefs.ErrGRPC, err)
			log.Warn("retry: connect grpc", zap.Error(err))
			return err
		}

		healthClient := healthpb.NewHealthClient(c)
		_, err = healthClient.Check(ctx, &healthpb.HealthCheckRequest{})

		if err != nil {
			// healthchech Unimplemented means server is available
			st, ok := healthst.FromError(err)
			if !ok || st.Code() != codes.Unimplemented {
				err = fmt.Errorf("%w: %v", errdefs.ErrGRPC, err)
				log.Warn("retry: connect grpc", zap.Error(err))
				return err
			}
		}

		// mark as initialized
		conn.conn = c

		return nil
	}

	if err := retry.RetryInfinite(ctx, initFn, 250*time.Millisecond); err != nil {
		return fmt.Errorf("GRPC conn: try connect: %w", err)
	}

	return nil
}

func (conn *GRPCConn) Invoke(
	ctx context.Context,
	method string,
	args any,
	reply any,
	opts ...grpc.CallOption,
) error {
	if err := conn.checkConnReady(); err != nil {
		return fmt.Errorf("grpc connection: invoke: %w", err)
	}
	return conn.conn.Invoke(ctx, method, args, reply, opts...)
}

func (conn *GRPCConn) NewStream(
	ctx context.Context,
	desc *grpc.StreamDesc,
	method string,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	if err := conn.checkConnReady(); err != nil {
		return nil, fmt.Errorf("grpc connection: new stream: %w", err)
	}
	return conn.conn.NewStream(ctx, desc, method, opts...)
}

func (conn *GRPCConn) checkConnReady() error {
	if conn == nil {
		return fmt.Errorf("%w: grpc connection", errdefs.ErrNilObjectCall)
	}

	conn.mu.RLock()
	defer conn.mu.RUnlock()
	if conn.conn == nil {
		return errdefs.ErrGRPCNotReady
	}
	return nil
}
