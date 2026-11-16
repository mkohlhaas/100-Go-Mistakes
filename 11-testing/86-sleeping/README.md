# 86 - Sleeping in tests

Part of the **100 Go Mistakes** series (chapter 11: testing, mistake 86).

## The mistake

`Handler.getBestFoo` publishes results from a background goroutine, then returns immediately:

```go
func (h Handler) getBestFoo(someInputs int) Foo {
    foos := getFoos(someInputs)
    best := foos[0]

    go func() {
        if len(foos) > h.n {
            foos = foos[:h.n]
        }
        h.publisher.Publish(foos)
    }()

    return best
}
```

A test that needs to assert on the published values must wait for that goroutine to finish.

### ❌ The sleeping approach (flaky)

`TestGetBestFoo` uses `time.Sleep`:

```go
time.Sleep(10 * time.Millisecond)
published := mock.Get()
```

Sleep-based waits are:

- **Flaky** — CI may be slower than your dev machine, causing intermittent failures.
- **Slow** — you pay the full sleep duration even when the work finishes instantly.
- **Fragile** — the magic number is a guessing game.

### ✅ The channel-based approach (correct)

`TestGetBestFoo2` uses a channel to signal completion:

```go
mock := publisherMock2{ch: make(chan []Foo)}
...
if v := len(<-mock.ch); v != 2 {
    t.Fatalf("expected 2, got %d", v)
}
```

The test blocks until the goroutine actually publishes, making it deterministic and fast.

## Run the tests

```bash
go test -v -race ./...
```

Run with `-race` to confirm the channel version is free of data races. The sleep version may pass locally but is timing-dependent — increase CI load or shorten the sleep to see it fail intermittently.

## Structure

- `main.go` — `Handler`, its `getBestFoo` method, and the `publisher` interface
- `main_test.go` — two tests comparing the sleeping approach vs the channel approach

## Takeaway

Never sleep in tests to wait for asynchronous work. Use channels, `sync.WaitGroup`, `context.Done()`, or similar synchronization primitives so your tests wait on the actual completion signal rather than a guessed duration.
