package subscr

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/common/xerrgroup"
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

var _ handler.SubscrService = (*Service)(nil)

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

func (s *Service) GetUserSub(ctx context.Context,
	p models.UserSubParams,
) (*models.UserSubResult, error) {
	if s == nil || s.storage == nil {
		return nil, errdefs.NilCall()
	}

	g, ctx := xerrgroup.WithContext(ctx)
	// find user
	var user *models.UserView
	g.Go(func() (err error) {
		user, err = s.storage.GetUserView(ctx, p.ID, p.Name)
		return
	})

	// get active nodes for user
	var userNodes []models.Node
	g.Go(func() (err error) {
		userNodes, err = s.storage.GetUserNodes(ctx, p.ID)
		return
	})

	// get dynamic config
	var dynConfig *models.DynamicConfig
	g.Go(func() (err error) {
		dynConfig, err = s.storage.GetDynamicConfig(ctx)
		return
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// reply NotFound for disabled users
	if user.User.TargetStatus != models.UserStatusEnabled {
		return nil, xerr.WrapWithStack(errdefs.ErrNotFound)
	}

	// get subscription content
	clientCfgs := createClientCfgs(user, userNodes, s.log)

	// get subscription headers
	clientHeaders := createClientHeaders(ctx, user, dynConfig)

	return &models.UserSubResult{
		Headers:       clientHeaders,
		ClientConfigs: clientCfgs,
	}, nil
}

func (s Service) SubHeadersPlaceholders() []string {
	return listPlaceholders()
}
