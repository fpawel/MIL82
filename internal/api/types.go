package api

import (
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/mil82/internal/data"
)

type LastPartyProduct struct {
	data.Product
	Place   int
	Checked bool
}

type TempPlusMinus struct {
	TempPlus, TempMinus float64
}

type AddrVarValue struct {
	Addr  modbus.Addr
	Var   modbus.Var
	Value float64
}
