package models

type UserConfigTemplate struct {
	Template      string
	UserIDField   string
	UserNameField string
	UserKeyField  string
}

type NodeConfig struct {
	UserConfigTemplate UserConfigTemplate
}
