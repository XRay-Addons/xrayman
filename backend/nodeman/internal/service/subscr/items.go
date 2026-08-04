package subscr

import (
	"github.com/XRay-Addons/xrayman/common/jsonval"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/template"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/go-faster/jx"
	"go.uber.org/zap"
)

func createClientCfgs(user *models.UserView, userNodes []models.Node,
	log *zap.Logger,
) []models.ClientConfigItem {
	var clientCfgs []models.ClientConfigItem
	for _, node := range userNodes {
		nodeClientConfigs, err := createNodeClientCfgs(
			user.User, node.Config.Settings.ClientConfigTemplate)
		if err != nil {
			// skip invalid node configs
			log.Warn("node client config", zap.Error(err))
			continue
		}
		clientCfgs = append(clientCfgs, nodeClientConfigs...)
	}

	return clientCfgs
}

func createNodeClientCfgs(user models.User,
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
