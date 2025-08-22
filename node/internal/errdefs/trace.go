package errdefs

import (
	"fmt"
	"runtime"
	"strings"
)

func getTrace(skip int) []string {
	callstack := getCallstack(skip + 1)
	return callstack
}

func getCallstack(skip int) []string {
	const maxStackDepth = 32
	pcs := make([]uintptr, maxStackDepth)
	n := runtime.Callers(skip+1, pcs)
	callFrames := runtime.CallersFrames(pcs[:n])

	var frames []string
	for {
		rawFrame, more := callFrames.Next()
		frames = append(frames, makePrettyFrame(rawFrame))
		if !more {
			break
		}
	}

	reverseFrames(frames)

	return frames
}

func makePrettyFrame(f runtime.Frame) string {
	return fmt.Sprintf("%s:%d %s", f.File, f.Line, getFuncName(f))
}

func getFuncName(f runtime.Frame) string {
	funcName := f.Func.Name()
	pos := strings.LastIndex(funcName, ".")
	if pos == -1 || pos == len(funcName)-1 {
		return "<anonymous func>"
	}
	return funcName[pos+1:]
}

func reverseFrames(frames []string) {
	for i, j := 0, len(frames)-1; i < j; i, j = i+1, j-1 {
		frames[i], frames[j] = frames[j], frames[i]
	}
}
