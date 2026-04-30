package config

type Config struct {
	Endpoint      string `env:"ENDPOINT"`
	DBConn        string `env:"DBCONN"`
	certsDir      string `env:"CERTS_DIR"`
	UserSPAPrefix string `env:"USER_SPA_PREFIX"`
	APIPrefix     string `env:"API_PREFIX"`
	AdminPassword string `env:"ADMIN_PASSWORD"`
	nodemanCrt    string
	nodemanKey    string
	rootCrt       string
}

func (c *Config) HasCerts() bool {
	return c.nodemanCrt != "" ||
		c.nodemanKey != "" ||
		c.rootCrt != ""
}

func (c *Config) NodemanCrt() string {
	return c.nodemanCrt
}

func (c *Config) NodemanKey() string {
	return c.nodemanKey
}

func (c *Config) RootCrt() string {
	return c.rootCrt
}
