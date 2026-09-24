package config

import (
	"strings"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/alecthomas/kong"
	"go.uber.org/zap/zapcore"
)

var kongVars = kong.Vars{
	"endpointHelp": "Nodeman server endpoint tcp address, like :8080, 127.0.0.1:80, localhost:22",

	"dbHelp": "Postgress connection string, like postgresql://user:password@127.0.0.1:4321/dbname",

	"jwtHelp": "JWT secret",

	"stateHelp": "Nodes state check and sync interval, every N s",

	"statsHelp": "Nodes statistics refresh interval, eveny N s",

	"apisrvHelp": `Public-facing base URL of nodeman API, as seen by clients
(e.g. behind a reverse proxy / TLS termination).
Used when the server needs to build a full external
URL to itself (e.g. https://api.example.com/api),
as opposed to --endpoint, which is the internal
address the process actually listens on. 
Should be short like /api or full https://api.example.com/apipath (optional)`,

	"userspaHelp": `Public-facing base URL of user page, as seen by clients
(e.g. behind a reverse proxy / TLS termination).
Used when the server needs to build a full external
URL to itself (e.g. https://u.example.com/user),
Should be short like /u or full https://u.example.com/user (optional)`,

	"adminspaHelp": `Public-facing base URL of admin page, as seen by clients
(e.g. behind a reverse proxy / TLS termination).
Used when the server needs to build a full external
URL to itself (e.g. https://a.example.com/adm),
Should be short like /adm or full https://adm.example.com/admin (optional)`,

	"storageTimeoutHelp": `Storage call timeout, s (optional)`,

	"nodeTimeoutHelp": `Node call timeout, s (optional)`,

	"metricsHelp": `Prometheus metrics endpoint (optional)`,
}

type CLI struct {
	DBConn   string `name:"db" help:"${dbHelp}"`
	Endpoint string `name:"endpoint" default:"localhost:80" help:"${endpointHelp}"`

	JwtSecret     string `name:"jwt" help:"${jwtHelp}"`
	AdminPassword string `name:"adm-pass" default:"" help:"${admpassHelp}"`

	ApiServiceUrl string `name:"api-service-url" default:"" help:"${apisrvHelp}"`
	UserSpaUrl    string `name:"user-spa-url" default:"" help:"${userspaHelp}"`
	AdminSpaUrl   string `name:"admin-spa-url" default:"" help:"${adminspaHelp}"`

	StateSyncInterval int `name:"state-sync-interval" default:"5" help:"${stateHelp}"`
	StatsSyncInterval int `name:"stats-sync-interval" default:"60" help:"${statsHelp}"`

	NodeCallTimeout    int `name:"node-call-timeout" default:"5" help:"${nodeTimeoutHelp}"`
	StorageCallTimeout int `name:"storage-call-timeout" default:"5" help:"${storageTimeoutHelp}"`

	MetricsEndpoint string        `name:"metrics-endpoint" default:"" help:"${metricsHelp}"`
	LogLevel        zapcore.Level `name:"log-level" default:"info" help:"Zap log level"`

	Version bool `short:"v" help:"Show version and exit."`
}

func LoadCLI() (*CLI, error) {
	var cli CLI

	ctx := kong.Parse(&cli,
		kongVars,
		kong.Name("xray-nodeman"),
		kong.Description("XRay node manager daemon"),
		kong.DefaultEnvars("xray_nodeman"),
	)

	if err := ctx.Validate(); err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	return &cli, nil
}

var _ zapcore.ObjectMarshaler = (*CLI)(nil)

func (c *CLI) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("DBConn", c.DBConn)
	enc.AddString("Endpoint", c.Endpoint)

	if c.JwtSecret != "" {
		enc.AddString("JwtSecret", strings.Repeat("*", len(c.JwtSecret)))
	}
	if c.AdminPassword != "" {
		enc.AddString("AdminPassword", strings.Repeat("*", len(c.AdminPassword)))
	}

	enc.AddString("ApiServiceUrl", c.ApiServiceUrl)
	enc.AddString("UserSpaUrl", c.UserSpaUrl)
	enc.AddString("AdminSpaUrl", c.AdminSpaUrl)

	enc.AddString("StateSyncInterval", (time.Duration(c.StateSyncInterval) * time.Second).String())
	enc.AddString("StatsSyncInterval", (time.Duration(c.StatsSyncInterval) * time.Second).String())

	enc.AddString("NodeCallTimeout", (time.Duration(c.NodeCallTimeout) * time.Second).String())
	enc.AddString("StorageCallTimeout", (time.Duration(c.StorageCallTimeout) * time.Second).String())

	enc.AddString("MetricsEndpoint", c.MetricsEndpoint)

	enc.AddString("LogLevel", c.LogLevel.String())

	return nil
}
