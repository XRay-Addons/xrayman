package bootstrap

import "context"

type AuthService interface {
	SetAdmin(ctx context.Context, password string) error
}
