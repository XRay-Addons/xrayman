package exec

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

func Run(ctx context.Context, name string, args []string) (stdout, stderr string, err error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()

	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()

	if err != nil {
		fullCmd := name + " " + strings.Join(args, " ")
		err = errdefs.Wrap(err, errdefs.WithStack(),
			errdefs.Withf("cmd: %s, out: %s", fullCmd, stdout))
	}

	return stdout, stderr, err
}
