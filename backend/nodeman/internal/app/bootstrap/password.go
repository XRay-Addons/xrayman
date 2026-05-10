package bootstrap

import (
	"context"
)

type Password interface {
	Update(ctx context.Context, password string) error
}
