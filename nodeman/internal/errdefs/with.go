package errdefs

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/xerr"
)

func WithFile(filename string) xerr.Option {
	return xerr.With(fmt.Sprintf("file: %s", filename))
}

func WrapWithFile(err error, path string) error {
	return xerr.Wrap(err, WithFile(path))
}
