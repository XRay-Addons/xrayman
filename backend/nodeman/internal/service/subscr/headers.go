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
)

func replacePlaceholders(hs models.Headers, u models.User) models.Headers {
	replaced := make([]models.Header, 0, len(hs))
	for _, h := range hs {
		t := fasttemplate.New(h.Value, "{{", "}}")
		v := t.ExecuteString(map[string]interface{}{
			UserIDPlaceholder:          fmt.Sprintf("%v", u.Profile.ID),
			UserNamePlaceholder:        u.Profile.Name,
			UserDisplayNamePlaceholder: u.Profile.DisplayName,
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
	}
}

func makePlaceholder(h string) string {
	return "{{" + h + "}}"
}
