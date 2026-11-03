# 02 — Nested Code

This example demonstrates the **arrow anti-pattern** (arrow code) and the
idiomatic Go alternative: **guard clauses with early returns**.

## The problem: `join1`

`join1` wraps every validity check inside an `else` block:

```go
if s1 == "" {
    return "", errors.New("s1 is empty")
} else {
    if s2 == "" {
        return "", errors.New("s2 is empty")
    } else {
        // ... more nesting ...
    }
}
```

The successful path drifts further and further to the right, so the code
resembles the head of an arrow. This is hard to read:

- The happy path is buried at the deepest indentation level.
- Adding another check means wrapping everything in yet another `else`.
- You must hold the whole chain of `else` branches in your head to understand
  when the final `return` actually runs.

## The fix: `join2`

`join2` uses guard clauses — each invalid condition returns immediately:

```go
if s1 == "" {
    return "", errors.New("s1 is empty")
}
if s2 == "" {
    return "", errors.New("s2 is empty")
}
concat, err := concatenate(s1, s2)
if err != nil {
    return "", err
}
if len(concat) > max {
    return concat[:max], nil
}
return concat, nil
```

This is the preferred Go style because:

- **Flat structure** — no deep indentation, less cognitive load.
- **Linear happy path** — the success case reads top-to-bottom at the left
  margin, without jumping through `else` blocks.
- **Easy to extend** — add a new check as one more guard clause instead of
  nesting the entire function further.

Both functions behave identically; only the *shape* of the code differs.

## Run it

```bash
cd 02-nested-code
go run main.go
```

(Note: `concatenate` is a stub that returns `"", nil`, so neither function
produces output on its own. This snippet is about code structure, not runtime
behavior.)

## Source

Inspired by *100 Go Mistakes and How to Avoid Them* — chapter on nested code
and project organization.
