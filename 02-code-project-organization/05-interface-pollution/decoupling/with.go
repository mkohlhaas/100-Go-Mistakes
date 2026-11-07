package main

// flexible
type customerStorer interface {
	StoreCustomer(Customer) error
}

type CustomerService2 struct {
	storer customerStorer // flexible: ANYTHING that fulfills the interface
}

func (cs CustomerService2) CreateNewCustomer(id string) error {
	customer := Customer{id: id}
	return cs.storer.StoreCustomer(customer)
}
