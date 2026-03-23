package dbstorage

const (
	UsersTable = "users"

	UserIDCol           = "user_id"
	DisplayNameCol      = "display_name"
	UserNameCol         = "user_name"
	VlessUUIDCol        = "vless_uuid"
	UserTargetStatusCol = "user_target_status"
)

const (
	NodesTable = "nodes"

	NodeIDCol               = "node_id"
	ClientConfigTemplateCol = "client_cfg_template"
	NodeEndpointCol         = "node_endpoint"
	NodeAccessKeyCol        = "node_access_key"
	NodeCurrentStatusCol    = "node_current_status"
	NodeTargetStatusCol     = "node_target_status"
)

const (
	SyncsTable = "syncs"

	UserCurrentStatusCol = "user_current_status"
)
