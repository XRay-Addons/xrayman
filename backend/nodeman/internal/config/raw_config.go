package config

type RawConfig struct {
	Endpoint      string `env:"ENDPOINT"`
	DBConn        string `env:"DBCONN"`
	AdminPassword string `env:"ADMIN_PASSWORD"`
	JwtSecret     string `env:"JWT_SECRET"`

	ApiServiceUrl string `env:"API_SERVICE_URL"`
	UserSpaUrl    string `env:"USER_SPA_URL"`
	AdminSpaUrl   string `env:"ADMIN_SPA_URL"`

	StateSyncInterval int `env:"STATE_SYNC_INTERVAL"`
	StatsSyncInterval int `env:"STATS_SYNC_INTERVAL"`

	NodeCallTimeout    int `env:"NODE_CALL_TIMEOUT"`
	StorageCallTimeout int `env:"STORAGE_CALL_TIMEOUT"`

	LogLevel string `env:"LOG_LEVEL"`
}
