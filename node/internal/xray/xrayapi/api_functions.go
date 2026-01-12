package xrayapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	handlerService "github.com/xtls/xray-core/app/proxyman/command"
	statsService "github.com/xtls/xray-core/app/stats/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
)

func addUser(
	ctx context.Context,
	hs handlerService.HandlerServiceClient,
	inboundTag string,
	user *protocol.User,
) error {
	_, err := hs.AlterInbound(ctx, &handlerService.AlterInboundRequest{
		Tag:       inboundTag,
		Operation: serial.ToTypedMessage(&handlerService.AddUserOperation{User: user}),
	})

	alreadyExistsErrPattern := fmt.Sprintf("User %s already exists", user.Email)
	if err == nil || strings.Contains(err.Error(), alreadyExistsErrPattern) {
		return nil
	}

	return fmt.Errorf("%w: add user: %v", errdefs.ErrGRPC, err)
}

func removeUser(
	ctx context.Context,
	hs handlerService.HandlerServiceClient,
	inboundTag string,
	email string,
) error {
	_, err := hs.AlterInbound(ctx, &handlerService.AlterInboundRequest{
		Tag:       inboundTag,
		Operation: serial.ToTypedMessage(&handlerService.RemoveUserOperation{Email: email}),
	})

	notFoundErrPattern := fmt.Sprintf("User %s not found", email)
	if err != nil && strings.Contains(err.Error(), notFoundErrPattern) {
		return nil
	}
	return fmt.Errorf("%w: remove user: %v", errdefs.ErrGRPC, err)
}

func ping(
	ctx context.Context,
	ssClient statsService.StatsServiceClient,
) error {
	_, err := ssClient.GetSysStats(ctx, &statsService.SysStatsRequest{})
	if err != nil {
		return fmt.Errorf("%w: ping call: %v", errdefs.ErrGRPC, err)
	}

	return nil
}
