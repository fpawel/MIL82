package data

//go:generate reform

// ProductCoefficient represents a row in product_coefficient table.
//reform:product_coefficient
type ProductCoefficient struct {
	ProductCoefficientID int64   `reform:"product_coefficient_id,pk"`
	ProductID            int64   `reform:"product_id"`
	Coefficient          int64   `reform:"coefficient"`
	Value                float64 `reform:"value"`
}
