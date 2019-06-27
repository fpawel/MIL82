package app

import (
	"context"
	"github.com/fpawel/gohelp/must"
	"github.com/fpawel/mil82/internal/api/notify"
	"github.com/fpawel/mil82/internal/cfg"
	"github.com/fpawel/mil82/internal/data"
	"github.com/getlantern/systray"
	"github.com/lxn/win"
	"github.com/powerman/structlog"
	"os"
	"os/exec"
	"path/filepath"
)

func Run() {
	initLog()
	data.Open(false)
	cfg.Open()
	notify.InitWindow("")
	go sysTray(notify.W.CloseWindow)

	var cancel func()
	ctxApp, cancel = context.WithCancel(context.TODO())

	closeHttpServer := startHttpServer()
	go runGUI()
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
	log.ErrIfFail(data.DB.Close, "defer", "close products db")
	cfg.Save()
}

func sysTray(onClose func()) {
	systray.Run(func() {
		systray.SetIcon(must.ReadFile(filepath.Join(filepath.Dir(os.Args[0]), "assets", "img", "app.ico")))
		systray.SetTitle("МИЛ-82")
		systray.SetTooltip("МИЛ-82")
		mRunGUIApp := systray.AddMenuItem("Показать", "Показать окно приложения")
		mQuitOrig := systray.AddMenuItem("Закрыть", "Закрыть приложение")

		go func() {
			for {
				select {
				case <-mRunGUIApp.ClickedCh:
					go runGUI()
				case <-mQuitOrig.ClickedCh:
					systray.Quit()
					onClose()
				}
			}
		}()
	}, func() {
	})
}

func runGUI() {
	if err := exec.Command(filepath.Join(filepath.Dir(os.Args[0]), "mil82gui.exe")).Start(); err != nil {
		panic(err)
	}
}

func initLog() {
	structlog.DefaultLogger.
		SetPrefixKeys(
			structlog.KeyApp, structlog.KeyPID, structlog.KeyLevel, structlog.KeyUnit, structlog.KeyTime,
		).
		SetDefaultKeyvals(
			structlog.KeyApp, filepath.Base(os.Args[0]),
			structlog.KeySource, structlog.Auto,
		).
		SetSuffixKeys(
			structlog.KeyStack,
		).
		SetSuffixKeys(structlog.KeySource).
		SetKeysFormat(map[string]string{
			structlog.KeyTime:   " %[2]s",
			structlog.KeySource: " %6[2]s",
			structlog.KeyUnit:   " %6[2]s",
			"config":            " %+[2]v",
			"запрос":            " %[1]s=`% [2]X`",
			"ответ":             " %[1]s=`% [2]X`",
			"работа":            " %[1]s=`%[2]s`",
		}).SetTimeFormat("15:04:05")
}
