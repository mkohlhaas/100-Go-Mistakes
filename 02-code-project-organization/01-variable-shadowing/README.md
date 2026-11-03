# Variable Shadowing

Source example for the **variable shadowing** mistake from *100 Go Mistakes and How to Avoid Them*.

## What it demonstrates

In Go, the short variable declaration `:=` only reuses an existing variable if that variable was declared **in the same block**. If a variable with the same name already exists in an **enclosing** block, `:=` creates a brand-new variable in the inner block that **shadows** the outer one. When the inner block ends, the shadowing variable is discarded — the outer variable keeps its zero value.

This is a silent bug: the code compiles, runs, and looks correct, but the value you assigned inside an `if`/`else` branch never reaches the outer variable.

## The bug — `listing1`

```go
var client *http.Client          // function-scope (outer) variable
if tracing {
    client, err := createClientWithTracing()   // NEW client + NEW err here
    ...
} else {
    client, err := createDefaultClient()       // NEW client + NEW err here
    ...
}
_ = client                       // still nil! the outer variable was never assigned
```

Because `client` was declared outside the `if/else`, the `:=` inside each branch creates a **new** `client` scoped to that branch. The outer `client` stays `nil`.

## The fixes — `listing2`, `listing3`, `listing4`

Each listing fixes the same problem a different way:

| Listing | Fix |
|---------|-----|
| `listing2` | Declare `client` outside, then assign through a throwaway `c` inside each branch: `client = c`. |
| `listing3` | Declare both `client` and `err` outside, then use `=` (simple assignment, not `:=`) inside the branches. |
| `listing4` | Cleanest: declare both up front, assign with `=`, check `err` once after the `if/else`. |

## Key takeaway

Use `:=` to **declare** new variables. To **assign** to a variable from an enclosing block inside an inner block, declare the variable in the outer scope first, then use `=` (simple assignment). Never assume `:=` will "reach out" to an outer scope.

## Run it

```bash
go vet ./...
go run .
```

`go vet` passes for all four listings — the shadowing bug compiles cleanly, which is exactly why it's dangerous.
