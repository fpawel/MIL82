package party

import (
	"github.com/fpawel/mil82/internal/cfg"
	"github.com/fpawel/mil82/internal/data"
)

type Product struct {
	data.Product
	Place   int
	Checked bool
}

func CheckedProducts() (products []Product) {
	xs := Products()
	for _, x := range xs {
		if x.Checked {
			products = append(products, x)
		}
	}
	return
}

func Products() (products []Product) {
	if err := data.DB.Select(&products, `
SELECT * FROM product 
WHERE party_id = (SELECT party_id FROM last_party) 
ORDER BY product_id`); err != nil {
		panic(err)
	}
	c := cfg.Get()
	for i := range products {
		products[i].Place = i
		products[i].Checked = c.PlaceChecked(i)
	}
	return
}
