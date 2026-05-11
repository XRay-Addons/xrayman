package dbstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (uow *uowctx) GetGlobalHeaders(ctx context.Context) ([]models.Header, error) {
	headers := []models.Header{
		models.Header{Key: "key1", Value: "value1"},
		models.Header{Key: "key2", Value: "value2"},
	}
	return headers, nil
}

func (uow *uowctx) SetGlobalHeaders(ctx context.Context, headers []models.Header) error {
	/*headers := []models.Header{
		models.Header{Key: "key1", Value: "value1"},
		models.Header{Key: "key2", Value: "value2"},
	}*/
	return nil
}
