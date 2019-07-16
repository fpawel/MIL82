package mil82

import (
	"github.com/fpawel/comm/modbus"
)

type Work string

type Gas int

type Temp string

const (
	WorkLin     Work = "work_lin"
	WorkTemp    Work = "work_temp"
	WorkCheckup Work = "work_checkup"
	WorkTex1    Work = "work_tex1"
	WorkTex2    Work = "work_tex2"

	VarConc modbus.Var = 0
	VarTemp modbus.Var = 2
	VarCurr modbus.Var = 4
	Var8    modbus.Var = 8
	Var10   modbus.Var = 10
	VarWork modbus.Var = 12
	VarRef  modbus.Var = 14
	Var16   modbus.Var = 16

	Gas1 Gas = 1
	Gas2 Gas = 2
	Gas3 Gas = 3
	Gas4 Gas = 4

	TempMinus Temp = "temp_minus"
	Temp20    Temp = "temp_20"
	TempPlus  Temp = "temp_plus"
	Temp90    Temp = "temp_90"
)

var Vars = []modbus.Var{
	VarConc, VarTemp, Var16, VarCurr, Var8, Var10, VarWork, VarRef,
}
