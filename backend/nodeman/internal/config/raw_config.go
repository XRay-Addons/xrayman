package config

type RawConfig struct {
	Endpoint      string `env:"ENDPOINT"`
	DBConn        string `env:"DBCONN"`
	AdminPassword string `env:"ADMIN_PASSWORD"`
	JwtSecret     string `env:"JWT_SECRET"`

	ApiServiceUrl string `env:"API_SERVICE_URL"`
	UserSpaUrl    string `env:"USER_SPA_URL"`
	AdminSpaUrl   string `env:"ADMIN_SPA_URL"`
}
