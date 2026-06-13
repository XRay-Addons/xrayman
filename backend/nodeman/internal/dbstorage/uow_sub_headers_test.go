// AI generated tests
package dbstorage

/*import (
	"context"
	"testing"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subheaders"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestSubHeaders(t *testing.T) {
	logger := zaptest.NewLogger(t)

	s, _, _, cleanup := setupTestDB(t, logger)
	defer cleanup()
	logger.Info("new test db inited")

	ctx := context.Background()

	h1 := models.Header{
		Key:   "X-Test-Header-1",
		Value: "value-1",
	}
	h2 := models.Header{
		Key:   "X-Test-Header-2",
		Value: "value-2",
	}

	// 1. add headers
	err := s.SubHeadersStorage().DoUoW(ctx, func(uowctx subheaders.UoWContext) error {
		if err := uowctx.NewSubHeader(ctx, &h1); err != nil {
			return err
		}
		if err := uowctx.NewSubHeader(ctx, &h2); err != nil {
			return err
		}
		return nil
	})
	require.NoError(t, err)

	require.Equal(t, models.HeaderID(1), h1.ID)
	require.Equal(t, models.HeaderID(2), h2.ID)

	// 2. list headers
	var headers []models.Header
	err = s.SubHeadersStorage().DoUoW(ctx, func(uowctx subheaders.UoWContext) error {
		var err error
		headers, err = uowctx.ListSubHeaders(ctx)
		return err
	})
	require.NoError(t, err)
	require.Len(t, headers, 2)
	require.Equal(t, h1, headers[0])
	require.Equal(t, h2, headers[1])

	// 3. upsert by key (same key, new value)
	updated := models.Header{
		Key:   h1.Key,
		Value: "value-1-updated",
	}

	err = s.SubHeadersStorage().DoUoW(ctx, func(uowctx subheaders.UoWContext) error {
		return uowctx.NewSubHeader(ctx, &updated)
	})
	require.NoError(t, err)

	require.Equal(t, h1.ID, updated.ID)

	// 4. list again
	err = s.SubHeadersStorage().DoUoW(ctx, func(uowctx subheaders.UoWContext) error {
		var err error
		headers, err = uowctx.ListSubHeaders(ctx)
		return err
	})
	require.NoError(t, err)
	require.Len(t, headers, 2)
	require.Equal(t, "value-1-updated", headers[0].Value)

	// 5. delete header
	err = s.SubHeadersStorage().DoUoW(ctx, func(uowctx subheaders.UoWContext) error {
		return uowctx.DeleteSubHeader(ctx, h1.ID)
	})
	require.NoError(t, err)

	// 6. list after delete
	err = s.SubHeadersStorage().DoUoW(ctx, func(uowctx subheaders.UoWContext) error {
		var err error
		headers, err = uowctx.ListSubHeaders(ctx)
		return err
	})
	require.NoError(t, err)
	require.Len(t, headers, 1)
	require.Equal(t, h2, headers[0])
}*/
