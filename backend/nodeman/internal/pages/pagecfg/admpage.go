package pagecfg

type AdminRoutes struct {
	ApiPrefix   string
	AdminPrefix string
	UserPrefix  string
}

type SubPlaceholders = []string

type AdminPageCfg struct {
	Routes                 AdminRoutes
	SubHeadersPlaceholders SubPlaceholders
}
