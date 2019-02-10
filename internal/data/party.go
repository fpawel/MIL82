package data

import (
	"time"
)

//go:generate reform

// Party represents a row in party table.
//reform:party
type Party struct {
	PartyID          int64     `reform:"party_id,pk"`
	CreatedAt        time.Time `reform:"created_at"`
	ProductType      string    `reform:"product_type"`
	Pgs1             float64   `reform:"pgs1"`
	Pgs2             float64   `reform:"pgs2"`
	Pgs3             float64   `reform:"pgs3"`
	Pgs4             float64   `reform:"pgs4"`
	TemperatureNorm  float64   `reform:"temperature_norm"`
	TemperatureMinus float64   `reform:"temperature_minus"`
	TemperaturePlus  float64   `reform:"temperature_plus"`
}
