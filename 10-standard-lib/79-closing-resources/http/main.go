package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

func (h handler) getBody1() (string, error) {
	resp, err := h.client.Get(h.url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (h handler) getBody2() (string, error) {
	resp, err := h.client.Get(h.url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response: %v\n", err)
		}
	}()

	return string(body), nil
}

func (h handler) getStatusCode1(body io.Reader) (int, error) {
	resp, err := h.client.Post(h.url, "application/json", body)
	if err != nil {
		return 0, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response: %v\n", err)
		}
	}()

	return resp.StatusCode, nil
}

func (h handler) getStatusCode2(body io.Reader) (int, error) {
	resp, err := h.client.Post(h.url, "application/json", body)
	if err != nil {
		return 0, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response: %v\n", err)
		}
	}()

	_, _ = io.Copy(io.Discard, resp.Body)

	return resp.StatusCode, nil
}

type handler struct {
	client http.Client
	url    string
}

func main() {
	// A tiny local server keeps the example self-contained and runnable
	// without any external network dependency.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello from the server\n"))
	}))
	defer srv.Close()

	h := handler{
		client: http.Client{},
		url:    srv.URL,
	}

	body, err := h.getBody1()
	if err != nil {
		log.Printf("getBody1: %v", err)
	} else {
		log.Printf("getBody1: %q", body) // 2026/11/15 16:19:36 getBody1: "hello from the server\n"
	}

	body, err = h.getBody2()
	if err != nil {
		log.Printf("getBody2: %v", err)
	} else {
		log.Printf("getBody2: %q", body) // 2026/11/15 16:19:36 getBody2: "hello from the server\n"
	}

	payload := strings.NewReader(`{"key":"value"}`)
	code, err := h.getStatusCode1(payload)
	if err != nil {
		log.Printf("getStatusCode1: %v", err)
	} else {
		log.Printf("getStatusCode1: %d", code) // 2026/11/15 16:19:36 getStatusCode1: 200
	}

	code, err = h.getStatusCode2(payload)
	if err != nil {
		log.Printf("getStatusCode2: %v", err)
	} else {
		log.Printf("getStatusCode2: %d", code) // 2026/11/15 16:19:36 getStatusCode2: 200
	}
}
