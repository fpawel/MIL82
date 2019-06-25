package data

import (
	"database/sql"
)

func GetLastPartyID() (partyID int64) {
	row := DB.QueryRow(`SELECT party_id FROM party ORDER BY created_at DESC LIMIT 1`)

	err := row.Scan(&partyID)
	if err == sql.ErrNoRows {
		DB.MustExec(`INSERT INTO party DEFAULT VALUES`)
		row = DB.QueryRow(`SELECT party_id FROM party ORDER BY created_at DESC LIMIT 1`)
		err = row.Scan(&partyID)
	}
	if err != nil {
		panic(err)
	}
	return partyID
}
