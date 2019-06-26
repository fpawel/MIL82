package notify

import (
	"github.com/fpawel/gohelp/copydata"
	"github.com/fpawel/mil82/internal"
)

// окно для отправки сообщений WM_COPYDATA дельфи-приложению
var W *copydata.NotifyWindow

func InitWindow(sourceWindowClassNameSuffix string) {
	W = copydata.NewNotifyWindow(
		internal.ServerWindowClassName+sourceWindowClassNameSuffix,
		internal.PeerWindowClassName)
}
