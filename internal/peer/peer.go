package peer

import (
	"github.com/fpawel/gohelp/copydata"
	"github.com/fpawel/gohelp/winapp"
	"github.com/lxn/win"
	"github.com/powerman/structlog"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	ServerWindowClassName = "Mil82ServerWindow"
	WindowClassName       = "TMainFormMil82"
)

// окно для отправки сообщений WM_COPYDATA дельфи-приложению
var (
	W *copydata.NotifyWindow
)

func Init(sourceWindowClassNameSuffix string) {
	W = copydata.NewNotifyWindow(
		ServerWindowClassName+sourceWindowClassNameSuffix,
		WindowClassName)
}

func RunGUI() {
	if err := exec.Command(filepath.Join(filepath.Dir(os.Args[0]), "mil82gui.exe")).Start(); err != nil {
		panic(err)
	}
}

func CloseGUI() {
	winapp.EnumWindowsWithClassName(func(hWnd win.HWND, winClassName string) {
		if winClassName == WindowClassName {
			r := win.PostMessage(hWnd, win.WM_CLOSE, 0, 0)
			log.Debug("close peer window", "syscall", r)
		}
	})
}

var log = structlog.New()
