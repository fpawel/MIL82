package data

//go:generate reform

// ProductTemp represents a row in product_temp table.
//reform:product_temp
type ProductTemp struct {
	ProductCoefficientID int64   `reform:"product_coefficient_id,pk"`
	ProductID            int64   `reform:"product_id"`
	GasPoint             int64   `reform:"gas_point"`
	TemperaturePoint     int64   `reform:"temperature_point"`
	Temperature          float64 `reform:"temperature"`
	Value                float64 `reform:"value"`
}
