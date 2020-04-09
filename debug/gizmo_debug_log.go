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
    Status = false
)

type LogWriter struct{}

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
