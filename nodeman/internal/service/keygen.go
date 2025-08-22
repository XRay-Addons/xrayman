package service

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/google/uuid"
)

func generateVlessUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errdefs.WrapWithStack(err)
	}
	return id.String(), nil
}
