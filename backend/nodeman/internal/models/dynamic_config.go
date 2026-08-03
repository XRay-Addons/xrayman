package models

type AppLink struct {
	Name      string
	Platforms string
	URL       string
}

type Settings struct {
	SubscrTitle    string
	UpdateInterval int
	UserPage       string
	UsersMessage   string
	TgPage         string
	Routing        string

	AppLinks      []AppLink
	CustomHeaders []SubHeader
}
