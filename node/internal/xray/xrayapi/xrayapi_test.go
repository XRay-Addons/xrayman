package xrayapi

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/XRay-Addons/xrayman/node/internal/logging"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testExecPath = "/usr/local/bin/xrayman/xray"

	testXRayCfg = `{
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

	testApiURL = "127.0.0.1:32999"
)

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

	log, err := logging.New()
	require.NoError(t, err)

	// create xray api
	xrayapi, err := New(testApiURL, testXRayInbounds, log)
	assert.NoError(t, err)
	defer func() {
		xrayapi.Close(ctx)
	}()

	// ping stopped service
	err = xrayapi.Ping(ctx)
	assert.Error(t, err)

	// write xray config to file,
	// remove it after execution
	testCfgPath := filepath.Join(t.TempDir(), "config.json")
	err = os.WriteFile(testCfgPath, []byte(testXRayCfg), 0o644)
	require.NoError(t, err)
	defer func() {
		err := os.Remove(testCfgPath)
		assert.NoError(t, err)
	}()

	// run xray
	xrayCmd := exec.Command(testExecPath, "-config", testCfgPath)
	err = xrayCmd.Start()
	require.NoError(t, err)
	defer func() {
		err := xrayCmd.Process.Kill()
		require.NoError(t, err, "xray kill error")
	}()

	// connect to xray api
	err = xrayapi.Connect(ctx)
	require.NoError(t, err)

	// ping xray service
	err = xrayapi.Ping(ctx)
	assert.NoError(t, err)

	// edit users
	err = xrayapi.EditUsers(ctx,
		[]models.User{testXRayUser},
		[]models.User{},
	)
	assert.NoError(t, err)

	// edit users again (expecting no error)
	err = xrayapi.EditUsers(ctx,
		[]models.User{testXRayUser},
		[]models.User{},
	)
	assert.NoError(t, err)
}
