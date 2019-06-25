package data

import (
	"database/sql"
	"github.com/ansel1/merry"
	"github.com/jmoiron/sqlx"
	"github.com/powerman/structlog"
	"os"
	"path/filepath"
)

var (
	DB  *sqlx.DB
	log = structlog.New()
)

func Open(createNew bool) {
	dir := os.Getenv("MIL82_DATA_DIR")
	if len(dir) == 0 {
		dir = filepath.Dir(os.Args[0])
	}
	fileName := filepath.Join(dir, "mil82.sqlite")
	if createNew {
		if _, err := os.Stat(fileName); err == nil {
			if err := os.Remove(fileName); err != nil {
				panic(merry.Appendf(err, "unable to delete database file: %s", fileName))
			}
		}
	}
	dbConn, err := sql.Open("sqlite3", fileName)
	if err != nil {
		panic(err)
	}
	dbConn.SetMaxIdleConns(1)
	dbConn.SetMaxOpenConns(1)
	dbConn.SetConnMaxLifetime(0)

	if _, err = dbConn.Exec(SQLCreate); err != nil {
		panic(err)
	}
	log.Info("open: " + fileName)
	DB = sqlx.NewDb(dbConn, "sqlite3")
}
