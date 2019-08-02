package app

import (
	"context"
	"github.com/fpawel/mil82/internal/cfg"
	"github.com/fpawel/mil82/internal/charts"
	"github.com/fpawel/mil82/internal/data"
	"github.com/fpawel/mil82/internal/peer"
	"github.com/lxn/win"
	"github.com/powerman/structlog"
	"sync"
)

func Run() {
	cfg.Open()

	var cancel func()
	ctxApp, cancel = context.WithCancel(context.TODO())
	closeHttpServer := startHttpServer()
	peer.Init("")
	// цикл оконных сообщений
	for {
		var msg win.MSG
		if win.GetMessage(&msg, 0, 0, 0) == 0 {
			break
		}
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
	cancel()
	closeHttpServer()
	peer.Close()
	log.ErrIfFail(data.DB.Close)
	log.ErrIfFail(charts.DB.Close)
}

type peerNotifier struct{}

func (_ peerNotifier) OnStarted() {
	peer.InitPeer()
}

func (_ peerNotifier) OnClosed() {
	peer.ResetPeer()
	cancelWorkFunc()
}

var (
	ctxApp         context.Context
	cancelWorkFunc = func() {}
	skipDelayFunc  = func() {}
	wgWork         sync.WaitGroup
	log            = structlog.New()
)
