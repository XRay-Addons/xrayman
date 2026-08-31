package config

import (
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/alecthomas/kong"
	"go.uber.org/zap/zapcore"
)

var kongVars = kong.Vars{
	"endpointHelp": "server endpoint tcp address, like :8080, 127.0.0.1:80, localhost:22",

	"dbHelp": "postgress connection string, like postgresql://user@password/127.0.0.1:4321/dbname",

	"jwtHelp": "jwt secret",

	"stateHelp": "state sync interval, s",

	"statsHelp": "stats sync interval, s",

	"apisrvHelp": `public base URL of the API as seen by browsers (used for CORS and SPAs config).
If empty or relative, the internal API base path is used.
should be like /internal/api or https://api.example.com/api (optional)`,

	"userspaHelp": `public base URL of the User SPA as seen by browsers (used for CORS and SPAs config).
If empty or relative, the internal User SPA base path is used.
should be like /user or https://u.example.com (optional)`,

	"adminspaHelp": `public base URL of the Admin SPA as seen by browsers (used for CORS and SPAs config).
If empty or relative, the internal Admin SPA base path is used.
should be like /admin or https://adm.example.com (optional)`,

	"admpassHelp": `admin password to change (optional, empty for keep prev pwd)`,

	"storageTimeoutHelp": `storage call timeout, s (optional)`,

	"nodeTimeoutHelp": `node call timeout, s (optional)`,

	"metricsHelp": `prometheus metrics endpoint (optional)`,
}

type CLI struct {
	DBConn    string `name:"db" help:"${dbHelp}"`
	JwtSecret string `name:"jwt" help:"${jwtHelp}"`

	Endpoint      string `name:"endpoint" default:"localhost:80" help:"${endpointHelp}"`
	AdminPassword string `name:"adm-pass" default:"" help:"${admpassHelp}"`

	ApiServiceUrl string `name:"api-service-url" default:"" help:"${apisrvHelp}"`
	UserSpaUrl    string `name:"user-spa-url" default:"" help:"${userspaHelp}"`
	AdminSpaUrl   string `name:"admin-spa-url" default:"" help:"${adminspaHelp}"`

	StateSyncInterval int `name:"state-sync-interval" default:"5" help:"${stateHelp}"`
	StatsSyncInterval int `name:"stats-sync-interval" default:"60" help:"${statsHelp}"`

	NodeCallTimeout    int `name:"node-call-timeout" default:"5" help:"${nodeTimeoutHelp}"`
	StorageCallTimeout int `name:"storage-call-timeout" default:"5" help:"${storageTimeoutHelp}"`

	MetricsEndpoint string        `name:"metrics-endpoint" default:"" help:"${metricsHelp}"`
	LogLevel        zapcore.Level `name:"log-level" default:"info" help:"zap log level"`

	Version bool `short:"v" help:"Show version and exit."`
}

func LoadCLI() (*CLI, error) {
	var cli CLI

	ctx := kong.Parse(&cli,
		kongVars,
		kong.Name("xray-node"),
		kong.Description("XRay node daemon"),
		kong.DefaultEnvars("xray_nodeman"),
	)

	if err := ctx.Validate(); err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	return &cli, nil
}
