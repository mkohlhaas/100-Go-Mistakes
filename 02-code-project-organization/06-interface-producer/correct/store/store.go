// Package store is the producer. It exports a concrete Store type and its
// methods, NOT a big interface that consumers must conform to.
package store

type Customer struct {
	ID          string
	Name        string
	Balance     float64
	HasContract bool
}

// Store is the concrete customer storage.
type Store struct{}

// NewStore creates a new Store.
func NewStore() *Store {
	return &Store{}
}

// GetAllCustomers returns every customer.
// 👍 fulfills the customersGetter interface of the client (consumer)
func (s *Store) GetAllCustomers() ([]Customer, error) {
	// ...
	return nil, nil
}

// --------------- //
// More Methods... //
// --------------- //

// StoreCustomer persists a new customer.
func (s *Store) StoreCustomer(customer Customer) error {
	// ...
	return nil
}

// GetCustomer returns a single customer by id.
func (s *Store) GetCustomer(id string) (Customer, error) {
	// ...
	return Customer{}, nil
}

// UpdateCustomer modifies an existing customer.
func (s *Store) UpdateCustomer(customer Customer) error {
	// ...
	return nil
}

// GetCustomersWithoutContract returns customers without an active contract.
func (s *Store) GetCustomersWithoutContract() ([]Customer, error) {
	// ...
	return nil, nil
}

// GetCustomersWithNegativeBalance returns customers whose balance is below zero.
func (s *Store) GetCustomersWithNegativeBalance() ([]Customer, error) {
	// ...
	return nil, nil
}
