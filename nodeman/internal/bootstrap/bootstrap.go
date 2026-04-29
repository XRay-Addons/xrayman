package bootstrap

import "context"

type Config struct {
	AdminPassword string
}

func Bootstrap(ctx context.Context, cfg Config, auth AuthService) error {
	if cfg.AdminPassword != "" {
		if err := setupAdmin(ctx, cfg.AdminPassword, auth); err != nil {
			return err
		}
	}
	return nil
}