package clientcfg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testClientCfg = `{
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
    }
  ]
}`

const testNameField = "Name"
const testVlessUUIDField = "VlessUUID"

func TestClientCfg(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "client_config.json")

	err := os.WriteFile(filePath, []byte(testClientCfg), 0644)
	require.NoError(t, err)

	clientCfg, err := New(filePath, testNameField, testVlessUUIDField)
	require.NoError(t, err)
	_, err = clientCfg.GetClientConfigTemplate()
	require.NoError(t, err)
}
