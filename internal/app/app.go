package app

import (
	"context"
	"github.com/fpawel/dseries"
	"github.com/fpawel/mil82/internal"
	"github.com/fpawel/mil82/internal/data"
	"github.com/fpawel/mil82/internal/peer"
	"github.com/lxn/win"
	"github.com/powerman/structlog"
	"path/filepath"
	"sync"
)

func Run() {

	peer.AssertRunOnes()

	dseries.Open(filepath.Join(internal.DataDir(), "mil82.series.sqlite"))
	log.Println("charts: updated at", dseries.UpdatedAt())

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
	log.ErrIfFail(dseries.Close)
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
