package main

import (
	"context"
	"fmt"
	"math"
	"sync"

	"golang.org/x/sync/errgroup"
)

func handler1(ctx context.Context, circles []Circle) ([]Result, error) {
	results := make([]Result, len(circles))
	wg := sync.WaitGroup{}
	wg.Add(len(results))

	var mu sync.Mutex
	var firstErr error

	for i, circle := range circles {
		go func() {
			defer wg.Done()

			result, err := computeCircleArea(ctx, circle)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
				return
			}
			results[i] = result
		}()
	}

	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}
	return results, nil
}

func handler2(ctx context.Context, circles []Circle) ([]Result, error) {
	results := make([]Result, len(circles))
	g, ctx := errgroup.WithContext(ctx)

	for i, circle := range circles {
		g.Go(func() error {
			result, err := computeCircleArea(ctx, circle)
			if err != nil {
				return err
			}
			results[i] = result
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}

func computeCircleArea(_ context.Context, circle Circle) (Result, error) {
	if circle.Radius <= 0 {
		return Result{}, fmt.Errorf("circle %q has non-positive radius %g", circle.Name, circle.Radius)
	}
	return Result{Name: circle.Name, Area: math.Pi * circle.Radius * circle.Radius}, nil
}

type (
	Point struct {
		X, Y float64
	}
	Circle struct {
		Name   string
		Center Point
		Radius float64
	}
	Result struct {
		Name string
		Area float64
	}
)

func main() {
	ctx := context.Background()
	circles := []Circle{
		{Name: "small", Center: Point{X: 0, Y: 0}, Radius: 1},
		{Name: "medium", Center: Point{X: 10, Y: 5}, Radius: 2.5},
		// ⚠️ (un)comment to see the errgroup in action
		{Name: "broken", Center: Point{X: -3, Y: 4}, Radius: -1}, // invalid
		{Name: "large", Center: Point{X: 100, Y: 200}, Radius: 5},
	}

	// with errgroup
	results2, err := handler2(ctx, circles)
	if err != nil {
		fmt.Println("handler2 error:", err)
		return
	}
	fmt.Println("handler2 results:")
	for _, r := range results2 {
		fmt.Printf("  %s: area = %.2f\n", r.Name, r.Area)
	}

	// manual book-keeping
	results, err := handler1(ctx, circles)
	if err != nil {
		fmt.Println("handler1 error:", err)
		return
	}
	fmt.Println("handler1 results:")
	for _, r := range results {
		fmt.Printf("  %s: area = %.2f\n", r.Name, r.Area)
	}
}
