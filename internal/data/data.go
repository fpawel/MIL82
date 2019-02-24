package data

import (
	"github.com/jmoiron/sqlx"
	"gopkg.in/reform.v1"
)

//go:generate go run github.com/fpawel/elco/cmd/utils/sqlstr/...

//reform-db -db-driver=postgres -db-source="host=localhost port=5432 user=mil82 dbname=mil82 password=1 sslmode=disable" init

func LastParty(db *sqlx.DB, party *Party) error {
	err := db.Get(party, `SELECT * FROM party ORDER BY created_at DESC LIMIT 1;`)
	if err == reform.ErrNoRows {
		return db.Get(party, `INSERT INTO party DEFAULT VALUES; SELECT * FROM party ORDER BY created_at DESC LIMIT 1`)
	}
	return err
}
