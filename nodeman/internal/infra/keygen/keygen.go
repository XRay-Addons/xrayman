package keygen

import (
	"crypto/rand"
	"fmt"

	"github.com/google/uuid"
)

type Keygen struct {
}

func New() *Keygen {
	return &Keygen{}
}

func (kg *Keygen) GenerateHS256Secret() ([]byte, error) {
	const size = 32
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("failed to generate random secret: %w", err)
	}
	return b, nil
}

func (kg *Keygen) GenerateVlessUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
