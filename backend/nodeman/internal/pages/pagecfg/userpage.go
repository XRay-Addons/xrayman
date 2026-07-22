package pagecfg

/*
	type UserRoutes struct {
		ApiPrefix  string
		UserPrefix string
	}
*/
type AppLink struct {
	Name      string
	Platforms string
	URL       string
}

type UserPageCfg struct {
	ApiPrefix   string
	UserPrefix  string
	SupportLink string
	AppLinks    []AppLink
}
