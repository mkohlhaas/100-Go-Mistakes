package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// A JSON message whose "id" field looks like an integer.
	const msg = `{"id": 42, "name": "widget", "price": 9.99}`

	// --- The mistake: ⚠️ map[string]any ⚠️ --------------------------------------
	var m map[string]any
	if err := json.Unmarshal([]byte(msg), &m); err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}

	fmt.Printf("m = %#v\n", m)                                      // m = map[string]interface {}{"id":42, "name":"widget", "price":9.99}
	fmt.Printf("m[\"id\"] type: %T, value: %v\n", m["id"], m["id"]) // m["id"] type: float64, value: 42

	// This type assertion FAILS at runtime: JSON numbers decode to float64,
	// never to int, uint, int64, etc.
	id, ok := m["id"].(int)
	if !ok {
		fmt.Println("BUG: m[\"id\"].(int) failed — JSON numbers are float64, not int") // BUG: m["id"].(int) failed — JSON numbers are float64, not int
	}
	_ = id

	// Even a direct comparison to an int literal misbehaves:
	fmt.Printf("m[\"id\"] == 42 is %v (surprising, but correct: float64 vs int)\n", m["id"] == 42) // m["id"] == 42 is false (surprising, but correct: float64 vs int)

	// --- The fix: 👍 a concrete struct 👍 ---------------------------------------
	type Message struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	var message Message
	if err := json.Unmarshal([]byte(msg), &message); err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}

	fmt.Printf("message.ID type: %T, value: %v\n", message.ID, message.ID) // message.ID type: int, value: 42
	fmt.Printf("message.ID == 42 is %v\n", message.ID == 42)               // message.ID == 42 is true
	fmt.Printf("message = %+v\n", message)                                 // message = {ID:42 Name:widget Price:9.99}
}
