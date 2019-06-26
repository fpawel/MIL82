package app

import (
	"context"
	"github.com/ansel1/merry"
	"github.com/fpawel/mil82/internal/api/notify"
	"github.com/fpawel/mil82/internal/data"
	"github.com/getlantern/systray"
	"github.com/lxn/win"
	"github.com/powerman/must"
	"github.com/powerman/structlog"
	"os"
	"os/exec"
	"path/filepath"
)

func Run() {
	initLog()
	data.Open(false)
	notify.InitWindow("")
	go sysTray(notify.W.CloseWindow)

	var cancel func()
	ctxApp, cancel = context.WithCancel(context.TODO())

	closeHttpServer := startHttpServer()

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

}

func sysTray(onClose func()) {
	systray.Run(func() {
		systray.SetIcon(must.ReadFile(filepath.Join(filepath.Dir(os.Args[0]), "assets", "img", "app.ico")))
		systray.SetTitle("Производство ЭХЯ CO")
		systray.SetTooltip("Производство ЭХЯ CO")
		mRunGUIApp := systray.AddMenuItem("Показать", "Показать окно приложения")
		mQuitOrig := systray.AddMenuItem("Закрыть", "Закрыть приложение")

		go func() {
			for {
				select {
				case <-mRunGUIApp.ClickedCh:
					if err := runGUI(); err != nil {
						panic(merry.Append(err, "не удалось запустить elcoui.exe"))
					}
				case <-mQuitOrig.ClickedCh:
					systray.Quit()
					onClose()
				}
			}
		}()
	}, func() {
	})
}

func runGUI() error {
	fileName := filepath.Join(filepath.Dir(os.Args[0]), "mil82gui.exe")
	err := exec.Command(fileName).Start()
	if err != nil {
		return merry.Append(err, fileName)
	}
	return nil
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
