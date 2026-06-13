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

func replacePlaceholders(hs models.Headers, u models.UserView) models.Headers {
	replaced := make([]models.Header, 0, len(hs))
	for _, h := range hs {
		t := fasttemplate.New(h.Value, "{{", "}}")
		v := t.ExecuteString(map[string]interface{}{
			UserIDPlaceholder:          fmt.Sprintf("%v", u.User.Profile.ID),
			UserNamePlaceholder:        u.User.Profile.Name,
			UserDisplayNamePlaceholder: u.User.Profile.DisplayName,
			TotalUploadPlaceholder:     u.Traffic.Total.Upload,
			TotalDownloadPlaceholder:   u.Traffic.Total.Download,
		})
		replaced = append(replaced, models.Header{
			ID:    h.ID,
			Key:   h.Key,
			Value: v,
		})
	}
	return replaced
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
