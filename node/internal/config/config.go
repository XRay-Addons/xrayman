package config

type Config struct {
	Endpoint  string `env:"ENDPOINT"`
	AccessKey string `env:"ACCESS_KEY"`

	XRayExecPath   string `env:"XRAY_EXEC"`
	XRayConfigPath string `env:"XRAY_CFG"`
}
