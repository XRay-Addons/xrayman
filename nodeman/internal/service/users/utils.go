package users

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/xerr"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

func generateVlessUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", xerr.WrapWithStack(err)
	}
	return id.String(), nil
}

func makeSlugName(name string) string {
	return slug.Make(name)
}
