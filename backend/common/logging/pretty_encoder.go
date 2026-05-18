package logging

import (
	"fmt"
	"strings"

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
