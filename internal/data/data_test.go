package data

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"testing"

	_ "github.com/lib/pq"
)

func TestLastParty(t *testing.T) {
	conn, err := sql.Open("postgres", "host=localhost port=5432 user=mil82 password=1 dbname=mil82 sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	db := sqlx.NewDb(conn, "postgres")
	if _, err = db.Exec(SQLCreate); err != nil {
		t.Fatal(err)
	}
	party := Party{}
	if err = LastParty(db, &party); err != nil {
		t.Fatal(err)
	}
	party.Pgs1 += 1

	_, err = db.Exec(`UPDATE party SET pgs1 = $1 WHERE party_id = $2`, party.Pgs1, party.PartyID)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("party_id:", party.PartyID, "pgs1", party.Pgs1)

}
