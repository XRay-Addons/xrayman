package xrayapitest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/common/logging"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/formats"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/xrayapi"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/xrayservice"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

const (
	testExecPath = "/usr/local/bin/xrayman/xray"

	testXRayCfg = `{
  "log": { "loglevel": "warning" },

  "api": {
    "tag": "api",
    "listen": "127.0.0.1:32997",
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

	testApiURL = "127.0.0.1:32997"
)

var testXRayUser = models.User{
	Name:      "username",
	VlessUUID: "aaaabbbbccccdddd",
}

var testXRayInbounds = []models.Inbound{
	{
		Tag:    "vlesstcp-reality",
		Format: &formats.VlessTCPReality{},
	},
}

// test service ctl
func TestXRayAPI(t *testing.T) {
	ctx := context.TODO()

	log, err := logging.New(zapcore.InfoLevel)
	require.NoError(t, err)

	// create xray api
	xrayapi, err := xrayapi.New(testApiURL, testXRayInbounds, xrayapi.WithLogger(log))
	assert.NoError(t, err)
	defer func() {
		err := xrayapi.Close(ctx)
		require.NoError(t, err)
	}()

	// ping stopped service
	err = xrayapi.Ping(ctx)
	assert.Error(t, err)

	// write xray config to file,
	// remove it after execution
	xray, err := xrayservice.New("")
	assert.NoError(t, err)
	err = xray.Start(t.Context(), testXRayCfg)
	assert.NoError(t, err)
	fmt.Println("started ok")
	time.Sleep(5 * time.Second)

	defer func() {
		err = xray.Close(t.Context())
		assert.NoError(t, err)
	}()

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
