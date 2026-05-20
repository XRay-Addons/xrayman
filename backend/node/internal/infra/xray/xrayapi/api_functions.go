package xrayapi

import (
	"context"
	"fmt"
	"strconv"
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

type UserStats struct {
	ID       models.UserID
	Uplink   int64
	Downlink int64
}

func getStats(
	ctx context.Context,
	ssClient statsService.StatsServiceClient,
	log *zap.Logger,
) error {
	reset := true
	resp, err := ssClient.QueryStats(context.Background(), &statsService.QueryStatsRequest{
		Reset_: reset,
	})
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	// Get traffic data
	stat := resp.GetStat()
	const splitTag = ">>>"
	const userTag = "user"
	const trafficTag = "traffic"

	userStatsMap := make(map[int]UserStats)
	for _, s := range stat {
		parts := strings.Split(s.Name, splitTag)
		if len(parts) != 4 || parts[0] != userTag || parts[2] != trafficTag {
			log.Warn("unparsed stat", zap.String("name", s.Name))
			continue
		}
		userID, _, err := extractUser(parts[1])
		if err != nil {
			log.Warn("unparsed user", zap.String("name", s.Name))
			continue
		}
		direction, err := extractDirection(parts[3])
		if err != nil {
			log.Warn("unparsed direction", zap.String("name", s.Name))
			continue
		}
		userStat := userStatsMap[userID]
		userStat.ID = models.UserID(userID)
		switch direction {
		case UplinkDirection:
			userStat.Uplink += s.Value
		case DownlinkDirection:
			userStat.Downlink += s.Value
		}
		userStatsMap[userID] = userStat
	}
	usersStats := make([]UserStats, 0, len(userStatsMap))
	for _, v := range userStatsMap {
		usersStats = append(usersStats, v)
	}
	for _, v := range usersStats {
		log.Info("stats", zap.Int("user", int(v.ID)),
			zap.Int64("in", v.Downlink), zap.Int64("out", v.Uplink))
	}

	return nil
}

func extractUser(s string) (userID int, userName string, err error) {
	defer func() {
		if err != nil {
			err = xerr.WrapWithf(err, "stats string: %s", s)
		}
	}()

	s = strings.TrimSpace(s)

	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return 0, "", xerr.New("invalid user format, expected id-name")
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", xerr.Wrap(err, xerr.WithStack())
	}

	name := strings.TrimSpace(parts[1])
	if name == "" {
		return 0, "", xerr.New("empty user name")
	}

	return id, name, nil
}

type Direction int

const (
	DirectionUnknown Direction = iota
	UplinkDirection
	DownlinkDirection
)

func extractDirection(s string) (Direction, error) {
	const uplinkTag = "uplink"
	const downlinkTag = "downlink"

	switch strings.ToLower(strings.TrimSpace(s)) {
	case uplinkTag:
		return UplinkDirection, nil
	case downlinkTag:
		return DownlinkDirection, nil
	default:
		return 0, xerr.Newf("unknown direction: %s", s)
	}
}
