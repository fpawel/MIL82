package mil82

import (
	"database/sql"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/mil82/internal/data"
)

type ReportNode struct {
	Title, Html string
	Nodes       []ReportNode
}

type reportNode struct {
	Table *reportTable `json:",omitempty"`
	Tree  *reportNodes `json:",omitempty"`
}

type reportTable struct {
	Title string
	Gases []Gas       `json:",omitempty"`
	Rows  []reportRow `json:",omitempty"`
}

type reportNodes struct {
	Title string
	Nodes []reportNode `json:",omitempty"`
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

func ReportParty(partyID int64) (result []ReportNode) {
	r := reportParty(partyID)
	for _, nd := range r {
		result = append(result, nd.ReportNode(""))
	}
	return
}

func reportParty(partyID int64) (result []reportNode) {
	var products []data.Product
	if err := data.DB.Select(&products, `SELECT * FROM product WHERE party_id = ?`, partyID); err != nil {
		panic(err)
	}
	for _, dataNd := range reportDataNodes() {
		nd := reportNode{Tree: &reportNodes{Title: dataNd.title}}

		if len(dataNd.nodes) == 0 {
			nd = dataNd.makeVarsReportNode(products)
		} else {

			for _, dataNd := range dataNd.nodes {
				nd2 := dataNd.makeVarsReportNode(products)
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

func (dataNd reportDataNode) makeVarsReportNode(products []data.Product) (nd reportNode) {
	nd.Tree = &reportNodes{Title: dataNd.title}
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
				Title: varName[Var],
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

func (nd reportNode) Tables(titlePrefix string) (result []reportTable) {
	if nd.Table != nil {
		result = []reportTable{*nd.Table}
	} else if nd.Tree != nil {
		for i := range nd.Tree.Nodes {
			result = nd.Tree.Nodes[i].Tables(nd.Tree.Title + ". ")
		}
	}
	for i := range result {
		result[i].Title = titlePrefix + result[i].Title
	}
	return
}

func (nd reportNode) Title() string {
	if nd.Tree != nil {
		return nd.Tree.Title
	} else if nd.Table != nil {
		return nd.Table.Title
	}
	panic("not a node")
}

func (nd reportNode) ReportNode(titlePrefix string) (r ReportNode) {
	r.Title = nd.Title()
	r.Html = nd.View(titlePrefix)
	r.Nodes = []ReportNode{}
	if nd.Tree != nil {
		for _, nd1 := range nd.Tree.Nodes {
			r.Nodes = append(r.Nodes, nd1.ReportNode(nd.Title()+". "))
		}
	}
	return
}

func reportNodesTables(nds []reportNode) (result []reportTable) {
	for _, nd := range nds {
		result = append(result, nd.Tables("")...)
	}
	return
}
