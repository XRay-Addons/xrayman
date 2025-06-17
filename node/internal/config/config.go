package config

type Config struct {
	XRayExecPath         string `env:"XRAY"`
	XRayServerConfigPath string `env:"XRAYSERVERCFG"`
	XRayClientConfigPath string `env:"XRAYCLIENTCFG"`
	Endpoint             string `env:"ENDPOINT"`
}
