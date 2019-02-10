package api

import (
	"github.com/fpawel/elco/internal/crud"
	"github.com/fpawel/elco/internal/data"
)

type Works struct {
	c crud.Works
}

func NewWorks(c crud.Works) *Works {
	return &Works{c}
}

func (x *Works) Years(_ struct{}, years *[]int) (err error) {
	*years, err = x.c.Years()
	return
}

func (x *Works) Months(r struct{ Year int }, months *[]int) (err error) {
	*months, err = x.c.Months(r.Year)
	return
}

func (x *Works) Days(r struct{ Year, Month int }, days *[]int) (err error) {
	*days, err = x.c.Days(r.Year, r.Month)
	return
}

func (x *Works) DayWorks(r struct{ Year, Month, Day int }, works *[]data.WorkInfo) (err error) {
	*works, err = x.c.DayWorks(r.Year, r.Month, r.Day)
	return
}

func (x *Works) WorkEntries(workID [1]int64, r *[]data.Entry) (err error) {
	*r, err = x.c.WorkEntries(workID[0])
	return
}

func (x *Works) DayEntries(d [3]int, entries *[]data.Entry) (err error) {
	*entries, err = x.c.DayEntries(d[0], d[1], d[2])
	return
}