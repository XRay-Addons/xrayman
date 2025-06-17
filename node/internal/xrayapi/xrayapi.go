package xrayapi

import (
	"context"
	"fmt"
	"sync"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/shared/models"
	handlerService "github.com/xtls/xray-core/app/proxyman/command"
	statsService "github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type XRayApi struct {
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
		return nil, fmt.Errorf("failed to connect to Xray API: %w", err)
	}

	hsClient := handlerService.NewHandlerServiceClient(apiConn)
	ssClient := statsService.NewStatsServiceClient(apiConn)

	return &XRayApi{
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

func (api *XRayApi) AddUsers(
	ctx context.Context,
	ins []models.Inbound,
	users []models.User,
) error {
	if api == nil {
		return fmt.Errorf("%w: api not exists", errdefs.ErrIPE)
	}

	api.mu.Lock()
	defer api.mu.Unlock()

	return AddUsers(ctx, api.hsClient, ins, users)
}

func (api *XRayApi) DelUsers(
	ctx context.Context,
	ins []models.Inbound,
	users []models.User,
) error {
	if api == nil {
		return fmt.Errorf("%w: api not exists", errdefs.ErrIPE)
	}

	api.mu.Lock()
	defer api.mu.Unlock()

	return DelUsers(ctx, api.hsClient, ins, users)
}

func (api *XRayApi) Ping(ctx context.Context) error {
	if api == nil {
		return fmt.Errorf("%w: api not exists", errdefs.ErrIPE)
	}

	api.mu.Lock()
	defer api.mu.Unlock()

	return Ping(ctx, api.ssClient)
}

/*func (api *XRayApi) addUser(user models.User, in models.Inbound) error {
	addUserOp := func(u *protocol.User) *serial.TypedMessage {
		return serial.ToTypedMessage(&handlerService.AddUserOperation{User: u})
	}
	err := s.alterUser(user, in, addUserOp)

	// already exists is not an error for us
	alreadyExistsErrPattern := fmt.Sprintf("User %s already exists", user.Name)
	if err != nil && strings.Contains(err.Error(), alreadyExistsErrPattern) {
		return nil
	}

	return err
}

func (api *XRayApi) delUser(user models.User, in models.Inbound) error {
	delUserOp := func(u *protocol.User) *serial.TypedMessage {
		return serial.ToTypedMessage(&handlerService.RemoveUserOperation{Email: u.Email})
	}

	err := s.alterUser(user, in, delUserOp)

	// not exists is not an error for us
	notFoundErrPattern := fmt.Sprintf("User %s not found", user.Email)
	if err != nil && strings.Contains(err.Error(), notFoundErrPattern) {
		s.log.Warnf("xray api: user %s not exists", user.Email)
		return nil
	}

	return err
}*/

/*func (s *XRayApiService) ListUserInbounds(email string) ([]models.VlessInbound, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.connected {
		return nil, fmt.Errorf("xray api not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	userInbounds := make([]models.VlessInbound, 0)
	for _, in := range s.config.ServerInbounds.Vless {
		res, err := s.hsClient.GetInboundUsers(ctx, &handlerService.GetInboundUserRequest{
			Tag:   in.Tag,
			Email: email,
		})
		if err != nil {
			return nil, err
		}
		if len(res.GetUsers()) > 0 {
			userInbounds = append(userInbounds, in)
		}
	}

	return userInbounds, nil
}

func (s *XRayApiService) GetUserStats(email string) (*models.UserStats, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.connected {
		return nil, fmt.Errorf("xray api not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	downlinkStats, err := s.ssClient.QueryStats(ctx, &statsService.QueryStatsRequest{
		Pattern: fmt.Sprintf("user>>>%s>>>traffic>>>downlink", email),
		Reset_:  false,
	})
	if err != nil {
		return nil, err
	}
	uplinkStats, err := s.ssClient.QueryStats(ctx, &statsService.QueryStatsRequest{
		Pattern: fmt.Sprintf("user>>>%s>>>traffic>>>uplink", email),
		Reset_:  false,
	})
	if err != nil {
		return nil, err
	}

	uplink := int64(0)
	downlink := int64(0)
	for _, stat := range uplinkStats.GetStat() {
		uplink += stat.Value
	}
	for _, stat := range downlinkStats.GetStat() {
		downlink += stat.Value
	}

	return &models.UserStats{
		Uplink:   uplink,
		Downlink: downlink,
	}, nil
}

func (api *XRayApi) addUser(user models.User, in models.Inbound) error {
	addUserOp := func(u *protocol.User) *serial.TypedMessage {
		return serial.ToTypedMessage(&handlerService.AddUserOperation{User: u})
	}
	err := api.alterUser(user, in, addUserOp)

	// already exists is not an error for us
	alreadyExistsErrPattern := fmt.Sprintf("User %s already exists", user.Name)
	if err != nil && strings.Contains(err.Error(), alreadyExistsErrPattern) {
		return nil
	}

	return err
}

func (api *XRayApi) delUser(user models.User, in models.Inbound) error {
	delUserOp := func(u *protocol.User) *serial.TypedMessage {
		return serial.ToTypedMessage(&handlerService.RemoveUserOperation{Email: u.Email})
	}

	err := api.alterUser(user, in, delUserOp)

	// not exists is not an error for us
	notFoundErrPattern := fmt.Sprintf("User %s not found", user.Email)
	if err != nil && strings.Contains(err.Error(), notFoundErrPattern) {
		s.log.Warnf("xray api: user %s not exists", user.Email)
		return nil
	}

	return err
}

type ApiUserOp = func(*protocol.User) *serial.TypedMessage

func (s *XRayApiService) alterUser(u models.User, in models.VlessInbound, op ApiUserOp) error {
	if s.cmdConn == nil {
		return fmt.Errorf("server api not exist")
	}

	config, err := s.toUserConfig(u, in)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	_, err = s.hsClient.AlterInbound(ctx, &handlerService.AlterInboundRequest{
		Tag:       in.Tag,
		Operation: op(config),
	})

	return err
}

func (s *XRayApiService) toUserConfig(u models.User, in models.Inbound) (*protocol.User, error) {
	var account *vless.Account
	if in.StreamSettings.Network == "tcp" && in.StreamSettings.Security == "reality" {
		account = &vless.Account{
			Id:         u.Properties.VlessUUID,
			Encryption: "none",
			Flow:       "xtls-rprx-vision",
		}
	} else if in.StreamSettings.Network == "xhttp" && in.StreamSettings.Security == "" {
		account = &vless.Account{
			Id:         u.Properties.VlessUUID,
			Encryption: "none",
		}
	} else {
		return nil, fmt.Errorf("unsupported in type %s", in)
	}

	return &protocol.User{
		Email:   u.Email,
		Account: serial.ToTypedMessage(account),
	}, nil
}*/
