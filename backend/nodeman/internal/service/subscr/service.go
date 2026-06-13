package subscr

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/jsonval"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/template"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/go-faster/jx"
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

	// find user
	user, err := s.findUser(ctx, p)
	if err != nil {
		return nil, err
	}

	// get active nodes for user
	var userNodes []models.Node
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		userNodes, err = uowctx.GetUserNodes(ctx, user.User.Profile.ID)
		return
	}); err != nil {
		return nil, err
	}

	// get subscription content
	clientCfgs := s.makeClientConfigs(*user, userNodes)

	// get subscription headers
	clientHeaders, err := s.makeClientHeaders(ctx, *user)
	if err != nil {
		return nil, err
	}

	return &models.UserSubResult{
		Headers:       clientHeaders,
		ClientConfigs: clientCfgs,
	}, nil
}

func (s *Service) findUser(ctx context.Context, p models.UserSubParams) (*models.UserView, error) {
	// find user with given id
	var user *models.UserView
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		user, err = uowctx.GetUserView(ctx, p.ID, p.Name)
		return
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) makeClientConfigs(user models.UserView,
	userNodes []models.Node,
) []models.ClientConfigItem {
	var clientCfgs []models.ClientConfigItem
	for _, node := range userNodes {
		nodeClientConfigs, err := s.makeNodeClientConfigs(
			user.User, node.Config.ClientConfigTemplate)
		if err != nil {
			// skip invalid node configs
			s.log.Warn("node client config", zap.Error(err))
			continue
		}
		clientCfgs = append(clientCfgs, nodeClientConfigs...)
	}

	return clientCfgs
}

func (s *Service) makeNodeClientConfigs(user models.User,
	cfgTemplate models.ClientConfigTemplate,
) ([]models.ClientConfigItem, error) {
	nodeConfigs := make([]models.ClientConfigItem, 0, len(cfgTemplate.Template))
	for _, item := range cfgTemplate.Template {
		tmpl, err := template.RenderTemplate(item.String(), map[string]string{
			cfgTemplate.VlessEmailField: user.Profile.VlessEmail(),
			cfgTemplate.VlessUUIDField:  user.Profile.VlessUUID,
		})
		if err != nil {
			return nil, err
		}
		nodeConfig := jx.Raw(tmpl)
		if err = jsonval.ValidateJsonData(nodeConfig); err != nil {
			return nil, err
		}
		nodeConfigs = append(nodeConfigs, nodeConfig)
	}
	return nodeConfigs, nil
}

func (s *Service) makeClientHeaders(ctx context.Context,
	u models.UserView,
) (models.Headers, error) {
	var headers models.Headers
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		headers, err = uowctx.ListSubHeaders(ctx)
		return
	}); err != nil {
		return nil, err
	}
	headers = replacePlaceholders(headers, u)
	return headers, nil
}

func (s Service) SubHeadersPlaceholders() []string {
	return listPlaceholders()
}
