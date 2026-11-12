package main

import (
	"context"
	"errors"
	"fmt"
)

type loc struct{}

func (l loc) getCoordinates1(ctx context.Context, address string) (
	lat, lng float32, err error,
) {
	isValid := l.validateAddress(address)
	if !isValid {
		return 0, 0, errors.New("invalid address")
	}

	if ctx.Err() != nil {
		return 0, 0, err // ⚠️ err still nil
	}

	// Get and return coordinates
	return 0, 0, nil
}

func (l loc) getCoordinates2(ctx context.Context, address string) (
	lat, lng float32, err error,
) {
	isValid := l.validateAddress(address)
	if !isValid {
		return 0, 0, errors.New("invalid address")
	}

	if err := ctx.Err(); err != nil {
		return 0, 0, err // ⚠️ shadowing
	}

	// Get and return coordinates
	return 0, 0, nil
}

func (l loc) getCoordinates3(ctx context.Context, address string) (
	lat, lng float32, err error,
) {
	isValid := l.validateAddress(address)
	if !isValid {
		return 0, 0, errors.New("invalid address")
	}

	if err = ctx.Err(); err != nil { // 👍 assigning to the named err return value
		return
	}

	// Get and return coordinates
	return 0, 0, nil
}

func (l loc) validateAddress(address string) bool {
	_ = address
	return true
}

func main() {
	l := loc{}

	// A cancelled context: any geocode lookup against it should fail.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	address := "1 Main St"

	lat, lng, err := l.getCoordinates1(ctx, address)
	fmt.Printf("getCoordinates1: lat=%v lng=%v err=%v (want a context error, got %v)\n",
		lat, lng, err, err == nil) // BUG: ctx.Err() result is discarded, err stays nil
	// Output: getCoordinates1: lat=0 lng=0 err=<nil> (want a context error, got true)

	lat, lng, err = l.getCoordinates2(ctx, address)
	fmt.Printf("getCoordinates2: lat=%v lng=%v err=%v\n", lat, lng, err)
	// Output: getCoordinates2: lat=0 lng=0 err=context canceled

	lat, lng, err = l.getCoordinates3(ctx, address)
	fmt.Printf("getCoordinates3: lat=%v lng=%v err=%v\n", lat, lng, err)
	// Output: getCoordinates3: lat=0 lng=0 err=context canceled
}
