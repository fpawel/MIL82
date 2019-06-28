package data

import (
	"database/sql"
)

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
