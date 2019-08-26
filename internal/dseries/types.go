package dseries

import (
	"github.com/fpawel/comm/modbus"
	"time"
)

type ChartsBucket struct {
	CreatedAt, UpdatedAt TimeDelphi
	BucketID             int64
	Name                 string
	IsLast               bool
}

type TimeDelphi struct {
	Year        int        `db:"year"`
	Month       time.Month `db:"month"`
	Day         int        `db:"day"`
	Hour        int        `db:"hour"`
	Minute      int        `db:"minute"`
	Second      int        `db:"second"`
	Millisecond int        `db:"millisecond"`
}

func (x TimeDelphi) Time() time.Time {
	return time.Date(
		x.Year, x.Month, x.Day,
		x.Hour, x.Minute, x.Second,
		x.Millisecond*int(time.Millisecond/time.Nanosecond),
		time.Local)
}

func timeDelphi(t time.Time) TimeDelphi {
	return TimeDelphi{
		Year:        t.Year(),
		Month:       t.Month(),
		Day:         t.Day(),
		Hour:        t.Hour(),
		Minute:      t.Minute(),
		Second:      t.Second(),
		Millisecond: t.Nanosecond() / 1000000,
	}
}

type bucket struct {
	BucketID        int64     `db:"bucket_id"`
	Name            string    `db:"name"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	CreatedAtYear   int       `db:"created_at_year"`
	CreatedAtMonth  int       `db:"created_at_month"`
	CreatedAtDay    int       `db:"created_at_day"`
	CreatedAtHour   int       `db:"created_at_hour"`
	CreatedAtMinute int       `db:"created_at_minute"`
	UpdatedAtYear   int       `db:"updated_at_year"`
	UpdatedAtMonth  int       `db:"updated_at_month"`
	UpdatedAtDay    int       `db:"updated_at_day"`
	UpdatedAtHour   int       `db:"updated_at_hour"`
	UpdatedAtMinute int       `db:"updated_at_minute"`
	IsLast          bool      `db:"is_last"`
}

type point struct {
	StoredAt time.Time
	Var      modbus.Var
	Addr     modbus.Addr
	Value    float64
}

type point3 struct {
	TimeDelphi
	Addr  modbus.Addr `db:"addr"`
	Var   modbus.Var  `db:"var"`
	Value float64     `db:"value"`
}

func (x point) point3() point3 {
	return point3{
		TimeDelphi: timeDelphi(x.StoredAt),
		Addr:       x.Addr,
		Var:        x.Var,
		Value:      x.Value,
	}
}
