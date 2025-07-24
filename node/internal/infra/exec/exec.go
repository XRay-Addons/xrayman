package exec

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

func Run(name string, args []string) (stdout, stderr string, err error) {
	cmd := exec.Command(name, args...)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()

	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()

	if err != nil {
		fullCmd := name + " " + strings.Join(args, " ")
		err = fmt.Errorf("%w: %s -> %s", errdefs.ErrExec, fullCmd, stdout)
	}

	return stdout, stderr, err
}
