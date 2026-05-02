package auth

import "context"

type Password interface {
	Verify(ctx context.Context, password string) error
}
