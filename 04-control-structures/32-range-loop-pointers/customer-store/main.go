package main

import "fmt"

type Customer struct {
	ID      string
	Balance float64
}

type Store struct {
	m map[string]*Customer
}

func main() {
	customers := []Customer{
		{ID: "1", Balance: 10},
		{ID: "2", Balance: -10},
		{ID: "3", Balance: 0},
	}

	fmt.Println("=== storeCustomers (buggy) ===")
	s1 := Store{m: make(map[string]*Customer)}
	s1.storeCustomers(customers)
	print(s1.m)

	fmt.Println("=== storeCustomers2 (copy before &) ===")
	s2 := Store{m: make(map[string]*Customer)}
	s2.storeCustomers2(customers)
	print(s2.m)

	fmt.Println("=== storeCustomers3 (index into slice) ===")
	s3 := Store{m: make(map[string]*Customer)}
	s3.storeCustomers3(customers)
	print(s3.m)
}

func (s *Store) storeCustomers(customers []Customer) {
	for _, customer := range customers {
		fmt.Printf("%p\n", &customer)
		s.m[customer.ID] = &customer
	}
}

func (s *Store) storeCustomers2(customers []Customer) {
	for _, customer := range customers {
		current := customer
		s.m[current.ID] = &current
	}
}

func (s *Store) storeCustomers3(customers []Customer) {
	for i := range customers {
		customer := &customers[i]
		s.m[customer.ID] = customer
	}
}

func print(m map[string]*Customer) {
	for k, v := range m {
		fmt.Printf("key=%s, value=%#v\n", k, v)
	}
}
