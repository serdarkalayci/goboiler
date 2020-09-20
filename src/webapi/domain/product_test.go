package domain_test

import (
	"testing"

	"github.com/serdarkalayci/goboiler/webapi/domain"
)

func Test_Unshelf(t *testing.T) {
	product := createProduct("Product1", "Product One", 7.75, 20)
	err := product.Unshelf(25)
	if err == nil || product.StockCount != 20 {
		t.Errorf("Error unshelving product when the claim is more than stock. Expected %d, got %d", 20, product.StockCount)
	}
	err = product.Unshelf(5)
	if err != nil || product.StockCount != 15 {
		t.Errorf("Error unshelving product when the claim is less than stock. Expected %d, got %d", 15, product.StockCount)
	}
}

func Test_Shelf(t *testing.T) {
	product := createProduct("Product1", "Product One", 7.75, 20)
	product.Shelf(5)
	if product.StockCount != 25 {
		t.Errorf("Error shelving product. Expected %d, got %d", 25, product.StockCount)
	}
}

func createProduct(id string, name string, price float64, stockCount int) domain.Product {
	return domain.Product{
		ID:         id,
		Name:       name,
		Price:      price,
		StockCount: stockCount,
	}
}
