package domain

import (
	"errors"
	"time"
)

// OrderRepository represents an interface for the outer layers to implement the actual low level operations
type OrderRepository interface {
	Store(order Order)
	Fetch(orderID string) Order
}

// Order defines the structure for an order
type Order struct {
	// the id of the order
	//
	// required: true
	ID string

	// the placement date of the order
	//
	// required: true
	Date time.Time
	// the list of the products within the order
	//
	// required: true
	Items []OrderItem
	// the total value of added items
	//
	// required: true
	Total float64
	// the customer refenrence of the order
	//
	// required: true
	Customer Customer
}

// OrderItem represents the products and their counts to be added to the order
type OrderItem struct {
	// describes how many of this produıct will be added to order
	//
	// required: true
	ItemCount int
	// describes which product will be added to the order
	//
	// required: true
	Item Product
}

// AddProduct adds new Product and increase the count if the order already has that spesific product
// The Product is passed as OrderItem which includes the Product and the count to be added
// Returns an error if the Customer's balance is not enough for the requested amount
func (order *Order) AddProduct(orderItem OrderItem) error {
	newAmount := float64(orderItem.ItemCount) * orderItem.Item.Price
	err := order.Customer.Rebalance(newAmount * -1)
	if err != nil {
		return errors.New("Customer balance is not enough for adding these items")
	}
	err = orderItem.Item.Unshelf(orderItem.ItemCount)
	found := false
	for i := 0; i < len(order.Items) && found == false; i++ {
		if order.Items[i].Item.ID == orderItem.Item.ID {
			order.Items[i].ItemCount += orderItem.ItemCount
			found = true
		}
	}
	if !found {
		order.Items = append(order.Items, orderItem)
	}
	order.Total += newAmount
	return nil
}

// RemoveProduct decrease the count od a Product in the order
// The Product is passed as OrderItem which includes the Product and the count to be added
// Returns an error if the Order does not contain that Product or the count is already below the requested amount
func (order *Order) RemoveProduct(orderItem OrderItem) error {
	found := false
	i := 0
	for i = 0; i < len(order.Items) && found == false; i++ {
		if order.Items[i].Item.ID == orderItem.Item.ID {
			found = true
		}
	}
	if !found {
		return errors.New("The order does not contain the Product requested")
	}
	if order.Items[i-1].ItemCount < orderItem.ItemCount { // There's not enough items to be removed
		return errors.New("The order does not contain enough of the Product requested")
	} else if order.Items[i-1].ItemCount == orderItem.ItemCount { // There's exactly the number of items to be removed
		order.Items = append(order.Items[:i-1], order.Items[i:]...) // Remove that item from the array completely
	} else { // There're more items than to be removed
		order.Items[i-1].ItemCount -= orderItem.ItemCount
	}
	orderItem.Item.Shelf(orderItem.ItemCount)
	newAmount := float64(orderItem.ItemCount) * orderItem.Item.Price
	order.Customer.Rebalance(newAmount)
	order.Total -= newAmount
	return nil
}
