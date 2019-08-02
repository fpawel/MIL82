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
	msgEndDelay
	msgStatus
)

func Panic(log *structlog.Logger, arg string) {
	if log != nil {
		msgPanic.Log(log)(peer.WindowClassName+": Panic: "+fmt.Sprintf("%+v", arg), "MSG", msgPanic)
	}
	go peer.NotifyStr(uintptr(msgPanic), arg)
}

func Panicf(log *structlog.Logger, format string, a ...interface{}) {
	if log != nil {
		msgPanic.Log(log)(peer.WindowClassName+": Panic: "+fmt.Sprintf(format, a...), "MSG", msgPanic)
	}
	go peer.Notifyf(uintptr(msgPanic), format, a...)
}
func ReadVar(log *structlog.Logger, arg types.AddrVarValue) {
	if log != nil {
		msgReadVar.Log(log)(peer.WindowClassName+": ReadVar: "+fmt.Sprintf("%+v", arg), "MSG", msgReadVar)
	}
	go peer.NotifyJson(uintptr(msgReadVar), arg)
}

func AddrError(log *structlog.Logger, arg types.AddrError) {
	if log != nil {
		msgAddrError.Log(log)(peer.WindowClassName+": AddrError: "+fmt.Sprintf("%+v", arg), "MSG", msgAddrError)
	}
	go peer.NotifyJson(uintptr(msgAddrError), arg)
}

func WorkStarted(log *structlog.Logger, arg string) {
	if log != nil {
		msgWorkStarted.Log(log)(peer.WindowClassName+": WorkStarted: "+fmt.Sprintf("%+v", arg), "MSG", msgWorkStarted)
	}
	go peer.NotifyStr(uintptr(msgWorkStarted), arg)
}

func WorkStartedf(log *structlog.Logger, format string, a ...interface{}) {
	if log != nil {
		msgWorkStarted.Log(log)(peer.WindowClassName+": WorkStarted: "+fmt.Sprintf(format, a...), "MSG", msgWorkStarted)
	}
	go peer.Notifyf(uintptr(msgWorkStarted), format, a...)
}
func WorkComplete(log *structlog.Logger, arg types.WorkResultInfo) {
	if log != nil {
		msgWorkComplete.Log(log)(peer.WindowClassName+": WorkComplete: "+fmt.Sprintf("%+v", arg), "MSG", msgWorkComplete)
	}
	go peer.NotifyJson(uintptr(msgWorkComplete), arg)
}

func Warning(log *structlog.Logger, arg string) {
	if log != nil {
		msgWarning.Log(log)(peer.WindowClassName+": Warning: "+fmt.Sprintf("%+v", arg), "MSG", msgWarning)
	}
	peer.NotifyStr(uintptr(msgWarning), arg)
}

func Warningf(log *structlog.Logger, format string, a ...interface{}) {
	if log != nil {
		msgWarning.Log(log)(peer.WindowClassName+": Warning: "+fmt.Sprintf(format, a...), "MSG", msgWarning)
	}
	peer.Notifyf(uintptr(msgWarning), format, a...)
}
func Delay(log *structlog.Logger, arg types.DelayInfo) {
	if log != nil {
		msgDelay.Log(log)(peer.WindowClassName+": Delay: "+fmt.Sprintf("%+v", arg), "MSG", msgDelay)
	}
	go peer.NotifyJson(uintptr(msgDelay), arg)
}

func EndDelay(log *structlog.Logger, arg string) {
	if log != nil {
		msgEndDelay.Log(log)(peer.WindowClassName+": EndDelay: "+fmt.Sprintf("%+v", arg), "MSG", msgEndDelay)
	}
	go peer.NotifyStr(uintptr(msgEndDelay), arg)
}

func EndDelayf(log *structlog.Logger, format string, a ...interface{}) {
	if log != nil {
		msgEndDelay.Log(log)(peer.WindowClassName+": EndDelay: "+fmt.Sprintf(format, a...), "MSG", msgEndDelay)
	}
	go peer.Notifyf(uintptr(msgEndDelay), format, a...)
}
func Status(log *structlog.Logger, arg string) {
	if log != nil {
		msgStatus.Log(log)(peer.WindowClassName+": Status: "+fmt.Sprintf("%+v", arg), "MSG", msgStatus)
	}
	go peer.NotifyStr(uintptr(msgStatus), arg)
}

func Statusf(log *structlog.Logger, format string, a ...interface{}) {
	if log != nil {
		msgStatus.Log(log)(peer.WindowClassName+": Status: "+fmt.Sprintf(format, a...), "MSG", msgStatus)
	}
	go peer.Notifyf(uintptr(msgStatus), format, a...)
}
