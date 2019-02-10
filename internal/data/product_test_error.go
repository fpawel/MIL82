package data

//go:generate reform

// ProductTestError represents a row in product_test_error table.
//reform:product_test_error
type ProductTestError struct {
	ProductTestErrorID int64   `reform:"product_test_error_id,pk"`
	ProductID          int64   `reform:"product_id"`
	TestErrorPoint     string  `reform:"test_error_point"`
	GasPoint           int64   `reform:"gas_point"`
	Value              float64 `reform:"value"`
}
