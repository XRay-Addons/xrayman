package cmd

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

func Run(name string, args []string) (stdout, stderr string, err error) {

	cmd := exec.Command(name, args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()

	if err != nil {
		err = fmt.Errorf("%w: %v", errdefs.ErrCmd, err)
	}
	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()

	return
}
