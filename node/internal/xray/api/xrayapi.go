package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/tx"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	handlerService "github.com/xtls/xray-core/app/proxyman/command"
	statsService "github.com/xtls/xray-core/app/stats/command"
	"github.com/xtls/xray-core/common/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type XRayApi struct {
	inbounds []models.Inbound
	apiConn  *grpc.ClientConn
	hsClient handlerService.HandlerServiceClient
	ssClient statsService.StatsServiceClient

	mu sync.Mutex
}

func New(apiURL string, inbounds []models.Inbound) (*XRayApi, error) {
	apiConn, err := grpc.NewClient(apiURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: connect xray api: %v", errdefs.ErrXRay, err)
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

func (api *XRayApi) Close() error {
	if api == nil {
		return nil
	}

	api.mu.Lock()
	defer api.mu.Unlock()

	if api.apiConn == nil {
		return nil
	}

	api.hsClient = nil
	api.ssClient = nil

	if err := api.apiConn.Close(); err != nil {
		return fmt.Errorf("%w: api connection closing: %v", errdefs.ErrXRay, err)
	}
	api.apiConn = nil

	return nil
}

func (api *XRayApi) EditUsers(
	ctx context.Context,
	add, remove []models.User,
) error {
	if api == nil || api.hsClient == nil {
		return fmt.Errorf("%w: xray api", errdefs.ErrNilObjectCall)
	}

	var editUsersTx tx.Tx
	for _, in := range api.inbounds {
		for _, u := range add {
			inUser, err := getInboundUser(u, in.Type)
			if err != nil {
				return fmt.Errorf("edit users: %w", err)
			}
			editUsersTx.AddItem(
				api.addFn(in.Tag, inUser),
				api.removeFn(in.Tag, inUser.Email),
			)
		}
		for _, u := range remove {
			inUser, err := getInboundUser(u, in.Type)
			if err != nil {
				return fmt.Errorf("edit users: %w", err)
			}
			editUsersTx.AddItem(
				api.removeFn(in.Tag, inUser.Email),
				api.addFn(in.Tag, inUser),
			)
		}
	}

	if err := editUsersTx.Run(ctx); err != nil {
		return fmt.Errorf("edit users: %w", err)
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
		return fmt.Errorf("%w: xray api", errdefs.ErrNilObjectCall)
	}

	return ping(ctx, api.ssClient)
}
