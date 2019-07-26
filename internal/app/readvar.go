package app

import (
	"context"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/gohelp"
	"github.com/fpawel/mil82/internal/api/notify"
	"github.com/fpawel/mil82/internal/api/types"
	"github.com/fpawel/mil82/internal/cfg"
	"github.com/fpawel/mil82/internal/charts"
	"time"
)

func readProductVar(addr modbus.Addr, VarCode modbus.Var) (float64, error) {
	return readProductVarWithContext(addr, VarCode, ctxWork)
}

func readProductVarWithContext(addr modbus.Addr, VarCode modbus.Var, ctx context.Context) (float64, error) {

	defer pauseWithContext(ctx.Done(), cfg.Get().InterrogateProductVarInterval())

	log := gohelp.LogPrependSuffixKeys(log, "адрес", addr, "var", VarCode)

	//if addr == 1 || addr == 2 || addr == 3 {
	//	return fakeReadAddrVarValue(ctx, addr, VarCode)
	//}

	value, err := modbus.Read3BCD(log, ctx, portProducts, addr, VarCode)
	if err == nil {
		notify.ReadVar(log, types.AddrVarValue{Addr: addr, VarCode: VarCode, Value: value})
		charts.AddPointToLastBucket(addr, VarCode, value)
		return value, nil
	}

	notify.AddrError(log, types.AddrError{Addr: addr, Message: err.Error()})

	return value, err
}

//func fakeReadAddrVarValue(ctx context.Context, addr modbus.Addr, Var modbus.Var) (float64, error){
//	value := math.Round(rand.Float64()*100) / 100
//	notify.ReadVar(log, types.AddrVarValue{Addr: addr, VarCode: Var, Value: value})
//	charts.AddPointToLastBucket(addr, Var, value)
//	return value, nil
//}

func pauseWithContext(chDone <-chan struct{}, d time.Duration) {
	timer := time.NewTimer(d)
	for {
		select {
		case <-timer.C:
			return
		case <-chDone:
			timer.Stop()
			return
		}
	}
}
