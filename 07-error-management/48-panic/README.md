# Panic (Mistake #48)

This example shows a common Go mistake: using `panic` where an `error` should be returned instead.

`panic` is meant for unrecoverable, programmer errors — not for routine error handling. Recovering from a panic is possible, but it makes control flow harder to follow and should be avoided in normal application logic.

## The mistake

```go
func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover", r)
		}
	}()

	f()
}

func f() {
	fmt.Println("a")
	panic("foo")   // panics instead of returning an error
	fmt.Println("b")
}
```

When `f()` panics with `"foo"`, execution stops immediately, the deferred `recover()` runs, and `"b"` is never printed. The caller has to wrap every call in a `defer/recover` block to keep the program alive — which is awkward and easy to forget.

## The fix

Return an `error` instead of panicking:

```go
func main() {
	if err := f(); err != nil {
		fmt.Println("error:", err)
	}
}

func f() error {
	fmt.Println("a")
	return fmt.Errorf("foo")   // returns an error instead of panicking
	fmt.Println("b")            // unreachable, but clearly signals intent
}
```

Now `main()` handles the failure with a plain `if err != nil` check — idiomatic, explicit, and no `recover` needed.

## Validating input with panic

`checkWriteHeaderCode` in this example panics on invalid input:

```go
func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}
```

Even here, returning an `error` is usually the better choice:

```go
func checkWriteHeaderCode(code int) error {
	if code < 100 || code > 999 {
		return fmt.Errorf("invalid WriteHeader code %v", code)
	}
	return nil
}
```

## When panic is appropriate

Reserve `panic` for situations the program cannot safely continue from, such as:

- An internal inconsistency that indicates a bug (e.g. an unreachable branch was hit)
- A failure during initialization that makes the program unusable
- A nil dereference you truly cannot guard against

For everything else — invalid input, network failures, missing files — return an `error` and let the caller decide what to do.

## Run it

```sh
go run main.go
```

Expected output:

```
a
recover foo
```

## Key takeaway

Don't use `panic` for ordinary error handling. Return `error` values instead, and let `panic` stay reserved for truly unrecoverable situations.