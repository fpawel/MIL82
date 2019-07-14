package charts

import (
	"encoding/binary"
	"fmt"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/gohelp"
	"github.com/fpawel/mil82/internal"
	"io"
	"path/filepath"
	"sync"
	"time"
)

//go:generate go run github.com/fpawel/goutils/dbutils/sqlstr/...

func LastBucket() (buck Bucket) {
	if err := DB.Get(&buck, `SELECT bucket_id, name, created_at, updated_at FROM last_bucket`); err != nil {
		panic(err)
	}
	return
}

func CreateNewBucket(name string) {
	DB.MustExec(`DELETE FROM bucket WHERE created_at = updated_at; INSERT INTO bucket (name) VALUES (?)`, name)
}

// AddPoint - добавить новую точку в кеш.
func AddPointToLastBucket(addr modbus.Addr, v modbus.Var, value float64) {
	muPoints.Lock()
	defer muPoints.Unlock()
	currentPoints = append(currentPoints, point{
		StoredAt: time.Now(),
		Addr:     addr,
		Var:      v,
		Value:    value,
	})
	if time.Since(currentPoints[0].StoredAt) > time.Minute {
		saveLastBucket()
	}
}

// SaveLastBucket - сохранить точки из кеша, очистить кеш.
// Можно вызывать кокурентно.
func SaveLastBucket() {
	muPoints.Lock()
	defer muPoints.Unlock()
	saveLastBucket()
}

func saveLastBucket() {
	if len(currentPoints) == 0 {
		return
	}
	queryInsertPoints := queryInsertPoints()
	currentPoints = nil
	go DB.MustExec(queryInsertPoints)
}

func WritePointsResponse(w io.Writer, bucketID int64) {

	var points []Point

	if err := DB.Select(&points, `
SELECT addr, var, value, year, month, day, hour, minute, second, millisecond 
FROM series_time 
WHERE bucket_id = ?`, bucketID); err != nil {
		panic(err)
	}

	if LastBucket().BucketID == bucketID {
		var points3 []Point
		muPoints.Lock()
		for _, p := range currentPoints {
			points3 = append(points3, p.Point())
		}
		muPoints.Unlock()
		points = append(points3, points...)
	}

	write := func(n interface{}) {
		if err := binary.Write(w, binary.LittleEndian, n); err != nil {
			panic(err)
		}
	}
	write(uint64(len(points)))
	for _, x := range points {
		write(byte(x.Addr))
		write(uint16(x.Var))
		write(uint16(x.Year))
		write(byte(x.Month))
		write(byte(x.Day))
		write(byte(x.Hour))
		write(byte(x.Minute))
		write(byte(x.Second))
		write(uint16(x.Millisecond))
		write(float64(x.Value))
	}
}

func queryInsertPoints() string {
	queryStr := `INSERT INTO series(bucket_id, Addr, var, Value, stored_at)  VALUES `
	for i, a := range currentPoints {

		s := fmt.Sprintf("(%d, %d, %d, %v,", LastBucket().BucketID, a.Addr, a.Var, a.Value) +
			"julianday(STRFTIME('%Y-%m-%d %H:%M:%f','" +
			a.StoredAt.Format("2006-01-02 15:04:05.000") + "')))"

		if i < len(currentPoints)-1 {
			s += ", "
		}
		queryStr += s
	}
	return queryStr
}

var (
	DB            = gohelp.OpenSqliteDBx(filepath.Join(internal.DataDir(), "mil82.series.sqlite"))
	currentPoints []point
	muPoints      sync.Mutex
)

func init() {
	DB.MustExec(SQLCreate)
}
