package config

import (
	"net/url"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Endpoint      string
	DBConn        string
	AdminPassword string
	JwtSecret     string

	NodeCallTimeout    time.Duration
	StorageCallTimeout time.Duration

	ApiServicePath string
	UserSpaPath    string
	AdminSpaPath   string

	ApiServiceUrl string
	UserSpaUrl    string
	AdminSpaUrl   string

	StateSyncInterval time.Duration
	StatsSyncInterval time.Duration

	AllowedOrigins  []string
	MetricsEndpoint string
	LogLevel        zapcore.Level
}

const (
	apiServicePath = "/api"
	userSpaPath    = "/u"
	adminSpaPath   = "/adm"
)

func NewConfig(cli *CLI) (*Config, error) {
	cfg, err := loadConfig(cli)
	if err != nil {
		return nil, xerr.WrapWithInfo(err, "load config")
	}

	if err := Validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadConfig(cli *CLI) (*Config, error) {

	cfg := Config{
		Endpoint:          cli.Endpoint,
		DBConn:            cli.DBConn,
		AdminPassword:     cli.AdminPassword,
		JwtSecret:         cli.JwtSecret,
		StateSyncInterval: time.Duration(cli.StateSyncInterval) * time.Second,
		StatsSyncInterval: time.Duration(cli.StatsSyncInterval) * time.Second,

		ApiServicePath: apiServicePath,
		UserSpaPath:    userSpaPath,
		AdminSpaPath:   adminSpaPath,

		NodeCallTimeout:    time.Duration(cli.NodeCallTimeout) * time.Second,
		StorageCallTimeout: time.Duration(cli.StorageCallTimeout) * time.Second,

		MetricsEndpoint: cli.MetricsEndpoint,
		LogLevel:        cli.LogLevel,
	}

	cfg.ApiServiceUrl = or(cli.ApiServiceUrl, cfg.ApiServicePath)
	cfg.UserSpaUrl = or(cli.UserSpaUrl, cfg.UserSpaPath)
	cfg.AdminSpaUrl = or(cli.AdminSpaUrl, cfg.AdminSpaPath)

	for _, u := range []string{cli.ApiServiceUrl, cli.UserSpaUrl, cli.AdminSpaUrl} {
		o, err := getUrlOrigin(u)
		if err != nil {
			return nil, err
		}
		if o != "" {
			cfg.AllowedOrigins = append(cfg.AllowedOrigins, o)
		}
	}

	return &cfg, nil
}

func or(a string, b string) string {
	if a != "" {
		return a
	}
	return b
}

// for empty or relative return ""
// else return origin
func getUrlOrigin(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", xerr.WrapWithStack(err)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return "", nil
	}

	return parsed.Scheme + "://" + parsed.Host, nil
}
