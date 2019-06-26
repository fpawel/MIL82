package api

import (
	"github.com/fpawel/mil82/internal/cfg"
	"github.com/fpawel/mil82/internal/data"
)

type LastPartySvc struct{}

type LastPartyProduct struct {
	data.Product
	Place   int    `db:"-"`
	Checked bool   `db:"-"`
	Error   string `db:"-"`
}

func (_ *LastPartySvc) Party(_ struct{}, r *data.Party) error {
	*r = data.LastParty()
	return nil
}

func (_ *LastPartySvc) DeleteProduct(productID [1]int64, r *[]LastPartyProduct) error {
	data.DB.MustExec(`DELETE FROM product WHERE product_id = ?`, productID[0])
	return getLastPartyProducts2(r)
}

func (_ *LastPartySvc) AddNewProduct(_ struct{}, r *[]LastPartyProduct) error {
	var products []LastPartyProduct
	if err := getLastPartyProducts(&products); err != nil {
		return err
	}
	addresses := make(map[int]struct{})
	serials := make(map[int]struct{})
	a := struct{}{}
	for _, x := range products {
		addresses[x.Addr] = a
		serials[x.Serial] = a

	}
	serial, addr := 1, 1
	for ; addr < 256; addr++ {
		if _, f := addresses[addr]; !f {
			break
		}
	}
	for serial = 1; serial < 100500; serial++ {
		if _, f := serials[serial]; !f {
			break
		}
	}
	data.DB.MustExec(`
INSERT INTO product( party_id, serial, addr) 
VALUES (?,?,?)`, data.LastParty().PartyID, serial, addr)
	return getLastPartyProducts2(r)
}

func (_ *LastPartySvc) Products(_ struct{}, r *[]LastPartyProduct) error {
	return getLastPartyProducts2(r)
}

func getLastPartyProducts(r *[]LastPartyProduct) error {
	return data.DB.Select(r,
		`
SELECT * FROM product 
WHERE party_id = (SELECT party_id FROM last_party) 
ORDER BY product_id`)
}

func getLastPartyProducts2(r *[]LastPartyProduct) error {
	if err := getLastPartyProducts(r); err != nil {
		return err
	}
	c := cfg.Get()
	for i := range *r {
		(*r)[i].Place = i
		_, f := c.PlacesUncheck[i]
		(*r)[i].Checked = !f
	}
	return nil
}
