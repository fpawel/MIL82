package data

import (
	"time"
)

type Party struct {
	PartyID          int64     `db:"party_id"`
	CreatedAt        time.Time `db:"created_at"`
	ProductType      string    `db:"product_type"`
	Pgs1             float32   `db:"pgs1"`
	Pgs2             float32   `db:"pgs2"`
	Pgs3             float32   `db:"pgs3"`
	PgsLin12         float32   `db:"pgs_lin_12"`
	PgsLin22         float32   `db:"pgs_lin_22"`
	TemperatureNorm  float32   `db:"temperature_norm"`
	TemperatureMinus float32   `db:"temperature_minus"`
	TemperaturePlus  float32   `db:"temperature_plus"`
}

type Product struct {
	ProductID    int64 `db:"product_id,pk"`
	PartyID      int64 `db:"party_id"`
	SerialNumber int16 `db:"serial_number"`
	Place        int16 `db:"place"`
	Addr         int16 `db:"addr"`
	Production   bool  `db:"production"`
}

type ProductCheckup struct {
	ProductID int64   `db:"product_id"`
	Test      string  `db:"test"` // FIXME unhandled database type "USER-DEFINED"
	Gas       string  `db:"gas"`  // FIXME unhandled database type "USER-DEFINED"
	Value     float32 `db:"value"`
}

type ProductCoefficient struct {
	ProductID   int64   `db:"product_id"`
	Coefficient int32   `db:"coefficient"`
	Value       float32 `db:"value"`
}

type ProductLin struct {
	ProductID int64   `db:"product_id"`
	Gas       string  `db:"gas"` // FIXME unhandled database type "USER-DEFINED"
	Value     float32 `db:"value"`
}

type ProductTemperatureComp struct {
	ProductID        int64   `reform:"product_id"`
	Gas              string  `reform:"gas"`         // FIXME unhandled database type "USER-DEFINED"
	Temperature      string  `reform:"temperature"` // FIXME unhandled database type "USER-DEFINED"
	TemperatureValue float32 `reform:"temperature_value"`
	Value            float32 `reform:"value"`
}
