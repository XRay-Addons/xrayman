package xraycfg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const testEmptyConfig = `
{
  "log": {},
  "stats": {},

  "inbounds": [
    {
      "protocol": "vless",
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

var testUsers = []User{
	{
		Name: "user1",
		UUID: "userid1",
	},
	{
		Name: "user2",
		UUID: "userid2",
	},
}

func TestCreateServerConfig(t *testing.T) {
	filledConfig, err := CreateServerConfig(testEmptyConfig, testUsers)
	require.NoError(t, err)
	require.JSONEq(t, testFilledConfig, filledConfig)
}
