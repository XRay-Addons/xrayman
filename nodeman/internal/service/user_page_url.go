package service

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func makeUserPageURL(u models.User) string {
	return fmt.Sprintf("./%d-%s", u.ID, u.Profile.SlugName)
}
