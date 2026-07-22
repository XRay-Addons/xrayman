package models

type DynamicConfigHeader struct {
	Key   string
	Value string
}

type AppLink struct {
	Name      string
	Platforms string
	URL       string
}

type DynamicConfig struct {
	SubscrTitle    string
	UpdateInterval int
	UserPage       string
	UsersMessage   string
	TgPage         string
	Routing        string

	AppLinks []AppLink
	Headers  []DynamicConfigHeader
}
