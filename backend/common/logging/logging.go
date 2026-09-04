package logging

import (
	"errors"
	"os"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(lvl zapcore.Level) *zap.Logger {
	// encode config
	encCfg := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02 15:04:05 UTC"))
		},
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encCfg)

	syncer := &logSyncer{
		zapcore.AddSync(os.Stdout),
	}
	stdoutCore := zapcore.NewCore(encoder, syncer, lvl)

	core := zapcore.NewTee(stdoutCore)
	logger := zap.New(core)

	return logger
}

// log syncer for avoid disgusting error
// 2026/09/03 17:17:56 sync /dev/stdout: inappropriate ioctl for device
type logSyncer struct {
	zapcore.WriteSyncer
}

func (s logSyncer) Sync() error {
	err := s.WriteSyncer.Sync()
	switch {
	case errors.Is(err, syscall.ENOTTY):
		return nil // linux, osx, ???
	case errors.Is(err, syscall.ENODEV):
		return nil // linux, osx, ???
	case errors.Is(err, syscall.EINVAL):
		return nil // linux, osx, ???
	case errors.Is(err, syscall.EBADF):
		return nil // linux, osx, ???
	case errors.Is(err, syscall.Errno(6)):
		return nil // windows
	default:
		return err
	}
}
