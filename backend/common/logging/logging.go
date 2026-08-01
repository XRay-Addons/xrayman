package logging

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.Logger, error) {
	const human = true
	var encoder zapcore.Encoder
	if human {
		encCfg := humanReadableConfig()
		encoder = &PrettyEncoder{
			Encoder: zapcore.NewConsoleEncoder(encCfg),
		}
	} else {
		encCfg := machineReadableConfig()
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	// All -> stdout
	stdoutCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	// Warn, Err -> stderr
	/*stderrCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stderr),
		zapcore.WarnLevel,
	)*/

	core := zapcore.NewTee(stdoutCore)

	logger := zap.New(core)

	return logger, nil
}

//
// ---------- Encoder configs ----------
//

func machineReadableConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
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
}

func humanReadableConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		MessageKey:    "msg",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    "\n",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02 15:04:05 UTC"))
		},
		EncodeLevel:  zapcore.CapitalColorLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
}
