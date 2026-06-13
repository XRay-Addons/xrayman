package xrayapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/node/internal/models"
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

const (
	userPattern = "user>>>"
	splitTag    = ">>>"
	userTag     = "user"
	trafficTag  = "traffic"
	uplinkTag   = "uplink"
	downlinkTag = "downlink"
)

func getStats(
	ctx context.Context,
	ssClient statsService.StatsServiceClient,
	log *zap.Logger,
) (*models.StatsResult, error) {
	resp, err := ssClient.QueryStats(context.Background(), &statsService.QueryStatsRequest{
		Pattern: userPattern,
		Reset_:  true,
	})
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	// Get traffic data
	stat := resp.GetStat()
	userStatsMap := make(map[int]models.UserStats)
	for _, s := range stat {
		parts := strings.Split(s.Name, splitTag)
		if len(parts) != 4 || parts[0] != userTag || parts[2] != trafficTag {
			log.Warn("unparsed stat", zap.String("name", s.Name))
			continue
		}
		userID, _, err := models.ParseVlessEmail(parts[1])
		if err != nil {
			log.Warn("unparsed user", zap.String("name", s.Name))
			continue
		}

		userStat := userStatsMap[userID]
		userStat.ID = userID

		switch parts[3] {
		case uplinkTag:
			userStat.Upload = s.Value
		case downlinkTag:
			userStat.Download = s.Value
		default:
			log.Warn("unparsed direction", zap.String("tag", parts[3]))
		}

		userStatsMap[userID] = userStat
	}

	usersStats := make([]models.UserStats, 0, len(userStatsMap))
	for _, v := range userStatsMap {
		usersStats = append(usersStats, v)
	}

	return &models.StatsResult{
		Users: usersStats,
	}, nil
}
