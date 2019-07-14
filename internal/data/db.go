package data

import (
	"database/sql"
	"github.com/fpawel/gohelp"
	"github.com/fpawel/mil82/internal"
	"path/filepath"
)

//go:generate go run github.com/fpawel/gohelp/cmd/sqlstr/...

var (
	DB = gohelp.OpenSqliteDBx(filepath.Join(internal.DataDir(), "mil82.sqlite"))
)

func init() {
	DB.MustExec(SQLCreate)
}

func LastParty() (party Party) {
	err := DB.Get(&party, `SELECT * FROM last_party`)
	if err == sql.ErrNoRows {
		DB.MustExec(`INSERT INTO party DEFAULT VALUES`)
		err = DB.Get(&party, `SELECT * FROM last_party`)
	}
	if err != nil {
		panic(err)
	}
	return
}
