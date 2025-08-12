package service

import (
	"crypto/rand"
	"fmt"

	"github.com/google/uuid"
)

func generateHS256Secret() ([]byte, error) {
	const size = 32
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("failed to generate random secret: %w", err)
	}
	return b, nil
}

func generateVlessUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
