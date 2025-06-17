package xraycfg

import (
	"testing"

	"github.com/XRay-Addons/xrayman/shared/models"
	"github.com/stretchr/testify/require"
)

const testEmptyConfig = `
{
  "log": {},
  "stats": {},

  "inbounds": [
    {
      "protocol": "vless",
      "tag": "reality-in",
      "settings": {
        "clients": [],
        "decryption": "none"
      },
      "streamSettings": {
        "network": "tcp",
        "security": "reality"
      }
    },
    {
      "protocol": "vless",
      "tag": "xhttp-in",
      "settings": {
        "clients": [],
        "decryption": "none"
      },
      "streamSettings": {
        "network": "xhttp"
      }
    }
  ],
  "outbounds": []
}
`

const testFilledConfig = `
{
  "log": {},
  "stats": {},

  "inbounds": [
    {
      "protocol": "vless",
      "tag": "reality-in",
      "settings": {
        "clients": [
          {
            "email": "user1",
            "flow": "xtls-rprx-vision",
            "id": "userid1"
          },
          {
            "email": "user2",
            "flow": "xtls-rprx-vision",
            "id": "userid2"
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
      "protocol": "vless",
      "tag": "xhttp-in",
      "settings": {
        "clients": [
          {
            "email": "user1",
            "id": "userid1"
          },
          {
            "email": "user2",
            "id": "userid2"
          }		
		    ],
        "decryption": "none"
      },
      "streamSettings": {
        "network": "xhttp"
      }
    }
  ],
  "outbounds": []
}
`

var testUsers = []models.User{
	{
		Name: "user1",
		UUID: "userid1",
	},
	{
		Name: "user2",
		UUID: "userid2",
	},
}

func TestAddUsers(t *testing.T) {
	inbounds, err := GetInbounds(testEmptyConfig)
	require.NoError(t, err)
	filledConfig, err := AddUsers(testEmptyConfig, inbounds, testUsers)
	require.NoError(t, err)
	require.JSONEq(t, testFilledConfig, filledConfig)
}
