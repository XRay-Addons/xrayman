package subheaders

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type option = func(s *Service)

func WithLogger(l *zap.Logger) option {
	return func(s *Service) {
		if l != nil {
			s.log = l
		}
	}
}

type Service struct {
	storage Storage
	log     *zap.Logger
}

var _ handler.SubHeadersService = (*Service)(nil)

func New(storage Storage, opts ...option) (*Service, error) {
	if storage == nil {
		return nil, errdefs.NilArg("storage")
	}
	s := &Service{
		storage: storage,
		log:     zap.NewNop(),
	}
	for _, o := range opts {
		o(s)
	}
	return s, nil
}

// NewHeader implements handler.SubscrService.
func (s *Service) NewHeader(ctx context.Context,
	p models.NewSubHeaderParams,
) (*models.Header, error) {
	if s == nil || s.storage == nil {
		return nil, errdefs.NilCall()
	}

	var header models.Header
	header.Key = p.Key
	header.Value = p.Value
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.NewSubHeader(ctx, &header)
		return
	}); err != nil {
		return nil, err
	}

	return &header, nil
}

func (s *Service) ListHeaders(ctx context.Context,
	p models.ListSubHeadersParams,
) (*models.ListSubHeadersResult, error) {
	if s == nil || s.storage == nil {
		return nil, errdefs.NilCall()
	}

	var headers []models.Header
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		headers, err = uowctx.ListSubHeaders(ctx)
		return
	}); err != nil {
		return nil, err
	}

	return &models.ListSubHeadersResult{
		Headers: headers,
	}, nil
}

func (s *Service) DeleteHeader(ctx context.Context,
	p models.DeleteSubHeaderParams,
) (*models.DeleteSubHeaderResult, error) {
	if s == nil || s.storage == nil {
		return nil, errdefs.NilCall()
	}

	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.DeleteSubHeader(ctx, p.ID)
		return
	}); err != nil {
		return nil, err
	}

	return &models.DeleteSubHeaderResult{}, nil
}
