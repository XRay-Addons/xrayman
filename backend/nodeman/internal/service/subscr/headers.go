package subscr

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

const (
	ProfileTitleHeader          = "profile-title"
	ProfileUpdateIntervalHeader = "profile-update-interval"
	TgPageHeader                = "support-url"
	WebPageHeader               = "profile-web-page-url"
	AnnounceHeader              = "announce"
	RoutingHeader               = "routing"
	TrafficStatsHeader          = "subscription-userinfo" // upload=0; download=2153701362; total=0; expire=1790951622
	TrafficStatsFmt             = "upload=%d; download=%d"
)

func createClientHeaders(ctx context.Context,
	u *models.UserView, cfg *models.DynamicConfig,
) models.SubHeaders {
	var headers []models.SubHeader

	// title
	if cfg.SubscrTitle != "" {
		headers = append(headers, models.SubHeader{
			Key:   ProfileTitleHeader,
			Value: replacePlaceholders(cfg.SubscrTitle, u),
		})
	}
	// update interval
	if cfg.UpdateInterval != 0 {
		headers = append(headers, models.SubHeader{
			Key:   ProfileUpdateIntervalHeader,
			Value: string(cfg.UpdateInterval),
		})
	}
	// tg page
	if cfg.TgPage != "" {
		headers = append(headers, models.SubHeader{
			Key:   TgPageHeader,
			Value: cfg.TgPage,
		})
	}
	// web page
	if cfg.UserPage != "" {
		headers = append(headers, models.SubHeader{
			Key:   WebPageHeader,
			Value: replacePlaceholders(cfg.UserPage, u),
		})
	}
	// announce header
	if cfg.UsersMessage != "" {
		headers = append(headers, models.SubHeader{
			Key:   AnnounceHeader,
			Value: replacePlaceholders(cfg.UsersMessage, u),
		})
	}
	// announce header
	if cfg.Routing != "" {
		headers = append(headers, models.SubHeader{
			Key:   RoutingHeader,
			Value: cfg.Routing,
		})
	}
	// traffic stats header
	ts := u.Traffic.Total
	headers = append(headers, models.SubHeader{
		Key:   TrafficStatsHeader,
		Value: fmt.Sprintf(TrafficStatsFmt, ts.Upload, ts.Download),
	})
	return headers
}
