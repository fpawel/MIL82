package crud

import (
	"github.com/fpawel/mil82/internal/data"
	"github.com/jmoiron/sqlx"
	"gopkg.in/reform.v1"
	"sync"
)

type DBContext struct {
	dbContext
}

type dbContext struct {
	mu  *sync.Mutex
	dbx *sqlx.DB
	dbr *reform.DB
}

func (x dbContext) LastParty(party *data.Party) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	exists := false
	if err :=x.dbx.Get(&exists, `SELECT EXISTS(*) FROM party`); err != nil {
		return err
	}
	if exists {
		return x.dbr.SelectOneTo(party, `ORDER BY created_at DESC LIMIT 1;`)
	}
	r, err :=x.dbx.Exec(`INSERT INTO party DEFAULT VALUES`)
	if err != nil {
		return err
	}
	partyID,err := r.LastInsertId()
	if err != nil {
		return err
	}
	return x.dbr.FindByPrimaryKeyTo(party, partyID)
}

func (x dbContext) Party(party *data.Party) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.dbr.FindByPrimaryKeyTo(party, party.PartyID)
}

func (x dbContext) Products(partyID int64, products *[]data.Product) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	return data.Products(x.dbr, partyID, products)
}

func (x dbContext) SaveParty(party *data.Party) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.dbr.Save(party)
}

func (x dbContext) SaveProduct(product *data.Product) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.dbr.Save(product)
}

func (x dbContext) SavePartyProducts(party *data.Party, products *[]data.Product) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.dbr.Save(product)
}