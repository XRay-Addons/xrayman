package xrayapi

import (
	"context"
	"sync"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/grpcconn"
	"github.com/XRay-Addons/xrayman/node/internal/infra/tx"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	handlerService "github.com/xtls/xray-core/app/proxyman/command"
	statsService "github.com/xtls/xray-core/app/stats/command"
	"github.com/xtls/xray-core/common/protocol"
	"go.uber.org/zap"
)

type XRayApi struct {
	inbounds []models.Inbound
	apiConn  *grpcconn.GRPCConn
	hsClient handlerService.HandlerServiceClient
	ssClient statsService.StatsServiceClient

	mu      sync.Mutex
	timeout time.Duration
}

func WithLogger(logger *zap.Logger) option {
	return func(o *options) {
		if logger == nil {
			return
		}
		o.log = logger
	}
}

func WithTimeout(t time.Duration) option {
	return func(o *options) {
		o.timeout = t
	}
}

type option func(o *options)

type options struct {
	log     *zap.Logger
	timeout time.Duration
}

const defaultTimeout = 5 * time.Second

func New(apiURL string, inbounds []models.Inbound, opts ...option) (*XRayApi, error) {
	o := &options{
		log:     zap.NewNop(),
		timeout: defaultTimeout,
	}
	for _, opt := range opts {
		opt(o)
	}

	apiConn, err := grpcconn.New(apiURL, o.log)
	if err != nil {
		return nil, err
	}

	hsClient := handlerService.NewHandlerServiceClient(apiConn)
	ssClient := statsService.NewStatsServiceClient(apiConn)

	return &XRayApi{
		inbounds: inbounds,
		apiConn:  apiConn,
		hsClient: hsClient,
		ssClient: ssClient,
		timeout:  o.timeout,
	}, nil
}

func (api *XRayApi) Close(ctx context.Context) error {
	if api == nil {
		return nil
	}

	if api.apiConn == nil {
		return nil
	}

	api.mu.Lock()
	defer api.mu.Unlock()

	api.hsClient = nil
	api.ssClient = nil

	ctx, cancel := context.WithTimeout(ctx, api.timeout)
	defer cancel()

	if err := api.apiConn.Close(ctx); err != nil {
		return xerr.WrapWithStack(err)
	}
	api.apiConn = nil

	return nil
}

func (api *XRayApi) EditUsers(
	ctx context.Context,
	add, remove []models.User,
) error {
	if api == nil || api.hsClient == nil {
		return errdefs.NilCall()
	}

	api.mu.Lock()
	defer api.mu.Unlock()

	var editUsersTx tx.Tx
	for _, in := range api.inbounds {
		for _, u := range add {
			inUser, err := getInboundUser(u, in.Type)
			if err != nil {
				return err
			}
			editUsersTx.AddItem(
				api.addFn(in.Tag, inUser),
				api.removeFn(in.Tag, inUser.Email),
			)
		}
		for _, u := range remove {
			inUser, err := getInboundUser(u, in.Type)
			if err != nil {
				return err
			}
			editUsersTx.AddItem(
				api.removeFn(in.Tag, inUser.Email),
				api.addFn(in.Tag, inUser),
			)
		}
	}

	ctx, cancel := context.WithTimeout(ctx, api.timeout)
	defer cancel()

	if err := editUsersTx.Run(ctx); err != nil {
		return err
	}

	return nil
}

func (api *XRayApi) addFn(inTag string, user *protocol.User) tx.Fn {
	return func(ctx context.Context) error {
		return addUser(ctx, api.hsClient, inTag, user)
	}
}

func (api *XRayApi) removeFn(inTag string, userEmail string) tx.Fn {
	return func(ctx context.Context) error {
		return removeUser(ctx, api.hsClient, inTag, userEmail)
	}
}

func (api *XRayApi) Ping(ctx context.Context) error {
	if api == nil || api.ssClient == nil {
		return errdefs.NilCall()
	}

	api.mu.Lock()
	defer api.mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, api.timeout)
	defer cancel()

	return ping(ctx, api.ssClient)
}
