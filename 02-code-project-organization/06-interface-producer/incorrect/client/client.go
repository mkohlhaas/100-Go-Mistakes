// Package client
package client

import "interface-producer-incorrect/store"

type customersGetter interface {
	GetAllCustomers() ([]store.Customer, error)
}
