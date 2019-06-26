package data

import (
	"database/sql"
)

func LastParty() (party Party) {
	err := DB.Get(&party, `SELECT party_id FROM last_party`)
	if err == sql.ErrNoRows {
		DB.MustExec(`INSERT INTO party DEFAULT VALUES`)
		err = DB.Get(&party, `SELECT party_id FROM last_party`)
	}
	if err != nil {
		panic(err)
	}
	return
}

func LastPartyProducts() (products []Product) {
	err := DB.Select(&products,
		`SELECT * FROM product WHERE party_id = (SELECT party_id FROM last_party)`)
	if err != nil {
		panic(err)
	}
	return
}
