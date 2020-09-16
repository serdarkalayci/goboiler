package domain

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
}

// ProductRepository represents an interface for the outer layers to implement the actual low level operations
type ProductRepository interface {
	Store(product Product)
	Fetch(id string) Product
}
