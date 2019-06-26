package notify

import (
	"fmt"
	"github.com/fpawel/mil82/internal"
	"github.com/powerman/structlog"
)

type msg int

const (
	msgPanic msg = iota
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
