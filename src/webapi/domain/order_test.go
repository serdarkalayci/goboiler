package domain_test

import (
	"testing"
	"time"

	"github.com/serdarkalayci/goboiler/webapi/domain"
)

func Test_AddProduct(t *testing.T) {
	customer := createCustomer("Customer1", "Customer Name1", 20)
	order := createOrder("Order1", customer)
	product := createProduct("Product1", "Product One", 7.75)
	orderItem := createOrderItem(product, 2)
	err := order.AddProduct(orderItem)
	if err != nil || order.Total != 15.50 {
		t.Errorf("Error while adding first items. Total expected: 15.50 got: %f", order.Total)
	}
}

func createOrder(id string, customer domain.Customer) domain.Order {
	return domain.Order{
		ID:       id,
		Date:     time.Now(),
		Customer: customer,
	}
}

func createProduct(id string, name string, price float64) domain.Product {
	return domain.Product{
		ID:    id,
		Name:  name,
		Price: price,
	}
}

func createOrderItem(product domain.Product, count int) domain.OrderItem {
	return domain.OrderItem{
		Item:      product,
		ItemCount: count,
	}
}
