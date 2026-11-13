# 53 - Not handling errors

Example from Chapter 7 (Error Management) of *100 Go Mistakes and How to Avoid Them*.

## The mistake

Ignoring the error returned by a function. Silent failures make bugs hard to diagnose and can corrupt state downstream.

```go
func listing1() {
	// ...

	notify() // ERROR return value is silently ignored
}
```

If `notify()` fails (e.g. network is down, SMTP is misconfigured), nothing tells the caller. The program continues as if nothing happened.

## The acceptable exception

There are legitimate cases where discarding an error is intentional — for example, best-effort notifications where missing one is acceptable. In those cases, be explicit about it and document why:

```go
func listing2() {
	// ...

	// Notifications are sent in best effort.
	// Hence, it's accepted to miss some of them in case of errors.
	_ = notify()
}
```

## Best practices

- **Always check errors** returned by functions, unless you have a documented reason not to.
- **Never use `_ = err` silently** — add a comment explaining why the error is being ignored.
- **Prefer wrapping or returning errors** so callers can decide how to handle them.
- **Log errors** at an appropriate level when you cannot recover from them.

## Running

```bash
go run main.go
```

## Related

- Chapter 7: Error Management
- #52 Ignoring errors from `defer` cleanup
- #54 Swallowing errors with blank identifiers