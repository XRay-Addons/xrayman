package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Log(log *zap.Logger, msg string, fields ...zap.Field) {
	// if fields contains error, log as error, else - as info
	lvl := zap.InfoLevel
	for _, f := range fields {
		if f.Type == zapcore.ErrorType {
			lvl = zap.ErrorLevel
			break
		}
	}
	log.Log(lvl, msg, fields...)
}
