// Package main demonstrates using NewClient with the store.
package main

import (
	"fmt"

	"interface-producer-correct/client"
	"interface-producer-correct/store"
)

func main() {
	// Producer: create a concrete Store.
	s := store.NewStore()

	// Consumer: wire the store into Client via NewClient.
	// *store.Store satisfies client.customersGetter because it provides
	// GetAllCustomers() ([]store.Customer, error).
	c := client.NewClient(s)

	if err := c.Report(); err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("report complete")
}
