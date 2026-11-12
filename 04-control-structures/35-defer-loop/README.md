# Mistake 35: `defer` in a loop

## The problem

`defer` schedules a call to run when the **enclosing function returns**, not at the end of the current loop iteration.

In `readFiles1`, every `defer file.Close()` inside the `for` loop piles up. The files stay open until `readFiles1` returns. If you loop over many paths, you leak file descriptors and can hit the OS open-file limit.

```go
func readFiles1(ch <-chan string) error {
	for path := range ch {
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close() // BUG: runs when readFiles1 returns, not each iteration

		// Do something with file
	}
	return nil
}
```

## The fix

Open and close each file within the same iteration, either by delegating to a helper function (`readFiles2`) or by wrapping the body in an immediately-invoked closure (`readFiles3`). In both cases `defer` runs when the helper/closure returns, i.e. at the end of each iteration.

```go
func readFiles2(ch <-chan string) error {
	for path := range ch {
		if err := readFile(path); err != nil {
			return err
		}
	}
	return nil
}

func readFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close() // OK: runs when readFile returns

	// Do something with file
	return nil
}
```

## Rule of thumb

Never `defer` inside a loop unless you truly mean to run it once when the enclosing function returns. Otherwise, close resources explicitly each iteration, or move them into a helper function.

## Run it

```bash
go run main.go
```
