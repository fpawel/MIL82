package data

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/powerman/must"
	"github.com/powerman/structlog"
	"os"
	"path/filepath"
)

//go:generate go run github.com/fpawel/gohelp/cmd/sqlstr/...

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
	if _, err := os.Stat(fileName); os.IsNotExist(err) || createNew {
		createNew = true
	}
	if createNew {
		must.Remove(fileName)
		log.Warn(fileName + " removed")
	}
	dbConn, err := sql.Open("sqlite3", fileName)
	if err != nil {
		panic(err)
	}
	dbConn.SetMaxIdleConns(1)
	dbConn.SetMaxOpenConns(1)
	dbConn.SetConnMaxLifetime(0)

	log.Info("open: " + fileName)
	DB = sqlx.NewDb(dbConn, "sqlite3")
	DB.MustExec(SQLCreate)
}
