package launchctl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testExecPath = "/usr/local/bin/xrayman/xray"
	testCfgPath  = "/usr/local/bin/xrayman/server.yaml"
	testPlist    = `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>xray</string>
    <key>ProgramArguments</key>
    <array>
      <string>/usr/local/bin/xrayman/xray</string>
      <string>-config</string>
      <string>/usr/local/bin/xrayman/server.yaml</string>
    </array>
    <key>RunAtLoad</key>
    <false/>
    
    <key>ProcessType</key>
    <string>Background</string>
  </dict>
</plist>`
)

func TestMakePlist(t *testing.T) {
	plist, err := makePlist(testExecPath, testCfgPath)
	require.NoError(t, err)
	require.Equal(t, testPlist, string(plist))
}

// test service ctl?
// create and install mock service?
