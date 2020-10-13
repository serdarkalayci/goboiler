package usecases

import (
	"errors"

	"github.com/serdarkalayci/goboiler/webapi/domain"
)

// OrderOperator is the struct that hold both OrderRepository and CustomerRepository
type OrderOperator struct {
	orderRepository    domain.OrderRepository
	customerRepository domain.CustomerRepository
	productRepository  domain.ProductRepository
}

// AddProduct adds a product to the order
// Returns error if Order's CustomerID does not match customerID
// Returns error if the Customer does not have enough credit
// Returns error if the productCount is above Product's StockCount
func (oo *OrderOperator) AddProduct(orderID, customerID, productID string, productCount int) error {
	order := oo.orderRepository.Fetch(orderID)
	if order.Customer.ID != customerID {
		return errors.New("The order does not belong to this customer, cannot add products")
	}
	product := oo.productRepository.Fetch(productID)
	orderItem := domain.OrderItem{
		Item:      product,
		ItemCount: productCount,
	}
	err := order.AddProduct(orderItem)
	if err != nil {
		return err
	}
	return nil
}
