package xrayapi

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/XRay-Addons/xrayman/node/internal/xrayctl/launchctl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	testExecPath = "/usr/local/bin/xrayman/xray"
)

const testXRayCfg = `{
  "log": { "loglevel": "warning" },

  "api": {
    "tag": "api",
    "listen": "127.0.0.1:32999",
    "services": ["HandlerService", "LoggerService", "StatsService", "ReflectionService"]
  },

  "inbounds": [
    {
	  "tag": "vlesstcp-reality",
      "port": 443,
      "protocol": "vless",
      "settings": {
        "clients": [],
        "decryption": "none"
      },
      "streamSettings": {
        "network": "tcp",
        "security": "reality",
        "realitySettings": {
          "show": false,
          "dest": "www.cloudflare.com:443",
          "xver": 0,
          "serverNames": ["www.cloudflare.com"],
          "privateKey": "4BHzOYgdeeG4de3oFimrg865ky_5X9cVoxLc_VmtEHc",
          "shortIds": [""]
        }
      }
    }  
  ],
  "outbounds": [{ "protocol": "freedom", "settings": {} }]
}`

const testApiURL = "127.0.0.1:32999"

var testXRayUser = models.User{
	Name:      "username",
	VlessUUID: "aaaabbbbccccdddd",
}

var testXRayInbounds = []models.Inbound{
	{Tag: "vlesstcp-reality", Type: models.VlessTcpReality},
}

// test service ctl
func TestXRayAPI(t *testing.T) {
	ctx := context.TODO()

	testCfgPath := filepath.Join(t.TempDir(), "config.json")

	// write xray config to file,
	// remove it after execution
	err := os.WriteFile(testCfgPath, []byte(testXRayCfg), 0o644)
	require.NoError(t, err)
	defer func() {
		err := os.Remove(testCfgPath)
		assert.NoError(t, err)
	}()

	// create xray service
	log := zap.NewNop()
	xrayctl, err := launchctl.New(testExecPath, testCfgPath, log)
	require.NoError(t, err)
	defer func() {
		err := xrayctl.Close(ctx)
		assert.NoError(t, err)
	}()

	// wait 1 second till service initialization
	time.Sleep(1 * time.Second)

	// start node
	err = xrayctl.Start(ctx)
	assert.NoError(t, err)

	// wait 1 second till service started
	time.Sleep(1 * time.Second)

	// check status is running now
	status, err := xrayctl.Status(ctx)
	assert.NoError(t, err)
	assert.Equal(t, models.ServiceRunning, status)

	// create xray api
	xrayapi, err := New(testApiURL)
	assert.NoError(t, err)

	// ping xray service
	err = xrayapi.Ping(ctx)
	assert.NoError(t, err)

	// edit users
	err = xrayapi.EditUsers(ctx,
		[]models.User{testXRayUser},
		[]models.User{},
		testXRayInbounds,
	)
	assert.NoError(t, err)
}
