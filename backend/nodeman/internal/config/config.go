package config

type Config struct {
	Endpoint       string `env:"ENDPOINT"`
	DBConn         string `env:"DBCONN"`
	UserSpaPrefix  string `env:"USER_SPA_PREFIX"`
	AdminSpaPrefix string `env:"ADMIN_SPA_PREFIX"`
	APIPrefix      string `env:"API_PREFIX"`
	AdminPassword  string `env:"ADMIN_PASSWORD"`
	JWTSecret      string `env:"JWT_SECRET"`
}
