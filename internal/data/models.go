package data

import (
	"time"
)

type Party struct {
	PartyID     int64     `db:"party_id"`
	CreatedAt   time.Time `db:"created_at"`
	ProductType string    `db:"product_type"`
	C1          float32   `db:"c1"`
	C2          float32   `db:"c2"`
	C3          float32   `db:"c3"`
	C4          float32   `db:"c4"`
}

type Product struct {
	ProductID int64 `db:"product_id"`
	PartyID   int64 `db:"party_id"`
	Serial    int   `db:"serial"`
	Addr      int   `db:"addr"`
}

type ProductValue struct {
	ProductID int64   `db:"product_id"`
	Work      string  `db:"work"`
	Var       int     `db:"var"`
	Gas       string  `db:"gas"`
	Temp      string  `db:"temp"`
	Value     float32 `db:"value"`
}

type ProductCoefficient struct {
	ProductID   int64   `db:"product_id"`
	Coefficient int32   `db:"coefficient"`
	Value       float32 `db:"value"`
}
