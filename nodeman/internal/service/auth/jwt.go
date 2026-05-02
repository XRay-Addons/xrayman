package auth

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type JWT interface {
	GenerateToken() (models.AuthResult, error)
}
