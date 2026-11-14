package main

import (
	"fmt"
	"sync"
)

func main() {
	customer := Customer{}
	_ = customer.UpdateAge2(-1)
	_ = customer.UpdateAge3(-1)
	_ = customer.UpdateAge1(-1)
}

type Customer struct {
	mutex sync.RWMutex
	id    string
	age   int
}

// UpdateAge1 - %v inside lock 👎
func (c *Customer) UpdateAge1(age int) error {
	fmt.Println("UpdateAge1: 👎")
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if age < 0 {
		return fmt.Errorf("age should be positive for customer %v", c)
	}

	c.age = age
	return nil
}

// UpdateAge2 - %v outside lock 👍
func (c *Customer) UpdateAge2(age int) error {
	fmt.Println("UpdateAge2: 👍")
	if age < 0 {
		return fmt.Errorf("age should be positive for customer %v", c)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.age = age
	return nil
}

// UpdateAge3 - %s inside lock 👍
func (c *Customer) UpdateAge3(age int) error {
	fmt.Println("UpdateAge3: 👍")
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if age < 0 {
		return fmt.Errorf("age should be positive for customer id %s", c.id)
	}

	c.age = age
	return nil
}

func (c *Customer) String() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return fmt.Sprintf("id %s, age %d", c.id, c.age)
}
