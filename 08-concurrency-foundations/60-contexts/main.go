package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"contexts/flight"
)

type publisher interface {
	Publish(ctx context.Context, position flight.Position) error
}

type publishHandler struct {
	pub publisher
}

func (h publishHandler) publishPosition(position flight.Position) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	return h.pub.Publish(ctx, position)
}

type key string

const isValidHostKey key = "isValidHost"

func checkValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validHost := r.Host == "acme"
		ctx := context.WithValue(r.Context(), isValidHostKey, validHost)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handler(ctx context.Context, ch chan Message) error {
	for {
		select {
		case msg := <-ch:
			// Do something with msg
			_ = msg
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

type Message struct{}

type mockPublisher struct{}

func (m mockPublisher) Publish(ctx context.Context, position flight.Position) error {
	_ = position
	// Simulate work that respects the context deadline
	select {
	case <-time.After(500 * time.Millisecond):
		return nil // work finished before the deadline
	case <-ctx.Done():
		return ctx.Err() // cancelled or timed out before work finished
	}
}

type slowPublisher struct {
	workTime time.Duration
}

func (s slowPublisher) Publish(ctx context.Context, position flight.Position) error {
	_ = position
	// Simulate slow work that exceeds the context deadline
	select {
	case <-time.After(s.workTime):
		return nil // work finished before the deadline
	case <-ctx.Done():
		return ctx.Err() // cancelled or timed out before work finished
	}
}

func main() {
	// 1. Timeout context -> publishPosition (work finishes before deadline)
	h := publishHandler{pub: mockPublisher{}}
	pos := flight.Position{Lat: 48.8566, Lng: 2.3522}
	if err := h.publishPosition(pos); err != nil {
		fmt.Println("publishPosition failed:", err)
		return
	}
	fmt.Println("publishPosition succeeded")

	// 1b. Timeout context -> publishPosition (work exceeds deadline: context.DeadlineExceeded)
	slowH := publishHandler{pub: slowPublisher{workTime: 5 * time.Second}}
	if err := slowH.publishPosition(pos); err != nil {
		fmt.Println("slowPublisher timed out:", err)
	}

	// 2. Request-scoped value -> checkValid
	_ = checkValid(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	// Demonstrate ctx.Done() being checked: cancel before work finishes
	cancelCtx, cancelCtxFn := context.WithCancel(context.Background())
	cancelCtxFn()
	err := mockPublisher{}.Publish(cancelCtx, pos)
	if err != nil {
		fmt.Println("manual cancel returned:", err)
	}

	// 3. Cancellation propagation -> handler
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan Message)
	done := make(chan error, 1)
	go func() {
		done <- handler(ctx, ch)
	}()
	cancel()
	if err := <-done; err != nil {
		fmt.Println("handler returned:", err)
	}
}
