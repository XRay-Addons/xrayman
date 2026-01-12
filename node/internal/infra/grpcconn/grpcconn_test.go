package grpcconn

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func newTestGRPCServer(t *testing.T, delay time.Duration) (addr string, cancel func()) {
	t.Helper()

	grpcServer := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, &grpc_health_v1.UnimplementedHealthServer{})

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Logf("failed to listen: %v", err)
	}

	go func() {
		time.Sleep(delay)
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			t.Logf("grpc server exited: %v", serveErr)
		}
	}()

	return lis.Addr().String(), func() {
		grpcServer.GracefulStop()
	}
}

func TestGRPCConn(t *testing.T) {
	addr, stop := newTestGRPCServer(t, 5*time.Second)
	defer stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log := zaptest.NewLogger(t)

	conn, err := New(addr, log)
	if err != nil {
		t.Fatalf("failed to connect to test grpc server: %v", err)
	}
	defer conn.Close(context.TODO())

	err = conn.Connect(ctx)
	require.NoError(t, err)

	err = conn.Disconnect(ctx)
	require.NoError(t, err)
}

func TestGRPCConnTimeout(t *testing.T) {
	addr, stop := newTestGRPCServer(t, 5*time.Second)
	defer stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	log := zaptest.NewLogger(t)

	conn, err := New(addr, log)
	if err != nil {
		t.Fatalf("failed to connect to test grpc server: %v", err)
	}
	defer conn.Close(context.TODO())

	err = conn.Connect(ctx)
	require.True(t, errors.Is(err, context.DeadlineExceeded))

	err = conn.Disconnect(ctx)
	require.NoError(t, err)
}
