package notify

import (
	"fmt"
	"github.com/fpawel/mil82/internal/api/types"
	"github.com/fpawel/mil82/internal/peer"
	"github.com/powerman/structlog"
)

type msg int

const (
	msgPanic msg = iota
	msgReadVar
	msgAddrError
	msgWorkStarted
	msgWorkComplete
	msgWarning
	msgDelay
)

func Panic(log *structlog.Logger, arg string) {
	if log != nil {
		log.Info(peer.WindowClassName+": Panic: "+fmt.Sprintf("%+v", arg), "MSG", msgPanic)
	}
	go peer.W.NotifyStr(uintptr(msgPanic), arg)
}

func Panicf(log *structlog.Logger, format string, a ...interface{}) {
	if log != nil {
		log.Info(peer.WindowClassName+": Panic: "+fmt.Sprintf(format, a...), "MSG", msgPanic)
	}
	go peer.W.Notifyf(uintptr(msgPanic), format, a...)
}
func ReadVar(log *structlog.Logger, arg types.AddrVarValue) {
	if log != nil {
		log.Info(peer.WindowClassName+": ReadVar: "+fmt.Sprintf("%+v", arg), "MSG", msgReadVar)
	}
	go peer.W.NotifyJson(uintptr(msgReadVar), arg)
}

func AddrError(log *structlog.Logger, arg types.AddrError) {
	if log != nil {
		log.Info(peer.WindowClassName+": AddrError: "+fmt.Sprintf("%+v", arg), "MSG", msgAddrError)
	}
	go peer.W.NotifyJson(uintptr(msgAddrError), arg)
}

func WorkStarted(log *structlog.Logger, arg string) {
	if log != nil {
		log.Info(peer.WindowClassName+": WorkStarted: "+fmt.Sprintf("%+v", arg), "MSG", msgWorkStarted)
	}
	go peer.W.NotifyStr(uintptr(msgWorkStarted), arg)
}

func WorkStartedf(log *structlog.Logger, format string, a ...interface{}) {
	if log != nil {
		log.Info(peer.WindowClassName+": WorkStarted: "+fmt.Sprintf(format, a...), "MSG", msgWorkStarted)
	}
	go peer.W.Notifyf(uintptr(msgWorkStarted), format, a...)
}
func WorkComplete(log *structlog.Logger, arg types.WorkResultInfo) {
	if log != nil {
		log.Info(peer.WindowClassName+": WorkComplete: "+fmt.Sprintf("%+v", arg), "MSG", msgWorkComplete)
	}
	go peer.W.NotifyJson(uintptr(msgWorkComplete), arg)
}

func Warning(log *structlog.Logger, arg string) {
	if log != nil {
		log.Info(peer.WindowClassName+": Warning: "+fmt.Sprintf("%+v", arg), "MSG", msgWarning)
	}
	go peer.W.NotifyStr(uintptr(msgWarning), arg)
}

func Warningf(log *structlog.Logger, format string, a ...interface{}) {
	if log != nil {
		log.Info(peer.WindowClassName+": Warning: "+fmt.Sprintf(format, a...), "MSG", msgWarning)
	}
	go peer.W.Notifyf(uintptr(msgWarning), format, a...)
}
func Delay(log *structlog.Logger, arg types.DelayInfo) {
	if log != nil {
		log.Info(peer.WindowClassName+": Delay: "+fmt.Sprintf("%+v", arg), "MSG", msgDelay)
	}
	go peer.W.NotifyJson(uintptr(msgDelay), arg)
}
