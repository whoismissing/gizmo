// Package debug provides methods for logging the code filename, line number, and function name
// at the point where the logging routine is called in order to assist the developer in debugging.
package debug

import (
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	LogB = log.New(LogWriter{}, "BEGIN: ", 0)
	LogE = log.New(LogWriter{}, "END: ", 0)
	// global variable Status is used to toggle debug logging
	Status = false
)

type LogWriter struct{}

// Write() obtains the code file, line number, and function name of its caller and outputs
// to stdout.
func (f LogWriter) Write(p []byte) (n int, err error) {
	pc, file, line, ok := runtime.Caller(4)
	if !ok {
		file = "?"
		line = 0
	}

	fn := runtime.FuncForPC(pc)
	var fnName string
	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".") + "()"
	}

	log.Printf("%s:%d %s: %s", filepath.Base(file), line, fnName, p)
	return len(p), nil
}

func LogBegin() {
	if Status {
		LogB.Println("beginning of function")
	}
}

func LogEnd() {
	if Status {
		LogE.Println("end of function")
	}
}
