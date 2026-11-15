# os-file — Closing resources (Go Mistake #79)

This example demonstrates a common mistake when closing file resources in Go:
**mishandling errors from `f.Close()` when using `defer`**.

## The problem

When you `defer f.Close()`, the close runs when the surrounding function
returns. Two things commonly go wrong:

### 1. A close error masks a real error (`writeToFile1`)

```go
func writeToFile1(filename string, content []byte) (err error) {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
    if err != nil {
        return err
    }

    defer func() {
        closeErr := f.Close()
        if err == nil {
            err = closeErr          // close error can overwrite the Write error
        }
    }()

    _, err = f.Write(content)
    return
}
```

Because the function uses a **named** return value `err`, the deferred
closure can mutate it. If `f.Write` fails, `err` is already non-nil so the
close error is ignored (the `Write` error is returned instead). If `Write`
succeeds but `Close` fails, the close error is returned — which is reasonable
behavior on its own, but the pattern is fragile: it relies on `defer`
mutating a named return, and when both `Write` and `Close` fail you only get
one of the errors.

### 2. Close errors are silently discarded (`writeToFile2`)

```go
func writeToFile2(filename string, content []byte) (err error) {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
    if err != nil {
        return err
    }

    defer func() {
        _ = f.Close()   // close error swallowed
    }()

    _, err = f.Write(content)
    if err != nil {
        return err
    }

    return f.Sync()
}
```

Here `_ = f.Close()` throws the close error away. If the flush/close fails,
you have no way of knowing your data may not have been persisted to disk.

## The fix

**Important:** `defer` runs *after* the `return` statement is evaluated, not before. So
`return closeErr` copies `closeErr` (still `nil`) into the result slot *before* the deferred
`Close` runs — the close error is lost. You need a **named return value** so `defer` can
modify the result that the caller receives.

Check both the operation error **and** the close error independently, without clobbering one
with the other:

```go
func writeToFile(filename string, content []byte) (err error) {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
    if err != nil {
        return
    }

    defer func() {
        if e := f.Close(); e != nil {
            if err == nil {
                err = e
            }
        }
    }()

    if _, err = f.Write(content); err != nil {
        return
    }

    if err = f.Sync(); err != nil {
        return
    }

    return
}
```

### With `errors.Join` (Go 1.20+)

If you want to surface both the operation error and the close error instead of picking just
one, use `errors.Join`:

```go
func writeToFileJoined(filename string, content []byte) (err error) {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
    if err != nil {
        return
    }

    defer func() {
        if e := f.Close(); e != nil {
            if err == nil {
                err = e
            } else {
                err = errors.Join(err, e)
            }
        }
    }()

    if _, err = f.Write(content); err != nil {
        return
    }

    if err = f.Sync(); err != nil {
        return
    }

    return
}
```

`errors.Join` combines multiple errors into a single multi-error, so neither the
write/sync failure nor the close failure is lost. Callers can inspect the
combined error with `errors.Is` or `errors.As` as usual. Don't forget to
import `"errors"`.

Key points:

- Don't use a local variable (`var closeErr error`) with `return closeErr` — `defer`
  runs after the return value is copied, so the close error is lost. Use a **named
  return value** instead.
- Don't use a named return value just so `defer` can silently overwrite a real error —
  let the operation error take precedence when it's non-nil.
- Never just `_ = f.Close()` — a failed close can mean data wasn't flushed
  to disk.

## Running

`main()` calls all five functions so you can run the example directly:

```bash
go run .
```

Note that the original functions open with `O_APPEND|O_WRONLY` but not
`O_CREATE`, so the target files must already exist — otherwise you'll see
`open ...: no such file or directory`. Create the files first if you want a
clean run:

```bash
touch /tmp/opencode/listing1.txt /tmp/opencode/file1.txt \
      /tmp/opencode/file2.txt /tmp/opencode/file_fixed.txt \
      /tmp/opencode/file_joined.txt
go run .
```

## Reference

*Mistake #79: Closing resources* from *100 Go Mistakes and How to Avoid Them*.