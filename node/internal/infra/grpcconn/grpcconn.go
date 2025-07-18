package grpcconn

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
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
	conn *grpc.ClientConn

	// for initialization loop
	initialized atomic.Bool
	wg          sync.WaitGroup
	cancel      context.CancelFunc
}

var _ grpc.ClientConnInterface = (*GRPCConn)(nil)

func New(target string, log *zap.Logger) (*GRPCConn, error) {
	if log == nil {
		return nil, fmt.Errorf("%w: init grpc client: logger", errdefs.ErrNilObjectCall)
	}
	ctx, cancel := context.WithCancel(context.Background())

	grpcConn := &GRPCConn{
		cancel: cancel,
	}

	// run installing service loop
	grpcConn.wg.Add(1)
	go grpcConn.createConnLoop(ctx, target, log)

	return grpcConn, nil
}

func (conn *GRPCConn) Close() error {
	if conn == nil {
		return nil
	}

	conn.cancel()
	conn.wg.Wait()

	if !conn.initialized.Load() {
		return nil
	}

	err := conn.conn.Close()
	conn.initialized.Store(false)

	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: close grpc conn: %v", errdefs.ErrGRPC, err)
}

func (conn *GRPCConn) createConnLoop(ctx context.Context, target string, log *zap.Logger) {
	defer conn.wg.Done()

	initFn := func(ctx context.Context) error {
		c, err := grpc.NewClient(target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			err = fmt.Errorf("%w: create grpc client: %v", errdefs.ErrGRPC, err)
			log.Warn(err.Error())
			return err
		}

		healthClient := healthpb.NewHealthClient(c)
		_, err = healthClient.Check(ctx, &healthpb.HealthCheckRequest{})

		if err != nil {
			// healthchech Unimplemented means server is available
			st, ok := healthst.FromError(err)
			if !ok || st.Code() != codes.Unimplemented {
				err = fmt.Errorf("%w: healthcheck grpc client: %v", errdefs.ErrGRPC, err)
				log.Warn(err.Error())
				return err
			}
		}

		// mark as initialized
		conn.conn = c
		conn.initialized.Store(true)

		return nil
	}

	retry.RetryInfinite(ctx, initFn, 1*time.Second)
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
	if !conn.initialized.Load() {
		return errdefs.ErrGRPCNotReady
	}
	return nil
}
