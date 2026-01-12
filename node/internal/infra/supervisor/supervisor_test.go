package supervisor

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/infra/supervisor/supervisorapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	testServiceName = "testxrayservice"
	testExecPath    = "/usr/local/bin/xrayman/xray"
)

const testXRayCfg = `{
  "log": { "loglevel": "warning" },
  "inbounds": [],
  "outbounds": [{ "protocol": "freedom", "settings": {} }]
}`

// test service ctl
func TestXrayctl(t *testing.T) {
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

	log := zaptest.NewLogger(t)

	// create xray service
	command := []string{testExecPath, "-config", testCfgPath}
	xrayctl, err := New(testServiceName, command, log)
	require.NoError(t, err)
	defer func() {
		err := xrayctl.Close(ctx)
		assert.NoError(t, err)
	}()

	// wait 1 second till service initialization
	time.Sleep(1 * time.Second)

	// check status is stopped now
	status, err := xrayctl.Status(ctx)
	assert.NoError(t, err)
	assert.Equal(t, supervisorapi.StatusStopped, status)

	// check scenario: [[start, status] x 2, [stop, status] x 2] x 2
	for range 2 {
		for range 2 {
			// start node
			err = xrayctl.Start(ctx)
			assert.NoError(t, err)

			// wait 1 second till service started
			time.Sleep(1 * time.Second)

			// check status is running
			status, err = xrayctl.Status(ctx)
			assert.NoError(t, err)
			assert.Equal(t, supervisorapi.StatusRunning, status)
		}

		// check stop two times
		for range 2 {
			// stop node
			err = xrayctl.Stop(ctx)
			assert.NoError(t, err)

			// wait 1 second till service stopped
			time.Sleep(1 * time.Second)

			// check status is running
			status, err = xrayctl.Status(ctx)
			assert.NoError(t, err)
			assert.Equal(t, supervisorapi.StatusStopped, status)
		}
	}
}
