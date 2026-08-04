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
}

type CLI struct {
	Version bool `short:"v" help:"Show version and exit."`

	Endpoint      string `short:"a" env:"ENDPOINT" default:"localhost:80" help:"${endpointHelp}"`
	DBConn        string `name:"db" env:"DBCONN" help:"${dbHelp}"`
	AdminPassword string `name:"admpass" env:"ADMIN_PASSWORD" default:"" help:"${admpassHelp}"`
	JwtSecret     string `name:"jwt" env:"JWT_SECRET" help:"${jwtHelp}"`

	ApiServiceUrl string `name:"apisrv" env:"API_SERVICE_URL" help:"${apisrvHelp}"`
	UserSpaUrl    string `name:"userspa" env:"USER_SPA_URL" help:"${userspaHelp}"`
	AdminSpaUrl   string `name:"adminspa" env:"ADMIN_SPA_URL" help:"${adminspaHelp}"`

	StateSyncInterval int `name:"state" env:"STATE_SYNC_INTERVAL" defautl:"5" help:"${stateHelp}"`
	StatsSyncInterval int `name:"stats" env:"STATS_SYNC_INTERVAL" default:"60" help:"${statsHelp}"`

	NodeCallTimeout    int `name:"node-timeout" env:"NODE_CALL_TIMEOUT" default:"5" help:"${nodeTimeoutHelp}"`
	StorageCallTimeout int `name:"storage-timeout" env:"STORAGE_CALL_TIMEOUT" default:"5" help:"${storageTimeoutHelp}"`

	LogLevel zapcore.Level `name:"log-lvl" env:"LOG_LEVEL" default:"info" help:"zap log level"`
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
