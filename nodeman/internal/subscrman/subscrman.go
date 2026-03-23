package subscrman

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/jxext"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/template"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/go-faster/jx"
	"go.uber.org/zap"
)

type SubscrMan interface {
	GetUserSub(ctx context.Context, user models.User) (*models.UserSubResult, error)
}

type option = func(s *subscrMan)

func WithLog(l *zap.Logger) option {
	return func(s *subscrMan) {
		if l != nil {
			s.log = l
		}
	}
}

type subscrMan struct {
	storage Storage
	log     *zap.Logger
}

func New(storage Storage, opts ...option) (SubscrMan, error) {
	if storage == nil {
		return nil, errdefs.NewNilArg("storage")
	}
	s := &subscrMan{
		storage: storage,
		log:     zap.NewNop(),
	}
	for _, o := range opts {
		o(s)
	}
	return s, nil
}

func (m *subscrMan) GetUserSub(ctx context.Context,
	user models.User,
) (*models.UserSubResult, error) {
	if m == nil || m.storage == nil {
		return nil, errdefs.NewNilCall()
	}

	// get active nodes for user
	var userNodes []models.Node
	if err := m.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		userNodes, err = uowctx.GetUserNodes(ctx, user.Profile.ID)
		return
	}); err != nil {
		return nil, err
	}

	// get subscription content
	clientCfgs, err := m.makeClientConfigs(user, userNodes)
	if err != nil {
		return nil, err
	}

	return &models.UserSubResult{
		ClientConfigs: clientCfgs,
	}, nil
}

func (m *subscrMan) makeClientConfigs(user models.User,
	userNodes []models.Node,
) ([]models.ClientConfigItem, error) {
	var clientCfgs []models.ClientConfigItem
	for _, node := range userNodes {
		nodeClientConfigs, err := m.makeNodeClientConfigs(
			user, node.Config.ClientConfigTemplate)
		if err != nil {
			// skip invalid node configs
			m.log.Warn("node client config", zap.Error(err))
			continue
		}
		clientCfgs = append(clientCfgs, nodeClientConfigs...)
	}

	return clientCfgs, nil
}

func (m *subscrMan) makeNodeClientConfigs(user models.User,
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
		if err = jxext.Validate(nodeConfig); err != nil {
			return nil, err
		}
		nodeConfigs = append(nodeConfigs, nodeConfig)
	}
	return nodeConfigs, nil
}
