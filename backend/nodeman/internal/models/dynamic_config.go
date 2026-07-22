package models

type DynamicConfigHeader struct {
	Key   string
	Value string
}

type PlatformApp struct {
	Name      string
	Platforms string
	Url       string
}

type DynamicConfig struct {
	SubscrTitle    string
	UpdateInterval int
	UserPage       string
	UsersMessage   string
	TgPage         string
	Routing        string

	PlatformApps []PlatformApp
	Headers      []DynamicConfigHeader
}
