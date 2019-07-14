package types

import (
	"fmt"
	"github.com/fpawel/comm/modbus"
)

type TempPlusMinus struct {
	TempPlus, TempMinus float64
}

type AddrVarValue struct {
	Addr    modbus.Addr
	VarCode modbus.Var
	Value   float64
}

type AddrError struct {
	Addr    modbus.Addr
	Message string
}

type WorkResult int

const (
	WrOk WorkResult = iota
	WrCanceled
	WrError
)

func (x WorkResult) String() string {
	switch x {
	case WrOk:
		return "Ok"
	case WrCanceled:
		return "Canceled"
	case WrError:
		return "Error"
	default:
		return fmt.Sprintf("?%d", x)
	}
}

type WorkResultInfo struct {
	Work    string
	Result  WorkResult
	Message string
}

type DelayInfo struct {
	Run     bool
	Seconds int
	What    string
}
