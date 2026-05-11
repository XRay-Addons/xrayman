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
) (*models.UserSubResult, bool, error) {
	if s == nil || s.storage == nil {
		return nil, false, errdefs.NilCall()
	}

	// find user
	user, exists, err := s.findUser(ctx, p)
	if err != nil || !exists {
		return nil, exists, err
	}

	// get active nodes for user
	var userNodes []models.Node
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		userNodes, err = uowctx.GetUserNodes(ctx, user.Profile.ID)
		return
	}); err != nil {
		return nil, false, err
	}

	// get subscription content
	clientCfgs, err := s.makeClientConfigs(*user, userNodes)
	if err != nil {
		return nil, false, err
	}

	// get subscription headers
	var rawHeaders []models.Header
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		rawHeaders, err = uowctx.GetGlobalHeaders(ctx)
		return
	}); err != nil {
		return nil, false, err
	}
	headers := s.makeSubscrHeaders(rawHeaders)
	if err != nil {
		return nil, false, err
	}

	return &models.UserSubResult{
		Headers:       headers,
		ClientConfigs: clientCfgs,
	}, true, nil
}

func (s *Service) SetGlobalHeaders(ctx context.Context, p models.Headers) (*models.SetGlobalHeadersResult, error) {
	if s == nil || s.storage == nil {
		return nil, errdefs.NilCall()
	}

	rawHeaders := s.makeRawHeaders(p)

	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.SetGlobalHeaders(ctx, rawHeaders)
		return
	}); err != nil {
		return nil, err
	}

	return &models.SetGlobalHeadersResult{}, nil
}

func (s *Service) findUser(ctx context.Context, p models.UserSubParams) (*models.User, bool, error) {
	// find user with given id
	var user *models.User
	var exists bool
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		user, exists, err = uowctx.GetUser(ctx, p.ID)
		return
	}); err != nil {
		return nil, false, err
	}

	// check user name
	if !exists || user.Profile.Name != p.Name {
		return nil, false, nil
	}

	return user, true, nil
}

func (s *Service) makeClientConfigs(user models.User,
	userNodes []models.Node,
) ([]models.ClientConfigItem, error) {
	var clientCfgs []models.ClientConfigItem
	for _, node := range userNodes {
		nodeClientConfigs, err := s.makeNodeClientConfigs(
			user, node.Config.ClientConfigTemplate)
		if err != nil {
			// skip invalid node configs
			s.log.Warn("node client config", zap.Error(err))
			continue
		}
		clientCfgs = append(clientCfgs, nodeClientConfigs...)
	}

	return clientCfgs, nil
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

func (s *Service) makeSubscrHeaders(raw []models.Header) models.Headers {
	rawHeaders := make(map[string]string)
	for _, h := range raw {
		rawHeaders[h.Key] = h.Value
	}
	headers := models.Headers{}
	if h, ok := rawHeaders[RoutingHeader]; ok {
		headers.Routing = &h
	}
	if h, ok := rawHeaders[ProfileTitleHeader]; ok {
		headers.ProfileTitle = &h
	}
	return headers
}

func (s *Service) makeRawHeaders(headers models.Headers) []models.Header {
	rawHeaders := make([]models.Header, 0, 2)
	if headers.Routing != nil {
		rawHeaders = append(rawHeaders, models.Header{
			Key:   RoutingHeader,
			Value: *headers.Routing,
		})
	}
	if headers.ProfileTitle != nil {
		rawHeaders = append(rawHeaders, models.Header{
			Key:   ProfileTitleHeader,
			Value: *headers.ProfileTitle,
		})
	}

	return rawHeaders
}
