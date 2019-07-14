package app

import (
	"context"
	"github.com/ansel1/merry"
	"github.com/fpawel/gohelp"
	"github.com/fpawel/gohelp/helpstr"
	"github.com/fpawel/mil82/internal/api/notify"
	"github.com/fpawel/mil82/internal/api/types"
	"github.com/fpawel/mil82/internal/cfg"
	"github.com/fpawel/mil82/internal/party"
	"time"
)

func delay(what string, duration time.Duration) error {

	originalLog := log
	defer func() {
		log = originalLog
	}()
	log = gohelp.LogWithKeys(log,
		"фоновый_опрос", what,
		"общая_длительность_задержки", helpstr.FormatDuration(duration),
	)
	t := time.Now()
	log.Info("начало задержки")
	err := doDelay(what, duration)
	log.Info("задержка окончена", "длительность", helpstr.FormatDuration(time.Since(t)), "error", err)
	return err

}

func doDelay(what string, duration time.Duration) error {

	var ctxDelay context.Context
	ctxDelay, skipDelayFunc = context.WithTimeout(ctxWork, duration)

	notify.Delay(log, types.DelayInfo{
		Run:     true,
		What:    what,
		Seconds: int(duration.Seconds()),
	})

	defer func() {
		notify.Delay(log, types.DelayInfo{Run: false})
	}()
	for {

		if ctxDelay.Err() != nil {
			return nil
		}

		if ctxWork.Err() != nil {
			return ctxWork.Err()
		}

		products := party.CheckedProducts()
		if len(products) == 0 {
			return merry.New("для опроса необходимо установить галочку для как минимум одиного прибора")
		}

		if len(products) == 0 {
			return merry.New("фоновый опрос: не выбрано ни одного прибора")
		}
	loopProducts:
		for _, p := range products {
			for _, v := range cfg.Get().Vars {
				_, err := readProductVarWithContext(p.Addr, v.Code, ctxDelay)

				if ctxDelay.Err() != nil {
					return nil
				}

				if ctxWork.Err() != nil {
					return ctxWork.Err()
				}

				if err == nil {
					continue
				}
				if merry.Is(err, context.DeadlineExceeded) {
					continue loopProducts
				}
				return err
			}
		}
	}
}
