// AI generated tests
package dbstorage

/*import (
	"context"
	"testing"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/auth/password"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestAdminAuth(t *testing.T) {
	logger := zaptest.NewLogger(t)

	s, _, _, cleanup := setupTestDB(t, logger)
	defer cleanup()
	logger.Info("new test db inited")

	ctx := context.Background()

	// 1. GetAuth on empty db
	err := s.PasswordStorage().DoUoW(ctx, func(uowctx password.UoWContext) error {
		auth, err := uowctx.GetAuth(ctx)
		require.Error(t, err)
		require.Nil(t, auth)
		return nil
	})
	require.NoError(t, err)

	// 2. SetAuth
	auth := models.Auth{
		PasswordHash: []byte("test_hash_1"),
	}

	err = s.PasswordStorage().DoUoW(ctx, func(uowctx password.UoWContext) error {
		return uowctx.SetAuth(ctx, auth)
	})
	require.NoError(t, err)

	// 3. GetAuth
	err = s.PasswordStorage().DoUoW(ctx, func(uowctx password.UoWContext) error {
		got, err := uowctx.GetAuth(ctx)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, auth.PasswordHash, got.PasswordHash)
		return nil
	})
	require.NoError(t, err)

	// 4. Update password
	updatedAuth := models.Auth{
		PasswordHash: []byte("test_hash_2"),
	}

	err = s.PasswordStorage().DoUoW(ctx, func(uowctx password.UoWContext) error {
		return uowctx.SetAuth(ctx, updatedAuth)
	})
	require.NoError(t, err)

	// 5. check updated password
	err = s.PasswordStorage().DoUoW(ctx, func(uowctx password.UoWContext) error {
		got, err := uowctx.GetAuth(ctx)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, updatedAuth.PasswordHash, got.PasswordHash)
		return nil
	})
	require.NoError(t, err)
}*/
