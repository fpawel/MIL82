package api

import (
	"github.com/fpawel/mil82/internal/charts"
)

type ChartsSvc struct{}

type ChartsYearMonth struct {
	Year  int `db:"year"`
	Month int `db:"month"`
}

func (_ *ChartsSvc) YearsMonths(_ struct{}, r *[]ChartsYearMonth) error {
	if err := charts.DB.Select(r, `SELECT DISTINCT year, month FROM bucket_time ORDER BY year DESC, month DESC`); err != nil {
		panic(err)
	}
	return nil
}

type ChartsBucket struct {
	Day      int    `db:"day"`
	Hour     int    `db:"hour"`
	Minute   int    `db:"minute"`
	BucketID int64  `db:"bucket_id"`
	Name     string `db:"name"`
	Last     bool   `db:"last"`
}

func (_ *ChartsSvc) BucketsOfYearMonth(x ChartsYearMonth, r *[]ChartsBucket) error {
	if err := charts.DB.Select(r, `
SELECT day, hour, minute, bucket_id, name, bucket_id = (SELECT bucket_id FROM last_bucket) AS last
FROM bucket_time
WHERE year = ?
  AND month = ?
ORDER BY created_at`, x.Year, x.Month); err != nil {
		panic(err)
	}
	return nil
}
