package bootstrap

import "context"

type AuthService interface {
	AuthAdmin(ctx context.Context, password string) error
}
