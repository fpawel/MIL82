package mil82

import (
	"database/sql"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/mil82/internal/data"
)

type reportNode struct {
	Table *reportTable `json:",omitempty"`
	Tree  *reportNodes `json:",omitempty"`
}

type reportTable struct {
	Var   modbus.Var
	Gases []Gas       `json:",omitempty"`
	Rows  []reportRow `json:",omitempty"`
}

type reportNodes struct {
	Title  string
	HtmlID string
	Level  int
	Nodes  []reportNode `json:",omitempty"`
}

type reportRow struct {
	Addr   modbus.Addr
	Serial int
	Values []*float64
}

type workTempGases struct {
	work  Work
	temp  Temp
	gases []Gas
}

type reportDataNode struct {
	workTempGases
	title string
	nodes []reportDataNode
}

func reportParty(partyID int64) (result []reportNode) {
	var products []data.Product
	if err := data.DB.Select(&products, `SELECT * FROM product WHERE party_id = ?`, partyID); err != nil {
		panic(err)
	}
	for _, dataNd := range reportDataNodes() {
		nd := reportNode{Tree: &reportNodes{Title: dataNd.title, HtmlID: string(dataNd.work)}}

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
	nd.Tree.HtmlID = string(dataNd.work) + "_" + string(dataNd.temp)

	hasValue := func(xs []*float64) bool {
		for _, x := range xs {
			if x != nil {
				return true
			}
		}
		return false
	}

	for _, Var := range Vars {
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

var varName = func() map[modbus.Var]string {
	type Var struct {
		Name string     `db:"name"`
		Var  modbus.Var `db:"var"`
	}
	var vars []Var
	if err := data.DB.Select(&vars, `SELECT * FROM var`); err != nil {
		panic(err)
	}

	m := make(map[modbus.Var]string)
	for _, v := range vars {
		m[v.Var] = v.Name
	}
	return m
}()

func reportDataNodes() []reportDataNode {
	lf := func(title string, work Work, temp Temp, gases []Gas) reportDataNode {
		return reportDataNode{workTempGases{work, temp, gases}, title, nil}
	}
	nd := func(title string, nodes ...reportDataNode) reportDataNode {
		return reportDataNode{title: title, nodes: nodes}
	}
	gases1234 := []Gas{1, 2, 3, 4}
	gases134 := []Gas{1, 3, 4}

	return []reportDataNode{
		lf("Линеаризация", WorkLin, Temp20, gases1234),
		nd("Термокомпенсация",
			lf("Нормальная температура", WorkTemp, Temp20, gases1234),
			lf("Пониженная температура", WorkTemp, TempMinus, gases134),
			lf("Повышенная температура", WorkTemp, TempPlus, gases134),
			lf("+90⁰С", WorkTemp, Temp90, gases134)),
		nd("Проверка",
			lf("Нормальная температура", WorkCheckup, Temp20, gases1234),
			lf("Пониженная температура", WorkCheckup, TempMinus, gases134),
			lf("Повышенная температура", WorkCheckup, TempPlus, gases134)),
		lf("Техпрогон 1", WorkTex1, Temp20, gases1234),
		lf("Техпрогон 2", WorkTex2, Temp20, gases1234),
	}
}
