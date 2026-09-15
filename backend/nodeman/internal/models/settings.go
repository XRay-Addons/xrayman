package models

type AppLink struct {
	Name      string
	Platforms string
	URL       string
}

type Settings struct {
	RecentDays int

	UsersLanguage string

	SubscrTitle    string
	UpdateInterval int
	UserPage       string
	UsersMessage   string
	TgPage         string
	Routing        string

	AppLinks      []AppLink
	CustomHeaders []SubHeader
}
