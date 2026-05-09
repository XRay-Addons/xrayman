package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Session struct {
	ID        uuid.UUID
	CreatedAt time.Time
	ExpiresAt time.Time
	Revoked   bool
}
