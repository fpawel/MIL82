package report

import (
	"database/sql"
	"fmt"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/mil82/internal/data"
	"github.com/fpawel/mil82/internal/mil82"
	"strconv"
)

type Table struct {
	Rows []Row
}

type Row struct {
	Cells []Cell
}

type Cell struct {
	ValueType    ValueType
	Text, Detail string
}

type ValueType int

const (
	vtNone ValueType = iota
	vtOk
	vtError
)

func PartyProductsValues(partyID int64, Var modbus.Var) (r Table) {
	type row []Cell
	cell1 := func(s string) Cell {
		return Cell{Text: s}
	}
	cell1f := func(format string, a ...interface{}) Cell {
		return cell1(fmt.Sprintf(format, a...))
	}

	row1 := func(s string) Row {
		return Row{row{cell1(s)}}
	}
	r = Table{
		[]Row{
			row1("ID"),
			row1("Сер.№"),
			row1("Адрес"),
		},
	}
	var products []data.Product
	if err := data.DB.Select(&products, `SELECT product_id, addr, serial FROM product WHERE party_id = ?`, partyID); err != nil {
		panic(err)
	}

	addCell := func(n int, fmtStr string, v interface{}) {
		r.Rows[n].Cells = append(r.Rows[n].Cells, cell1f(fmtStr, v))
	}
	for _, p := range products {
		addCell(0, "%d", p.ProductID)
		addCell(1, "%d", p.Serial)
		addCell(2, "%02d", p.Addr)
	}

	for _, dn := range dataTables {

		rows := []Row{row1(dn.title)}

		for _, gas := range dn.gases {
			row := Row{Cells: make([]Cell, len(products)+1)}
			row.Cells[0] = cell1f("ПГС%d", gas)
			hasValue := false
			for i, p := range products {
				var v float64
				err := data.DB.Get(&v,
					`SELECT value FROM product_value WHERE product_id = ? AND work=? AND temp=? AND var=? AND gas=?`,
					p.ProductID, dn.work, dn.temp, Var, gas)
				if err == sql.ErrNoRows {
					continue
				}
				if err != nil {
					panic(err)
				}
				hasValue = true
				row.Cells[i+1] = cell1(strconv.FormatFloat(v, 'f', -1, 64))
			}
			if hasValue {
				rows = append(rows, row)
			}
		}

		if len(rows) > 1 {
			r.Rows = append(r.Rows, rows...)
		}
	}
	return
}

func errors(partyID int64) {
	party := data.GetParty(partyID)
	productType := mil82.ProductTypeByName(party.ProductType)

	maxErr20 := func(Cn float64) float64 {
		if productType.Component != mil82.CO2 {
			return 2.5 + 0.05*Cn
		}
		switch productType.Scale {
		case 4:
			return 0.2 + 0.05*Cn
		case 10:
			return 0.5
		default:
			return 1
		}
	}
}

type dataTable struct {
	title string
	work  mil82.Work
	temp  mil82.Temp
	gases []mil82.Gas
}

var dataTables = func() []dataTable {
	gases1234 := []mil82.Gas{1, 2, 3, 4}
	gases134 := []mil82.Gas{1, 3, 4}

	return []dataTable{
		{"Линеаризация", mil82.WorkLin, mil82.Temp20, gases1234},
		{"+20⁰С, НКУ, ", mil82.WorkTemp, mil82.Temp20, gases1234},
		{"«-»⁰С, пониженная температура", mil82.WorkTemp, mil82.TempMinus, gases134},
		{"«+»⁰С, повышенная температура", mil82.WorkTemp, mil82.TempPlus, gases134},
		{"+90⁰С", mil82.WorkTemp, mil82.Temp90, gases134},
		{"Проверка +20⁰С, НКУ", mil82.WorkCheckup, mil82.Temp20, gases1234},
		{"Проверка «-»⁰С, пониженная температура", mil82.WorkCheckup, mil82.TempMinus, gases134},
		{"Проверка «+»⁰С, повышенная температура", mil82.WorkCheckup, mil82.TempPlus, gases134},
		{"1. Первый техпрогон", mil82.WorkTex1, mil82.Temp20, gases134},
		{"2. Второй техпрогон", mil82.WorkTex2, mil82.Temp20, gases134},
	}
}()

func errorLimit(t int) float64 {

	//let concErrorlimit (t:ProductType) concValue =
	//	let scale = t.Scale
	//if t.IsCH then 2.5m+0.05m * concValue
	//elif scale=Sc4 then 0.2m + 0.05m * concValue
	//elif scale=Sc10 then 0.5m
	//elif scale=Sc20 then 1.0m else 0.m
}
