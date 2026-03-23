package service

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

func generateVlessUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errdefs.WrapWithStack(err)
	}
	return id.String(), nil
}

func makeSlugName(name string) string {
	return slug.Make(name)
}

func makeUserPageURL(id models.UserID, name string) string {
	return fmt.Sprintf("./%d-%s", id, name)
}
