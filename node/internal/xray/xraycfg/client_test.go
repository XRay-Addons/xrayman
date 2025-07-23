package xraycfg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testValidClientCfg = `{
  "outbounds": [
    {
      "protocol": "vless",
      "settings": {
        "vnext": [
          {
            "users": [
              {
                "email": "{{ .Name }}",
                "encryption": "none",
                "flow": "xtls-rprx-vision",
                "id": "{{ .VlessUUID }}"
              }
    	      ]
          }
        ]
      }
    },
    {
      "protocol": "vless",
      "settings": {
        "vnext": [
          {
            "users": [
              {
                "email": "{{ .Name }}",
                "encryption": "none",
                "flow": "xtls-rprx-vision",
                "id": "{{ .VlessUUID }}"
              }
    	      ]
          }
        ]
      }
    }
  ]
}`

const testInvalidClientCfg = `{
  "outbounds": [
    {
      "protocol": "vless",
      "settings": {
        "vnext": [
          {
            "users": [
              {
                "email": "{{ .Name }}",
                "encryption": "none",
                "flow": "xtls-rprx-vision",
                "id": "{{ .VlessUUIDA }}"
              },
              {
                "email": "{{ .Name }}",
                "encryption": "none",
                "flow": "xtls-rprx-vision",
                "id": "{{ .VlessUUIDA }}"
              }
    	      ]
          },
          {
            "users": [
              {
                "email": "{{ .Name }}",
                "encryption": "none",
                "flow": "xtls-rprx-vision",
                "id": "{{ .VlessUUIDA }}"
              },
              {
                "email": "{{ .Name }}",
                "encryption": "none",
                "flow": "xtls-rprx-vision",
                "id": "{{ .VlessUUIDA }}"
              }
    	      ]
          }
        ]
      }
    },
    {
      "protocol": "vless",
      "settings": {
        "vnext": [
          {
            "users": [
              {
                "email": "{{ .Name }}",
                "encryption": "none",
                "flow": "xtls-rprx-vision",
                "id": "{{ .VlessUUIDB }}"
              }
    	    ]
          }
        ]
      }
    }
  ]
}`

const testNameField = "Name"
const testVlessUUIDField = "VlessUUID"

func TestValidClientCfg(t *testing.T) {
	// test fields extraction
	uuidField, err := extractVlessUUIDField(testValidClientCfg)
	require.NoError(t, err)
	require.Equal(t, testVlessUUIDField, uuidField)

	nameField, err := extractNameField(testValidClientCfg)
	require.NoError(t, err)
	require.Equal(t, testNameField, nameField)

	// test full config
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "client_config.json")

	err = os.WriteFile(filePath, []byte(testValidClientCfg), 0644)
	require.NoError(t, err)

	cfg, err := NewClientCfg(filePath)
	require.NoError(t, err)
	require.Equal(t, testNameField, cfg.cfg.UserNameField)
	require.Equal(t, testVlessUUIDField, cfg.cfg.VlessUUIDField)
}

func TestInvalidClientCfg(t *testing.T) {
	// test fields extraction
	_, err := extractVlessUUIDField(testInvalidClientCfg)
	require.Error(t, err)

	nameField, err := extractNameField(testInvalidClientCfg)
	require.NoError(t, err)
	require.Equal(t, testNameField, nameField)

	// test full config
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "client_config.json")

	err = os.WriteFile(filePath, []byte(testInvalidClientCfg), 0644)
	require.NoError(t, err)

	_, err = NewClientCfg(filePath)
	require.Error(t, err)
}
