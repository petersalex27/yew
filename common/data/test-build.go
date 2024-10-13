//go:build test
// +build test

package data

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/petersalex27/yew/api"
)

var mu sync.Mutex

const logPath string = "../../logs/log-fail-caller.txt"

var nWrites byte = 0 // max out writes at wrap around to 0

func logCaller(pc uintptr, fileName string, line int, ok bool) {
	mu.Lock()
	defer mu.Unlock()
	if !ok {
		return
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	if nWrites == 0 { // clear file if it's the first write or after wrap around
		f.Truncate(0)
	}
	nWrites++ // increment number of writes, possibly wrapping around

	pid := os.Getpid()
	caller := runtime.FuncForPC(pc).Name()
	var reason string = "unknown"
	rsnPC, _, _, ok := runtime.Caller(1) // get the caller of this function
	if ok {
		if fn := runtime.FuncForPC(rsnPC); fn != nil {
			reason = fn.Name()
		}
	}
	fmt.Fprintf(f, "reason: %s\n\tpid: %d / file name: %q / line: %d / caller: %s\n", reason, pid, fileName, line, caller)
}

func Fail[a api.Node](msg string, positioned api.Positioned) Either[Ers, a] {
	logCaller(runtime.Caller(1))
	return __Fail[a](msg, positioned)
}

func PassErs[b api.Node](e Ers) Either[Ers, b] {
	logCaller(runtime.Caller(1))
	return __PassErs[b](e)
}