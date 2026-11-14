package main

import "fmt"

func merge1(ch1, ch2 <-chan int) <-chan int {
	ch := make(chan int, 1)

	go func() {
		for v := range ch1 {
			ch <- v
		}
		for v := range ch2 {
			ch <- v
		}
		close(ch)
	}()

	return ch
}

func merge2(ch1, ch2 <-chan int) <-chan int {
	ch := make(chan int, 1)

	go func() {
		for {
			select {
			case v := <-ch1:
				ch <- v
			case v := <-ch2:
				ch <- v
			}
		}
	}()

	return ch
}

func merge3(ch1, ch2 <-chan int) <-chan int {
	ch := make(chan int, 1)
	ch1Closed := false
	ch2Closed := false

	go func() {
		for {
			select {
			case v, open := <-ch1:
				if !open {
					ch1Closed = true
					break
				}
				ch <- v
			case v, open := <-ch2:
				if !open {
					ch2Closed = true
					break
				}
				ch <- v
			}

			if ch1Closed && ch2Closed {
				close(ch)
				return
			}
		}
	}()

	return ch
}

func merge4(ch1, ch2 <-chan int) <-chan int {
	ch := make(chan int, 1)

	go func() {
		for ch1 != nil || ch2 != nil {
			select {
			case v, open := <-ch1:
				if !open {
					ch1 = nil
					break
				}
				ch <- v
			case v, open := <-ch2:
				if !open {
					ch2 = nil
					break
				}
				ch <- v
			}
		}
		close(ch)
	}()

	return ch
}

func demo(name string, fn func(<-chan int, <-chan int) <-chan int) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		ch1 <- 1
		ch1 <- 2
		close(ch1)
	}()
	go func() {
		ch2 <- 3
		ch2 <- 4
		close(ch2)
	}()

	fmt.Printf("%s: ", name)
	for v := range fn(ch1, ch2) {
		fmt.Printf("%d ", v)
	}
	fmt.Println()
}

func main() {
	// merge1, merge3, and merge4 all terminate once both inputs are closed.
	demo("merge1", merge1)
	demo("merge3", merge3)
	demo("merge4", merge4)

	// merge2 never terminates: its select has no way to detect that both
	// inputs are closed, so close(ch) is never reached. It is shown here
	// for reference only — running it would hang forever.
	fmt.Println("merge2: (never terminates — see the README for why)")
	_ = merge2
}
