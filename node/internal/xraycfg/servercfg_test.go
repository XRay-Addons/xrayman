package xraycfg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/stretchr/testify/require"
)

const testServerCfg = `{
  "api": {
    "tag": "api",
    "listen": "127.0.0.1:32999"
  },

  "inbounds": [
    {
      "tag": "reality-in",
      "listen": "0.0.0.0",
      "port": 443,
      "protocol": "vless",
      "settings": {
        "clients": [],
        "fallbacks": [
          {
            "dest": "@xhttp-input-socket",
            "xver": 1
          }
        ],
        "decryption": "none"
      },
      "streamSettings": {
        "network": "tcp",
        "security": "reality"
      }
    },
    {
      "tag": "xhttp-in",
      "listen": "@xhttp-input-socket",
      "protocol": "vless",
      "settings": {
        "clients": [],
        "decryption": "none"
      },
      "streamSettings": {
        "network": "xhttp",
        "xhttpSettings": {
          "mode": "stream-one",
          "path": "come-on-xhttp"
        },
        "sockopt": {
          "acceptProxyProtocol": true
        }
      }
    }
  ]
}`

const testApiURL = "127.0.0.1:32999"

var testInbounds = []models.Inbound{
	{Tag: "reality-in", Type: models.VlessTcpReality},
	{Tag: "xhttp-in", Type: models.VlessXHTTP},
}

func TestServiceCfg(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "service_config.json")

	err := os.WriteFile(filePath, []byte(testServerCfg), 0644)
	require.NoError(t, err)

	serviceCfg, err := NewServerCfg(filePath)
	require.NoError(t, err)

	apiURL, err := serviceCfg.GetApiURL()
	require.NoError(t, err)
	require.Equal(t, testApiURL, apiURL)

	inbounds, err := serviceCfg.GetInbounds()
	require.NoError(t, err)
	require.Equal(t, testInbounds, inbounds)
}
