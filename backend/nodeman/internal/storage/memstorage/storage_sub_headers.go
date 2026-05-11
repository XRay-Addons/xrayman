package memstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (s *Storage) NewSubHeader(ctx context.Context, header *models.Header) error {
	header.ID = models.HeaderID(len(s.headers))
	s.headers = append(s.headers, *header)
	return nil
}

func (s *Storage) DeleteSubHeader(ctx context.Context, id models.HeaderID) error {
	s.headers[id].Key = ""
	s.headers[id].Value = ""
	return nil
}

func (s *Storage) ListSubHeaders(ctx context.Context) ([]models.Header, error) {
	var headers []models.Header
	headers = append(headers, s.headers...)
	return headers, nil
}
