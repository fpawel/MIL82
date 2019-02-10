package crud

import (
	"database/sql"
	"github.com/fpawel/mil82/internal/data"
	"gopkg.in/reform.v1"
)

type Works struct {
	dbContext
}

func (x Works) Years() (years []int, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	err = x.dbx.Select(&years, `SELECT DISTINCT year FROM work_info ORDER BY year ASC;`)
	return
}

func (x Works) Months(y int) (months []int, err error ) {
	x.mu.Lock()
	defer x.mu.Unlock()
	err = x.dbx.Select( &months,
		`SELECT DISTINCT month FROM work_info WHERE year = ? ORDER BY month ASC;`, y)
	return
}

func (x Works) Days(year, month int) (days []int, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	err = x.dbx.Select( &days,
		`SELECT DISTINCT day FROM work_info WHERE year = ? AND month = ? ORDER BY day ASC;`,
		year, month)
	return
}

func (x Works) DayWorks(year, month, day int) ([]data.WorkInfo, error ) {
	x.mu.Lock()
	defer x.mu.Unlock()

	rows, err := x.dbr.SelectRows(data.WorkInfoTable,
		"WHERE year = ? AND month = ? AND day = ?", year, month, day)
	if err != nil {
		return nil, err
	}
	var works []data.WorkInfo
	for {
		var work data.WorkInfo
		err = x.dbr.NextRow(&work, rows)
		if err == reform.ErrNoRows {
			break
		}
		if err != nil {
			return nil, err
		}
		works = append(works, work)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return works, nil
}

func (x Works) WorkEntries(workID int64) ([]data.Entry, error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	rows, err := x.dbr.SelectRows(data.EntryTable, "WHERE work_id = ?", workID)
	if err != nil {
		return nil, err
	}
	return x.getEntries(rows)
}

func (x Works) DayEntries(year, month, day int ) ([]data.Entry, error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	rows, err := x.dbr.SelectRows(data.EntryTable, `
WHERE CAST(STRFTIME('%Y', DATETIME(created_at, '+3 hours')) AS INTEGER) = ? AND
      CAST(STRFTIME('%m', DATETIME(created_at, '+3 hours')) AS INTEGER) = ? AND
      CAST(STRFTIME('%d', DATETIME(created_at, '+3 hours')) AS INTEGER) = ?
`, year, month, day)
	if err != nil {
		return nil, err
	}
	return x.getEntries(rows)
}

func (x Works) getEntries(rows *sql.Rows ) ([]data.Entry, error) {
	var entries []data.Entry
	for {
		var entry data.Entry
		err := x.dbr.NextRow(&entry, rows)
		if err == reform.ErrNoRows {
			break
		}
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	return entries,nil
}