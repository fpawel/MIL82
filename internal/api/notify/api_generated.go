package notify

import (
	"fmt"
	"github.com/fpawel/mil82/internal"
	"github.com/fpawel/mil82/internal/api"
	"github.com/powerman/structlog"
)

type msg int

const (
	msgPanic msg = iota
	msgReadVar
	msgError
)

func Panic(log *structlog.Logger, arg string) {
	if log != nil {
		log.Debug(internal.PeerWindowClassName, "Panic", arg, "MSG", msgPanic)
	}
	W.NotifyStr(uintptr(msgPanic), arg)
}

func Panicf(log *structlog.Logger, format string, a ...interface{}) {
	if log != nil {
		log.Debug(internal.PeerWindowClassName, "Panic", fmt.Sprintf(format, a...), "MSG", msgPanic)
	}
	W.Notifyf(uintptr(msgPanic), format, a...)
}
func ReadVar(log *structlog.Logger, arg api.AddrVarValue) {
	if log != nil {
		log.Debug(internal.PeerWindowClassName, "ReadVar", arg, "MSG", msgReadVar)
	}
	W.NotifyJson(uintptr(msgReadVar), arg)
}

func Error(log *structlog.Logger, arg string) {
	if log != nil {
		log.Debug(internal.PeerWindowClassName, "Error", arg, "MSG", msgError)
	}
	W.NotifyStr(uintptr(msgError), arg)
}

func Errorf(log *structlog.Logger, format string, a ...interface{}) {
	if log != nil {
		log.Debug(internal.PeerWindowClassName, "Error", fmt.Sprintf(format, a...), "MSG", msgError)
	}
	W.Notifyf(uintptr(msgError), format, a...)
}
