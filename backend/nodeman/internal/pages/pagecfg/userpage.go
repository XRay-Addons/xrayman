package pagecfg

import "github.com/XRay-Addons/xrayman/nodeman/internal/models"

type UserPageCfg struct {
	ApiPrefix   string
	UserPrefix  string
	SupportLink string
	AppLinks    []models.AppLink
}
