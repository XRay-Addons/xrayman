package xraycfg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testValidClientConfig = `[{
  "outbounds": [
    {
      "protocol": "vless",
      "settings": {
        "vnext": [
          {
            "users": [
              {
                "email": "{{ .VlessEmail }}",
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
                "email": "{{ .VlessEmail }}",
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
}]`

const testInvalidClientConfig = `[{
  "outbounds": [
    {
      "protocol": "vless",
      "settings": {
        "vnext": [
          {
            "users": [
              {
                "email": "{{ .VlessEmail }}",
                "encryption": "none",
                "flow": "xtls-rprx-vision",
                "id": "{{ .VlessUUIDA }}"
              },
              {
                "email": "{{ .VlessEmail }}",
                "encryption": "none",
                "flow": "xtls-rprx-vision",
                "id": "{{ .VlessUUIDA }}"
              }
    	      ]
          },
          {
            "users": [
              {
                "email": "{{ .VlessEmail }}",
                "encryption": "none",
                "flow": "xtls-rprx-vision",
                "id": "{{ .VlessUUIDA }}"
              },
              {
                "email": "{{ .VlessEmail }}",
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
                "email": "{{ .VlessEmail }}",
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
}]`

const testVlessEmailField = "VlessEmail"
const testVlessUUIDField = "VlessUUID"

func TestValidClientConfig(t *testing.T) {
	// test fields extraction
	uuidField, err := extractVlessUUIDField(testValidClientConfig)
	require.NoError(t, err)
	require.Equal(t, testVlessUUIDField, uuidField)

	emailField, err := extractVlessEmailField(testValidClientConfig)
	require.NoError(t, err)
	require.Equal(t, testVlessEmailField, emailField)

	// test full config
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "client_config.json")

	err = os.WriteFile(filePath, []byte(testValidClientConfig), 0o600)
	require.NoError(t, err)

	cfg, err := NewClientConfig(filePath)
	require.NoError(t, err)
	require.Equal(t, testVlessEmailField, cfg.cfg.VlessEmailField)
	require.Equal(t, testVlessUUIDField, cfg.cfg.VlessUUIDField)
}

func TestInvalidClientConfig(t *testing.T) {
	// test fields extraction
	_, err := extractVlessUUIDField(testInvalidClientConfig)
	require.Error(t, err)

	nameField, err := extractVlessEmailField(testInvalidClientConfig)
	require.NoError(t, err)
	require.Equal(t, testVlessEmailField, nameField)

	// test full config
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "client_config.json")

	err = os.WriteFile(filePath, []byte(testInvalidClientConfig), 0o600)
	require.NoError(t, err)

	_, err = NewClientConfig(filePath)
	require.Error(t, err)
}
