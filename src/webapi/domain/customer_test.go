package domain_test

import (
	"testing"

	"github.com/serdarkalayci/goboiler/webapi/domain"
)

func Test_RebalanceCustomer(t *testing.T) {
	customer := createCustomer("Customer_ID", "Test Customer", 23.85)
	customer.Rebalance(13.67)
	if customer.Balance != 37.52 {
		t.Errorf("Customer balance is not correct. Expected 37.52, got %f", customer.Balance)
	}
	customer.Rebalance(-13.67)
	if customer.Balance != 23.85 {
		t.Errorf("Customer balance is not correct. Expected 23.85, got %f", customer.Balance)
	}
	customer.Rebalance(-50)
	if customer.Balance < 0 {
		t.Errorf("Customer balance should not go below 0. Expected 23.85, got %f", customer.Balance)
	}
}

func createCustomer(id string, name string, balance float64) domain.Customer {
	return domain.Customer{
		ID:      id,
		Name:    name,
		Balance: balance,
	}
}
