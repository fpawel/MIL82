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
	Cells []string
}

func PartyProductsValues(partyID int64, Var modbus.Var) (r Table) {

	type row []string

	r = Table{
		[]Row{
			{row{"ID"}},
			{row{"Сер.№"}},
			{row{"Адрес"}},
		},
	}

	var products []data.Product
	if err := data.DB.Select(&products, `SELECT product_id, addr, serial FROM product WHERE party_id = ?`, partyID); err != nil {
		panic(err)
	}

	addCell := func(n int, fmtStr string, v interface{}) {
		r.Rows[n].Cells = append(r.Rows[n].Cells, fmt.Sprintf(fmtStr, v))
	}
	for _, p := range products {
		addCell(0, "%d", p.ProductID)
		addCell(1, "%d", p.Serial)
		addCell(2, "%02d", p.Addr)
	}

	for _, dn := range dataTables {

		rows := []Row{{Cells: row{dn.title}}}

		for _, gas := range dn.gases {
			row := make([]string, len(products)+1)
			row[0] = fmt.Sprintf("ПГС%d", gas)
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
				row[i+1] = strconv.FormatFloat(v, 'f', -1, 64)
			}
			if hasValue {
				rows = append(rows, Row{row})
			}
		}

		if len(rows) > 1 {
			r.Rows = append(r.Rows, rows...)
		}
	}
	return
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
