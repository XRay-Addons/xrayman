package logging

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// pretty errors formatter
type PrettyEncoder struct {
	zapcore.Encoder
}

func (e *PrettyEncoder) Clone() zapcore.Encoder {
	return &PrettyEncoder{
		Encoder: e.Encoder.Clone(),
	}
}

func (e *PrettyEncoder) EncodeEntry(
	ent zapcore.Entry,
	fields []zapcore.Field,
) (*buffer.Buffer, error) {
	var okFields, errFields []zapcore.Field
	for _, f := range fields {
		if (f.Type == zapcore.ErrorType) && (f.Interface.(error) != nil) {
			errFields = append(errFields, f)
		} else {
			okFields = append(okFields, f)
		}
	}

	// non-errors messages:
	buf, err := e.Encoder.EncodeEntry(ent, okFields)

	// add errors messages:
	for _, f := range errFields {
		buf.AppendString(e.formatError(f.Key, f.Interface.(error)))
	}

	return buf, err
}

func (e *PrettyEncoder) formatError(name string, err error) string {
	errLines := strings.Split(fmt.Sprintf("%+v", err), "\n")
	for i, line := range errLines {
		if i > 0 {
			errLines[i] = "\t\t" + line
		}
	}
	return fmt.Sprintf("\t%s: %s\n", name, strings.Join(errLines, "\n"))
}

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

	stdout := zapcore.AddSync(os.Stdout)
	stderr := zapcore.AddSync(os.Stderr)

	core := zapcore.NewCore(
		encoder,
		stdout,
		zapcore.InfoLevel,
	)

	logger := zap.New(
		core,
		zap.ErrorOutput(stderr),
	)

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
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
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
		EncodeTime:    zapcore.TimeEncoderOfLayout("15:04:05"),
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}
}
