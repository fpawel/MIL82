package app

import (
	"context"
	"fmt"
	"github.com/ansel1/merry"
	"github.com/fpawel/gohelp"
	"github.com/fpawel/gohelp/myfmt"
	"github.com/fpawel/mil82/internal/api/notify"
	"github.com/fpawel/mil82/internal/api/types"
	"github.com/fpawel/mil82/internal/cfg"
	"github.com/fpawel/mil82/internal/last_party"
	"time"
)

func delayf(x worker, duration time.Duration, format string, a ...interface{}) error {
	return delay(x, duration, fmt.Sprintf(format, a...))
}

func delay(x worker, duration time.Duration, name string) error {
	fd := myfmt.FormatDuration
	startTime := time.Now()
	x.log = gohelp.LogPrependSuffixKeys(x.log, "start", startTime.Format("15:04:05"))

	{
		var skipDelay context.CancelFunc
		x.ctx, skipDelay = context.WithTimeout(x.ctx, duration)
		skipDelayFunc = func() {
			skipDelay()
			x.log.Info("задержка прервана", "elapsed", myfmt.FormatDuration(time.Since(startTime)))
		}
	}

	ctxWork := x.ctx
	return x.performf("%s: %s", name, fd(duration))(func(x worker) error {
		x.log.Info("задержка начата")
		defer func() {
			x.log.Debug("задержка окончена", "elapsed", fd(time.Since(startTime)))
			notify.EndDelay(x.log, "")
		}()
		for {
			products := last_party.CheckedProducts()
			if len(products) == 0 {
				return merry.New("для опроса необходимо установить галочку для как минимум одиного прибора")
			}

			for _, p := range products {
				for _, v := range cfg.Get().Vars {
					_, err := readProductVar(x, p.Addr, v.Code)
					notify.Delay(nil, types.DelayInfo{
						What:           name,
						TotalSeconds:   int(duration.Seconds()),
						ElapsedSeconds: int(time.Since(startTime).Seconds()),
					})
					if merry.Is(err, context.DeadlineExceeded) {
						return nil // задержка истекла
					}
					if merry.Is(err, context.Canceled) {
						if x.ctx.Err() == context.Canceled {
							return nil // задержка пропущена пользователем
						}
						if ctxWork.Err() == context.Canceled {
							return context.Canceled // прервано пользователем
						}
						return nil
					}
					pause(x.ctx.Done(), millis(cfg.Get().InterrogateProductVarIntervalMillis))
				}
			}
		}
	})
}

//func delay(what string, duration time.Duration) error {
//
//	originalLog := log
//	defer func() {
//		log = originalLog
//	}()
//	log = gohelp.LogPrependSuffixKeys(log,
//		"фоновый_опрос", what,
//		"общая_длительность_задержки", myfmt.FormatDuration(duration),
//	)
//	t := time.Now()
//	log.Info("начало задержки")
//	err := doDelay(what, duration)
//	log.Info("задержка окончена", "длительность", myfmt.FormatDuration(time.Since(t)), "error", err)
//	return err
//
//}
//
//func doDelay(what string, duration time.Duration) error {
//
//	var ctxDelay context.Context
//	ctxDelay, skipDelayFunc = context.WithTimeout(ctxWork, duration)
//
//	notify.Delay(log, types.DelayInfo{
//		Run:     true,
//		What:    what,
//		Seconds: int(duration.Seconds()),
//	})
//
//	defer func() {
//		notify.Delay(log, types.DelayInfo{Run: false})
//	}()
//	for {
//
//		if ctxDelay.Err() != nil {
//			return nil
//		}
//
//		if ctxWork.Err() != nil {
//			return ctxWork.Err()
//		}
//
//		products := last_party.CheckedProducts()
//		if len(products) == 0 {
//			return merry.New("для опроса необходимо установить галочку для как минимум одиного прибора")
//		}
//
//		if len(products) == 0 {
//			return merry.New("фоновый опрос: не выбрано ни одного прибора")
//		}
//	loopProducts:
//		for _, p := range products {
//			for _, v := range cfg.Get().Vars {
//				_, err := readProductVarWithContext(p.Addr, v.Code, ctxDelay)
//
//				if ctxDelay.Err() != nil {
//					return nil
//				}
//
//				if ctxWork.Err() != nil {
//					return ctxWork.Err()
//				}
//
//				if err == nil {
//					continue
//				}
//				if merry.Is(err, context.DeadlineExceeded) {
//					continue loopProducts
//				}
//				return err
//			}
//		}
//	}
//}
