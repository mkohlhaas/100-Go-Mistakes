// Package client is the consumer. It defines its own small, focused interface
// close to where it is used, rather than depending on a bloated producer
// interface.
package client

import store "interface-producer-correct/store"

// customersGetter is the minimal interface this client needs. It lives in the
// consumer package, not the producer package.
type customersGetter interface {
	GetAllCustomers() ([]store.Customer, error)
}

// Client wraps whatever satisfies customersGetter. The concrete store package
// satisfies it automatically because it exports GetAllCustomers returning
// []store.Customer.
type Client struct {
	getter customersGetter
}

// NewClient builds a Client from any value that provides GetAllCustomers.
func NewClient(getter customersGetter) *Client {
	return &Client{getter: getter}
}

// Report prints all customers.
func (c *Client) Report() error {
	customers, err := c.getter.GetAllCustomers()
	if err != nil {
		return err
	}
	for _, cust := range customers {
		_ = cust
	}
	return nil
}
