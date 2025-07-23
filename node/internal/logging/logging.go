package logging

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.Logger, error) {
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// no caller
	cfg.EncoderConfig.CallerKey = ""
	// no stacktrace
	cfg.EncoderConfig.StacktraceKey = ""

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("logger init: %w: %v", errdefs.ErrIPE, err)
	}
	return logger, nil
}
