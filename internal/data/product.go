package data

//go:generate reform

// Product represents a row in product table.
//reform:product
type Product struct {
	ProductID  int64 `reform:"product_id,pk"`
	PartyID    int64 `reform:"party_id"`
	Serial     int64 `reform:"serial"`
	Place      int64 `reform:"place"`
	Addr       int64 `reform:"addr"`
	Production bool  `reform:"production"`
}
