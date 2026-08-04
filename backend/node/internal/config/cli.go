package config

import (
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/alecthomas/kong"
	"go.uber.org/zap/zapcore"
)

var kongVars = kong.Vars{
	"endpointHelp": "server endpoint tcp address, like :8080, 127.0.0.1:80, localhost:22",

	"xrayDataHelp": "xray data dir, should contains geoip, geodat if routing uses it",

	"xrayConfigHelp": `xray configs dir, must contains xray_server.json and xray_client.json.
xray_server.json should be valid xray server config,
xray_clinet.json should be user config template,
supported template params:
  - {{ .VlessEmail }}
  - {{ .VlessUUID }}`,

	"persistentHelp": `persistent config dir. persistent objects
(certs, secrets, config to connect to node)
should be generated on-demand`,
}

type CLI struct {
	Version       bool          `short:"v" help:"Show version and exit."`
	Endpoint      string        `short:"a" default:"localhost:8080" env:"ENDPOINT" help:"${endpointHelp}"`
	XRayDataDir   string        `short:"d" env:"XRAY_DATA_DIR" help:"${xrayDataHelp}"`
	XRayConfigDir string        `short:"c" env:"XRAY_CONFIG_DIR" help:"${xrayConfigHelp}"`
	PersistentDir string        `short:"p" env:"PERSISTENT_DIR" help:"${persistentHelp}"`
	LogLevel      zapcore.Level `name:"log-lvl" default:"info" env:"LOG_LEVEL" help:"zap log level"`
}

func LoadCLI() (*CLI, error) {
	var cli CLI

	ctx := kong.Parse(&cli,
		kongVars,
		kong.Name("xray-node"),
		kong.Description("XRay node daemon"),
	)

	if err := ctx.Validate(); err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	return &cli, nil
}
