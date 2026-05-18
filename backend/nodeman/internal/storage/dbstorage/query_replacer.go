package dbstorage

import (
	"strings"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

var (
	queryReplacer = strings.NewReplacer(
		"{users}", UsersTable,
		"{user_id}", UserIDCol,
		"{display_name}", DisplayNameCol,
		"{user_name}", UserNameCol,
		"{vless_uuid}", VlessUUIDCol,
		"{user_target_status}", UserTargetStatusCol,

		"{nodes}", NodesTable,
		"{node_id}", NodeIDCol,
		"{client_config_template}", ClientConfigTemplateCol,
		"{node_endpoint}", NodeEndpointCol,
		"{node_access_key}", NodeAccessKeyCol,
		"{node_current_status}", NodeCurrentStatusCol,
		"{node_target_status}", NodeTargetStatusCol,

		"{syncs}", SyncsTable,
		"{user_current_status}", UserCurrentStatusCol,

		"{user_status_enabled}", models.UserStatusEnabled.StringInt(),
		"{user_status_disabled}", models.UserStatusDisabled.StringInt(),
		"{node_status_running}", models.NodeStatusRunning.StringInt(),

		"{admin_auth}", AdminAuthTable,
		"{admin_id}", AdminIdCol,
		"{password_hash}", PasswordHashCol,

		"{sub_headers}", SubHeadersTable,
		"{header_id}", HeaderIDCol,
		"{header_key}", HeaderKeyCol,
		"{header_value}", HeaderValueCol,

		"{created_at}", CreatedAtCol,
		"{updated_at}", UpdatedAtCol,
		"{deleted_at}", DeletedAtCol,
	)
)
