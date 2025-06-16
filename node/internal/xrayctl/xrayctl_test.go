package xrayctl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testExec   = "xray"
	testParams = "/var/etc/xray/conf.json"
)

var expectedCfg = `[Unit]
Description = xray service
After       = network.target

[Service]
Type      = simple
ExecStart = xray /var/etc/xray/conf.json
`

func TestCreateServiceConfig(t *testing.T) {
	cfg, err := createServiceConfig(testExec, testParams)
	require.NoError(t, err)
	require.Equal(t, expectedCfg, cfg)
}
