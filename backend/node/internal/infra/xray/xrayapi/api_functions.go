package xrayapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/XRay-Addons/xrayman/common/xerr"
	handlerService "github.com/xtls/xray-core/app/proxyman/command"
	statsService "github.com/xtls/xray-core/app/stats/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"go.uber.org/zap"
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

	return xerr.WrapWithStack(err)
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
	return xerr.WrapWithStack(err)
}

func ping(
	ctx context.Context,
	ssClient statsService.StatsServiceClient,
) error {
	_, err := ssClient.GetSysStats(ctx, &statsService.SysStatsRequest{})
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func getStats(
	ctx context.Context,
	ssClient statsService.StatsServiceClient,
	log *zap.Logger,
) error {
	resp, err := ssClient.QueryStats(context.Background(), &statsService.QueryStatsRequest{
		// 这里是查询语句，例如 “user>>>love@xray.com>>>traffic>>>uplink” 表示查询用户 email 为 love@xray.com 在所有入站中的上行流量
		// Pattern: ptn,
		// 是否重置流量信息(true, false)，即完成查询后是否把流量统计归零
		// Reset_: reset, // reset traffic data everytime
	})
	if err != nil {
		return nil
	}
	// Get traffic data
	stat := resp.GetStat()
	log.Info("stats", zap.String("stats", fmt.Sprintf("%+v", stat)))

	return nil
}
