package main

// -------- //
// Customer //
// -------- //

type Customer struct {
	id string
}

// ----- //
// Store //
// ----- //

type Store struct{}

func (s Store) StoreCustomer(customer Customer) error {
	return nil
}

// --------------- //
// CustomerService //
// --------------- //

type CustomerService struct {
	store Store // ⚠️ hard-wired
}

func (cs CustomerService) CreateNewCustomer(id string) error {
	customer := Customer{id: id}
	return cs.store.StoreCustomer(customer)
}
