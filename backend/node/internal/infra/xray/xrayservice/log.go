package xrayservice

import (
	"fmt"
	"strings"

	xapplog "github.com/xtls/xray-core/app/log"
	"github.com/xtls/xray-core/common/log"
	xlog "github.com/xtls/xray-core/common/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logHandler struct {
	log *zap.Logger
}

// Handle implements log.Handler.
func (h *logHandler) Handle(msg xlog.Message) {
	h.log.Log(h.translate(msg))
}

/*
xray log levels:

	0: "Unknown",
	1: "Error",
	2: "Warning",
	3: "Info",
	4: "Debug",
*/
func (h *logHandler) translate(msg xlog.Message) (zapcore.Level, string) {
	text := msg.String()
	const tag = "[xray] "
	switch {
	case strings.HasPrefix(text, "[Unknown]"):
		return zapcore.InfoLevel, tag + text[len("[Unknown]")+1:]
	case strings.HasPrefix(text, "[Error]"):
		return zapcore.ErrorLevel, tag + text[len("[Error]")+1:]
	case strings.HasPrefix(text, "[Warning]"):
		return zapcore.WarnLevel, tag + text[len("[Warning]")+1:]
	case strings.HasPrefix(text, "[Info]"):
		return zapcore.InfoLevel, tag + text[len("[Info]")+1:]
	case strings.HasPrefix(text, "[Debug]"):
		return zapcore.DebugLevel, tag + text[len("[Debug]")+1:]
	default:
		return zapcore.InfoLevel, tag + text
	}
}

var _ xlog.Handler = (*logHandler)(nil)

// Redirect XRay logs to zap logger
func redirectXRayLogs(l *zap.Logger) {
	lh := logHandler{log: l}
	hc := func(_ xapplog.LogType, o xapplog.HandlerCreatorOptions) (log.Handler, error) {
		if o.Path != "" {
			l.Warn(fmt.Sprintf("XRay logs redirected from %s", o.Path))
		}
		return &lh, nil
	}
	xlog.RegisterHandler(&lh)
	xapplog.RegisterHandlerCreator(xapplog.LogType_None, hc)
	xapplog.RegisterHandlerCreator(xapplog.LogType_Console, hc)
	xapplog.RegisterHandlerCreator(xapplog.LogType_File, hc)
	xapplog.RegisterHandlerCreator(xapplog.LogType_Event, hc)
}
