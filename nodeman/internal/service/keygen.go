package service

import (
	"github.com/google/uuid"
)

func generateVlessUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
