package domain

import (
	"errors"
)

// Product defines the structure for a product
type Product struct {
	// the id of the product
	//
	// required: true
	ID string

	// the name of the product
	//
	// required: true
	Name string

	// the Price of the product
	//
	// required: false
	Price float64

	// The count of items in the stock
	//
	// required: true
	StockCount int
}

// ProductRepository represents an interface for the outer layers to implement the actual low level operations
type ProductRepository interface {
	Store(product Product)
	Fetch(id string) Product
}

// Unshelf removes product from the stock when it's added to an order
// Returns error when the requested count is above the stock count
func (product *Product) Unshelf(count int) error {
	if product.StockCount < count {
		return errors.New("Not enough of that product in the stock")
	}
	product.StockCount -= count
	return nil
}

// Shelf adds product to the stock when it's removed from an order
func (product *Product) Shelf(count int) {
	product.StockCount += count
}
