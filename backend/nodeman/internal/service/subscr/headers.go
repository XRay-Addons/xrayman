package subscr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

const (
	ProfileTitleHeader          = "profile-title"
	ProfileUpdateIntervalHeader = "profile-update-interval"
	TgPageHeader                = "support-url"
	WebPageHeader               = "profile-web-page-url"
	AnnounceHeader              = "announce"
	RoutingHeader               = "routing"
	TrafficStatsHeader          = "subscription-userinfo"
	TrafficStatsFmt             = "upload=%d; download=%d; total=0; expire=0"
)

func createClientHeaders(ctx context.Context,
	u *models.UserView, settings *models.Settings,
) models.SubHeaders {
	var headers []models.SubHeader

	// title
	if settings.SubscrTitle != "" {
		headers = append(headers, models.SubHeader{
			Key:   ProfileTitleHeader,
			Value: replacePlaceholders(settings.SubscrTitle, u),
		})
	}
	// update interval
	if settings.UpdateInterval != 0 {
		headers = append(headers, models.SubHeader{
			Key:   ProfileUpdateIntervalHeader,
			Value: strconv.Itoa(settings.UpdateInterval),
		})
	}
	// tg page
	if settings.TgPage != "" {
		headers = append(headers, models.SubHeader{
			Key:   TgPageHeader,
			Value: settings.TgPage,
		})
	}
	// web page
	if settings.UserPage != "" {
		headers = append(headers, models.SubHeader{
			Key:   WebPageHeader,
			Value: replacePlaceholders(settings.UserPage, u),
		})
	}
	// announce header
	if settings.UsersMessage != "" {
		headers = append(headers, models.SubHeader{
			Key:   AnnounceHeader,
			Value: replacePlaceholders(settings.UsersMessage, u),
		})
	}
	// routing header
	if settings.Routing != "" {
		headers = append(headers, models.SubHeader{
			Key:   RoutingHeader,
			Value: settings.Routing,
		})
	}
	// custom headers
	headers = append(headers, settings.CustomHeaders...)
	// traffic stats header
	ts := u.Traffic.Total
	headers = append(headers, models.SubHeader{
		Key:   TrafficStatsHeader,
		Value: fmt.Sprintf(TrafficStatsFmt, ts.Upload, ts.Download),
	})

	return headers
}
