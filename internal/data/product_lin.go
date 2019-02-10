package data

//go:generate reform

// ProductLin represents a row in product_lin table.
//reform:product_lin
type ProductLin struct {
	ProductLinID int64   `reform:"product_lin_id,pk"`
	ProductID    int64   `reform:"product_id"`
	LinPoint     int64   `reform:"lin_point"`
	Value        float64 `reform:"value"`
}
