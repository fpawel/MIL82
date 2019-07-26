package app

import (
	"context"
	"fmt"
	"github.com/ansel1/merry"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/gohelp"
	"github.com/fpawel/mil82/internal/api/notify"
)

func switchGas(n int) error {

	log := gohelp.LogPrependSuffixKeys(log, "пневмоблок", n)
	log.Info("переключение")

	_, err := modbus.Request{
		Addr:     5,
		ProtoCmd: 0x10,
		Data: []byte{
			0, 0x10, 0, 1, 2, 0, byte(n),
		},
	}.GetResponse(log, ctxWork, portGas, nil)
	if err == nil {
		return nil
	}

	s := "Не удалось "
	if n == 0 {
		s += "отключить газ"
	} else {
		s += fmt.Sprintf("подать ПГС%d", n)
	}

	s += ": " + err.Error() + ".\n\n"

	if n == 0 {
		s += "Отключите газ"
	} else {
		s += fmt.Sprintf("Подайте ПГС%d", n)
	}
	s += " вручную."
	notify.Warning(log, s)
	if merry.Is(ctxWork.Err(), context.Canceled) {
		return err
	}
	log.Warn("проигнорирована ошибка связи", "error", err)

	return nil
}
