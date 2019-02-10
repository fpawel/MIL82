package data

import (
	"gopkg.in/reform.v1"
)

//go:generate go run github.com/fpawel/goutils/dbutils/sqlstr/...

func Products(db *reform.DB, partyID int64, products *[]Product) error {
	rows, err := db.SelectRows(ProductTable, "WHERE party_id = ? ORDER BY place",
		partyID)
	if err != nil {
		return err
	}
	defer func() {
		_ = rows.Close()
	}()

	for {
		var product Product
		err = db.NextRow(&product, rows)
		if err == reform.ErrNoRows {
			break
		}
		if err != nil {
			return err
		}
		*products = append(*products, product)
	}
	return nil
}

func  LastParty(db *reform.DB, party *Party) error {
	err :=db.SelectOneTo(party, `ORDER BY created_at DESC LIMIT 1;`)
	if err == reform.ErrNoRows{
		*party = Party{}
		return db.Insert(party)
	}
	return err
}