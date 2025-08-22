package xrayapi

import (
	"context"
	"sync"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/grpcconn"
	"github.com/XRay-Addons/xrayman/node/internal/infra/tx"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	handlerService "github.com/xtls/xray-core/app/proxyman/command"
	statsService "github.com/xtls/xray-core/app/stats/command"
	"github.com/xtls/xray-core/common/protocol"
	"go.uber.org/zap"
)

// TODO: WithLog
type XRayApi struct {
	inbounds []models.Inbound
	apiConn  *grpcconn.GRPCConn
	hsClient handlerService.HandlerServiceClient
	ssClient statsService.StatsServiceClient

	mu sync.Mutex
}

func New(apiURL string, inbounds []models.Inbound, log *zap.Logger) (*XRayApi, error) {
	if log == nil {
		return nil, errdefs.NewNilArg("log")
	}
	apiConn, err := grpcconn.New(apiURL, log)
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

	if err := api.apiConn.Close(ctx); err != nil {
		return errdefs.WrapWithStack(err)
	}
	api.apiConn = nil

	return nil
}

func (api *XRayApi) Connect(ctx context.Context) error {
	if api == nil {
		return errdefs.NewNilCall()
	}

	api.mu.Lock()
	defer api.mu.Unlock()

	if err := api.apiConn.Connect(ctx); err != nil {
		return errdefs.WrapWithStack(err)
	}
	return nil
}

func (api *XRayApi) Disconnect(ctx context.Context) error {
	if api == nil {
		return errdefs.NewNilCall()
	}

	api.mu.Lock()
	defer api.mu.Unlock()

	if err := api.apiConn.Disconnect(ctx); err != nil {
		return errdefs.WrapWithStack(err)
	}
	return nil
}

func (api *XRayApi) EditUsers(
	ctx context.Context,
	add, remove []models.User,
) error {
	if api == nil || api.hsClient == nil {
		return errdefs.NewNilCall()
	}

	api.mu.Lock()
	defer api.mu.Unlock()

	var editUsersTx tx.Tx
	for _, in := range api.inbounds {
		in := in
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
		return errdefs.NewNilCall()
	}

	api.mu.Lock()
	defer api.mu.Unlock()

	return ping(ctx, api.ssClient)
}
