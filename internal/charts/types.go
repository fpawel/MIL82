package charts

import (
	"github.com/fpawel/comm/modbus"
	"time"
)

type Bucket struct {
	BucketID  int64      `db:"bucket_id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	Year      int        `db:"year"`
	Month     time.Month `db:"month"`
	Day       int        `db:"day"`
}

type point struct {
	StoredAt time.Time
	Var      modbus.Var
	Addr     modbus.Addr
	Value    float64
}

type Point struct {
	Addr        modbus.Addr `db:"addr"`
	Var         modbus.Var  `db:"var"`
	Value       float64     `db:"value"`
	Year        int         `db:"year"`
	Month       time.Month  `db:"month"`
	Day         int         `db:"day"`
	Hour        int         `db:"hour"`
	Minute      int         `db:"minute"`
	Second      int         `db:"second"`
	Millisecond int         `db:"millisecond"`
}

func (x point) Point() Point {
	t := x.StoredAt

	return Point{
		Addr:        x.Addr,
		Var:         x.Var,
		Value:       x.Value,
		Year:        t.Year(),
		Month:       t.Month(),
		Day:         t.Day(),
		Hour:        t.Hour(),
		Minute:      t.Minute(),
		Second:      t.Second(),
		Millisecond: t.Nanosecond() / 1000000,
	}
}
