package keygen

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

type Keygen struct {
}

func New() *Keygen {
	return &Keygen{}
}

func GenerateHS256Secret() (string, error) {
	const size = 32
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random secret: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (kg *Keygen) GenerateVlessUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
