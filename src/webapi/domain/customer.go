package domain

import "errors"

// CustomerRepository represents an interface for the outer layers to implement the actual low level operations
type CustomerRepository interface {
	Store(customer Customer)
	Fetch(customerID string) Customer
}

// Customer defines the structure for an customer
type Customer struct {
	// the id of the customer
	//
	// required: true
	ID string

	// the name of the customer
	//
	// required: true
	Name string

	// current balance of the customer within our shop
	//
	// required: false
	Balance float64
}

// Rebalance calculates the current balance of the customer with the given (positive or negative) amount
func (customer *Customer) Rebalance(amount float64) error {
	if customer.Balance+amount < 0 {
		return errors.New("The customer balance cannot go below 0")
	}
	customer.Balance += amount
	return nil
}
