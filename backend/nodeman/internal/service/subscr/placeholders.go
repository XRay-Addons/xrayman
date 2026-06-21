package subscr

import (
	"fmt"

	"github.com/valyala/fasttemplate"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

const (
	UserIDPlaceholder          = "UserID"
	UserNamePlaceholder        = "UserName"
	UserDisplayNamePlaceholder = "DisplayName"
	TotalUploadPlaceholder     = "TotalUpload"
	TotalDownloadPlaceholder   = "TotalDownload"
)

func replacePlaceholders(s string, u *models.UserView) string {
	return fasttemplate.New(s, "{{", "}}").ExecuteString(map[string]interface{}{
		UserIDPlaceholder:          fmt.Sprintf("%v", u.User.Profile.ID),
		UserNamePlaceholder:        u.User.Profile.Name,
		UserDisplayNamePlaceholder: u.User.Profile.DisplayName,
		TotalUploadPlaceholder:     u.Traffic.Total.Upload,
		TotalDownloadPlaceholder:   u.Traffic.Total.Download,
	})
}

func listPlaceholders() []string {
	return []string{
		makePlaceholder(UserIDPlaceholder),
		makePlaceholder(UserNamePlaceholder),
		makePlaceholder(UserDisplayNamePlaceholder),
		makePlaceholder(TotalUploadPlaceholder),
		makePlaceholder(TotalDownloadPlaceholder),
	}
}

func makePlaceholder(h string) string {
	return "{{" + h + "}}"
}
