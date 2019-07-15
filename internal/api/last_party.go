package api

import (
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/mil82/internal/data"
	"github.com/fpawel/mil82/internal/party"
	"strconv"
)

type LastPartySvc struct{}

func (_ *LastPartySvc) Party(_ struct{}, r *data.Party) error {
	*r = data.LastParty()
	return nil
}

func (_ *LastPartySvc) SetPartySettings(x struct{ A data.PartySettings }, _ *struct{}) error {
	data.DB.MustExec(`
UPDATE party SET product_type = ?, c1 = ?, c2 = ?, c3 = ?, c4 = ?
WHERE party_id = (SELECT party_id FROM last_party)`,
		x.A.ProductType, x.A.C1, x.A.C2, x.A.C3, x.A.C4)
	return nil
}

func (_ *LastPartySvc) SetProductSerial(x struct {
	ProductID int64
	SerialStr string
}, _ *struct{}) error {

	serial, err := strconv.Atoi(x.SerialStr)
	if err != nil {
		return err
	}
	_, err = data.DB.Exec(`UPDATE product SET serial = ? WHERE product_id = ?`, serial, x.ProductID)
	if err != nil {
		return err
	}
	return nil
}

func (_ *LastPartySvc) SetProductAddr(x struct {
	ProductID int64
	AddrStr   string
}, _ *struct{}) error {
	addr, err := strconv.Atoi(x.AddrStr)
	if err != nil {
		return err
	}
	_, err = data.DB.Exec(`UPDATE product SET addr = ? WHERE product_id = ?`, addr, x.ProductID)
	if err != nil {
		return err
	}
	return nil
}

func (_ *LastPartySvc) DeleteProduct(productID [1]int64, r *[]party.Product) error {
	data.DB.MustExec(`DELETE FROM product WHERE product_id = ?`, productID[0])
	*r = party.Products()
	return nil
}

func (_ *LastPartySvc) AddNewProduct(_ struct{}, r *[]party.Product) error {
	products := party.Products()
	addresses := make(map[modbus.Addr]struct{})
	serials := make(map[int]struct{})
	a := struct{}{}
	for _, x := range products {
		addresses[x.Addr] = a
		serials[x.Serial] = a

	}
	serial, addr := 1, modbus.Addr(1)
	for ; addr <= modbus.Addr(255); addr++ {
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

	*r = party.Products()
	return nil
}

func (_ *LastPartySvc) Products(_ struct{}, r *[]party.Product) error {
	*r = party.Products()
	if *r == nil {
		*r = []party.Product{}
	}
	return nil
}
