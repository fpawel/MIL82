package data

import (
	"database/sql"
	"fmt"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/reform.v1"
)


func TestLastParty(t *testing.T) {
	conn, err := sql.Open("sqlite3", `:memory:`)
	if err != nil {
		t.Fatal(err)
	}
	db := reform.NewDB(conn, sqlite3.Dialect, nil)
	if _,err = db.Exec(SQLCreate); err != nil {
		t.Fatal(err)
	}
	party := Party{}
	if err = LastParty(db, &party); err != nil {
		t.Fatal(err)
	}
	fmt.Println(party.PartyID)
}
