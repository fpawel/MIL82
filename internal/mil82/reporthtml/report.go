package reporthtml

import (
	"database/sql"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/mil82/internal/data"
	"github.com/fpawel/mil82/internal/mil82"
)

type reportNode struct {
	Table *reportTable `json:",omitempty"`
	Tree  *reportNodes `json:",omitempty"`
}

type reportTable struct {
	Var   modbus.Var
	Gases []mil82.Gas `json:",omitempty"`
	Rows  []reportRow `json:",omitempty"`
}

type reportNodes struct {
	Title string
	Level int
	Nodes []reportNode `json:",omitempty"`
}

type reportRow struct {
	Addr   modbus.Addr
	Serial int
	Values []*float64
}

type workTempGases struct {
	work  mil82.Work
	temp  mil82.Temp
	gases []mil82.Gas
}

type reportDataNode struct {
	workTempGases
	title string
	nodes []reportDataNode
}

func reportParty(partyID int64) (result []reportNode) {
	var products []data.Product
	if err := data.DB.Select(&products, `SELECT product_id, addr, serial FROM product WHERE party_id = ?`, partyID); err != nil {
		panic(err)
	}
	for _, dataNd := range reportDataNodes {
		nd := reportNode{Tree: &reportNodes{Title: dataNd.title}}

		if len(dataNd.nodes) == 0 {
			dataNd.makeVarsReportNode(products, &nd)
		} else {

			for _, dataNd := range dataNd.nodes {
				nd2 := reportNode{Tree: &reportNodes{Title: dataNd.title, Level: 1}}
				dataNd.makeVarsReportNode(products, &nd2)
				if len(nd2.Tree.Nodes) > 0 {
					nd.Tree.Nodes = append(nd.Tree.Nodes, nd2)
				}
			}
		}
		if len(nd.Tree.Nodes) > 0 {
			result = append(result, nd)
		}
	}
	return
}

func (dataNd reportDataNode) makeVarsReportNode(products []data.Product, nd *reportNode) {
	nd.Table = nil
	nd.Tree.Title = dataNd.title

	hasValue := func(xs []*float64) bool {
		for _, x := range xs {
			if x != nil {
				return true
			}
		}
		return false
	}

	for _, Var := range mil82.Vars {
		varNd := reportNode{
			Table: &reportTable{
				//Title: fmt.Sprintf( "[%02d] %s",Var, varName[Var]),
				Var:   Var,
				Gases: dataNd.workTempGases.gases,
			},
		}
		for _, p := range products {
			varRow := reportRow{Addr: p.Addr, Serial: p.Serial, Values: make([]*float64, len(dataNd.gases))}
			for i, gas := range varNd.Table.Gases {
				var v float64
				err := data.DB.Get(&v,
					`SELECT value FROM product_value WHERE product_id = ? AND work=? AND temp=? AND var=? AND gas=?`,
					p.ProductID, dataNd.workTempGases.work, dataNd.workTempGases.temp, Var, gas)
				if err == sql.ErrNoRows {
					continue
				}
				if err != nil {
					panic(err)
				}
				varRow.Values[i] = &v
			}
			if hasValue(varRow.Values) {
				varNd.Table.Rows = append(varNd.Table.Rows, varRow)
			}
		}
		if len(varNd.Table.Rows) > 0 {
			nd.Tree.Nodes = append(nd.Tree.Nodes, varNd)
		}
	}
	return
}

var reportDataNodes = func() []reportDataNode {
	lf := func(title string, work mil82.Work, temp mil82.Temp, gases []mil82.Gas) reportDataNode {
		return reportDataNode{workTempGases{work, temp, gases}, title, nil}
	}
	nd := func(title string, nodes ...reportDataNode) reportDataNode {
		return reportDataNode{title: title, nodes: nodes}
	}
	gases1234 := []mil82.Gas{1, 2, 3, 4}
	gases134 := []mil82.Gas{1, 3, 4}

	return []reportDataNode{
		lf("Линеаризация", mil82.WorkLin, mil82.Temp20, gases1234),
		nd("Термокомпенсация",
			lf("Нормальная температура", mil82.WorkTemp, mil82.Temp20, gases1234),
			lf("Пониженная температура", mil82.WorkTemp, mil82.TempMinus, gases134),
			lf("Повышенная температура", mil82.WorkTemp, mil82.TempPlus, gases134),
			lf("+90⁰С", mil82.WorkTemp, mil82.Temp90, gases134)),
		nd("Проверка",
			lf("Нормальная температура", mil82.WorkCheckup, mil82.Temp20, gases1234),
			lf("Пониженная температура", mil82.WorkCheckup, mil82.TempMinus, gases134),
			lf("Повышенная температура", mil82.WorkCheckup, mil82.TempPlus, gases134)),
		lf("Техпрогон 1", mil82.WorkTex1, mil82.Temp20, gases1234),
		lf("Техпрогон 2", mil82.WorkTex2, mil82.Temp20, gases1234),
	}
}()
