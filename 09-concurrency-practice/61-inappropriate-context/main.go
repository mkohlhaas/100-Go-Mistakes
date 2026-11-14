package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handler1(w http.ResponseWriter, r *http.Request) {
	response, err := doSomeTask(r.Context(), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func() {
		err := publish(r.Context(), response)
		// Do something with err
		_ = err
	}()

	writeResponse(response)
}

func handler2(w http.ResponseWriter, r *http.Request) {
	response, err := doSomeTask(r.Context(), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func() {
		err := publish(context.Background(), response)
		// Do something with err
		_ = err
	}()

	writeResponse(response)
}

func handler3(w http.ResponseWriter, r *http.Request) {
	response, err := doSomeTask(r.Context(), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func() {
		err := publish(detach{ctx: r.Context()}, response)
		// Do something with err
		_ = err
	}()

	writeResponse(response)
}

type detach struct {
	ctx context.Context
}

func (d detach) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (d detach) Done() <-chan struct{} {
	return nil
}

func (d detach) Err() error {
	return nil
}

func (d detach) Value(key any) any {
	return d.ctx.Value(key)
}

func doSomeTask(ctx context.Context, r *http.Request) (string, error) {
	select {
	case <-time.After(50 * time.Millisecond):
	case <-ctx.Done():
		return "", ctx.Err()
	}
	return fmt.Sprintf("processed %s", r.URL.Path), nil
}

func publish(ctx context.Context, response string) error {
	select {
	case <-time.After(200 * time.Millisecond):
	case <-ctx.Done():
		return ctx.Err()
	}
	log.Printf("published: %s", response)
	return nil
}

func writeResponse(response string) {
	log.Printf("responded: %s", response)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/handler1", handler1)
	mux.HandleFunc("/handler2", handler2)
	mux.HandleFunc("/handler3", handler3)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Shut down the server on SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
