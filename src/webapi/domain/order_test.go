package domain_test

import (
	"testing"
	"time"

	"github.com/serdarkalayci/goboiler/webapi/domain"
)

func Test_AddProduct(t *testing.T) {
	customer := createCustomer("Customer1", "Customer Name1", 30)
	order := createOrder("Order1", customer)
	product1 := createProduct("Product1", "Product One", 7.75, 20)
	orderItem1 := createOrderItem(product1, 2)
	err := order.AddProduct(orderItem1)
	if err != nil || order.Total != 15.50 || order.Customer.Balance != 14.50 {
		t.Errorf("Error while adding first items. Total expected: 15.50 got: %f, Customer balance expected: 14.50 got: %f", order.Total, order.Customer.Balance)
	}
	product2 := createProduct("Product2", "Product Two", 1.25, 20)
	orderItem2 := createOrderItem(product2, 2)
	err = order.AddProduct(orderItem2)
	if err != nil || order.Total != 18 || order.Customer.Balance != 12 {
		t.Errorf("Error while adding new items. Total expected: 18.00 got: %f, Customer balance expected: 12.00 got: %f", order.Total, order.Customer.Balance)
	}
	err = order.AddProduct(orderItem1)
	if err == nil || order.Total != 18 || order.Customer.Balance != 12 {
		t.Errorf("Error while adding items that is out of Customer's balance. Total expected: 18.00 got: %f, Customer balance expected: 12.00 got: %f", order.Total, order.Customer.Balance)
	}
	err = order.AddProduct(orderItem2)
	if err != nil || order.Total != 20.50 || order.Customer.Balance != 9.50 {
		t.Errorf("Error while adding items of the same kind. Total expected: 20.50 got: %f, Customer balance expected: 9.50 got: %f", order.Total, order.Customer.Balance)
	}
}

func Test_RemoveProduct(t *testing.T) {
	customer := createCustomer("Customer1", "Customer Name1", 30)
	order := createOrder("Order1", customer)
	product1 := createProduct("Product1", "Product One", 7.75, 20)
	orderItem1 := createOrderItem(product1, 3)
	order.AddProduct(orderItem1)
	product2 := createProduct("Product2", "Product Two", 1.25, 20)
	orderItem2 := createOrderItem(product2, 2)
	removeItem := createOrderItem(product1, 1)
	err := order.RemoveProduct(removeItem)
	if err != nil || order.Total != 15.50 || order.Customer.Balance != 14.50 {
		t.Errorf("Error while removing first items. Total expected: 15.50 got: %f, Customer balance expected: 14.50 got: %f", order.Total, order.Customer.Balance)
	}
	err = order.RemoveProduct(orderItem2)
	if err == nil || order.Total != 15.50 || order.Customer.Balance != 14.50 {
		t.Errorf("Error while trying to remove non-order items. Total expected: 15.50 got: %f, Customer balance expected: 14.50 got: %f", order.Total, order.Customer.Balance)
	}
	removeItem = createOrderItem(product1, 3)
	err = order.RemoveProduct(removeItem)
	if err == nil || order.Total != 15.50 || order.Customer.Balance != 14.50 {
		t.Errorf("Error while trying to remove items more than those included in the Order. Total expected: 15.50 got: %f, Customer balance expected: 14.50 got: %f", order.Total, order.Customer.Balance)
	}
	removeItem = createOrderItem(product1, 2)
	err = order.RemoveProduct(removeItem)
	if err != nil || order.Total != 0 || order.Customer.Balance != 30 || len(order.Items) != 0 {
		t.Errorf("Error while trying to remove all items in the Order. Total expected: 0 got: %f, Customer balance expected: 30 got: %f", order.Total, order.Customer.Balance)
	}
}

func createOrder(id string, customer domain.Customer) domain.Order {
	return domain.Order{
		ID:       id,
		Date:     time.Now(),
		Customer: customer,
	}
}

func createOrderItem(product domain.Product, count int) domain.OrderItem {
	return domain.OrderItem{
		Item:      product,
		ItemCount: count,
	}
}
