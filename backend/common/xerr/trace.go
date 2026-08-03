package xerr

import (
	"fmt"
	"runtime"
	"strings"
)

type trace struct {
	frames []string
}

func (t *trace) Format(f fmt.State, verb rune) {
	if len(t.frames) == 0 {
		fmt.Fprint(f, "callstack is empty")
		return
	}
	format := "%" + string(verb)
	traceText := "\n\t-> " + strings.Join(t.frames, "\n\t-> ")
	fmt.Fprintf(f, format, traceText)
}

func getTrace(skip int) trace {
	const maxStackDepth = 32
	pcs := make([]uintptr, maxStackDepth)
	n := runtime.Callers(skip+1, pcs)
	callFrames := runtime.CallersFrames(pcs[:n])

	var frames []string
	for {
		f, more := callFrames.Next()
		frame := fmt.Sprintf("%s:%d %s", f.File, f.Line, getFuncName(f))
		frames = append(frames, frame)
		if !more {
			break
		}
	}

	return trace{frames}
}

func getFuncName(f runtime.Frame) string {
	funcName := f.Func.Name()
	pos := strings.LastIndex(funcName, ".")
	if pos == -1 || pos == len(funcName)-1 {
		return "<anonymous func>"
	}
	return funcName[pos+1:]
}
