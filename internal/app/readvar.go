package app

import (
	"context"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/gohelp"
	"github.com/fpawel/mil82/internal/api/notify"
	"github.com/fpawel/mil82/internal/api/types"
	"github.com/fpawel/mil82/internal/charts"
	"math"
	"math/rand"
	"time"
)

func readProductVar(addr modbus.Addr, VarCode modbus.Var) (float64, error) {
	return readProductVarWithContext(addr, VarCode, ctxWork)
}

func readProductVarWithContext(addr modbus.Addr, VarCode modbus.Var, ctx context.Context) (float64, error) {

	log := gohelp.LogWithKeys(log, "адрес", addr, "var", VarCode)

	if addr == 1 || addr == 2 || addr == 3 {

		timer := time.NewTimer(time.Millisecond * 200)
	pauseLoop:
		for {
			select {
			case <-timer.C:
				break pauseLoop
			case <-ctx.Done():
				timer.Stop()
				return 0, ctx.Err()
			}
		}

		value := math.Round(rand.Float64()*100) / 100
		notify.ReadVar(log, types.AddrVarValue{Addr: addr, VarCode: VarCode, Value: value})
		charts.AddPointToLastBucket(addr, VarCode, value)
		return value, nil
	}

	value, err := modbus.Read3BCD(log, responseReaderProducts(ctx), addr, VarCode)
	if err == nil {
		notify.ReadVar(log, types.AddrVarValue{Addr: addr, VarCode: VarCode, Value: value})
		charts.AddPointToLastBucket(addr, VarCode, value)
		return value, nil
	}

	notify.AddrError(log, types.AddrError{Addr: addr, Message: err.Error()})
	return value, err
}
