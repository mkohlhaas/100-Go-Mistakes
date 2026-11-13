package main

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

// task simulates CPU-bound work on a chunk of data.
// It is intentionally expensive enough that spreading chunks across
// multiple workers can reduce total runtime — provided there are
// enough chunks to amortize goroutine and channel overhead.
func task(b []byte) int {
	var sum int
	for _, v := range b {
		sum += int(v)
	}
	for i := range b {
		sum ^= int(b[i]) * 31
	}
	return sum
}

// dummyReader produces fixed-size chunks. In real code the reader could be
// a file, network connection, or any io.Reader. Note that io.Reader is not
// safe for concurrent use, so only one goroutine may call Read at a time.
type dummyReader struct {
	i int
}

func (c *dummyReader) Read(p []byte) (n int, err error) {
	const total = 100
	if c.i >= total {
		return 0, io.EOF
	}
	for i := range p {
		p[i] = byte(c.i)
	}
	c.i++
	return len(p), nil
}

// read1 processes chunks sequentially.
// Best when work per chunk is cheap or ordering matters.
func read1(r io.Reader) (int, error) {
	count := 0
	for {
		b := make([]byte, 1024)
		n, err := r.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return 0, err
		}
		// Process only the bytes actually read, not the whole buffer.
		count += task(b[:n])
	}
	return count, nil
}

// read2 processes chunks with a fixed pool of worker goroutines.
// Useful when task() is CPU-bound and chunks are independent:
// the producer keeps filling the channel while workers crunch earlier chunks.
// Not helpful when task() is trivial — the coordination overhead dominates.
func read2(r io.Reader) (int, error) {
	var count int64
	wg := sync.WaitGroup{}
	workers := 10

	ch := make(chan []byte, workers)
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for b := range ch {
				v := task(b)
				atomic.AddInt64(&count, int64(v))
			}
		}()
	}

	for {
		b := make([]byte, 1024)
		n, err := r.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return 0, err
		}
		ch <- b[:n] // send only the bytes actually read
	}

	close(ch)
	wg.Wait()
	return int(count), nil
}

func main() {
	res1, _ := read1(&dummyReader{})
	fmt.Println(res1) // 5068800

	res2, _ := read2(&dummyReader{})
	fmt.Println(res2) // 5068800
}
