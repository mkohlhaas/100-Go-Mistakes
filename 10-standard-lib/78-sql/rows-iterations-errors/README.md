# rows-iterations-errors

Part of the **100 Go Mistakes** series — Chapter 10 (Standard Library), Mistake 78: SQL.

## The mistake

When iterating over `sql.Rows` with `rows.Next()`, a `rows.Scan()` call can fail partway through the result set (for example, due to a type conversion error or a network hiccup). When that happens, `rows.Next()` returns `false` and the `for` loop exits silently — **you never learn that an error occurred.**

The code below reads a department and age from an `EMP` table:

```go
func get1(ctx context.Context, db *sql.DB, id string) (string, int, error) {
	rows, err := db.QueryContext(ctx, "SELECT DEP, AGE FROM EMP WHERE ID = ?", id)
	if err != nil {
		return "", 0, err
	}
	defer rows.Close()

	var (
		department string
		age        int
	)
	for rows.Next() {
		err := rows.Scan(&department, &age)
		if err != nil {
			return "", 0, err
		}
	}

	return department, age, nil
}
```

If `Scan` fails on the second (or later) row, `get1` returns the values from the first row and `nil` error — silently hiding the failure.

## The fix

Always call `rows.Err()` after the `rows.Next()` loop finishes to detect any error that happened during iteration:

```go
func get2(ctx context.Context, db *sql.DB, id string) (string, int, error) {
	rows, err := db.QueryContext(ctx, "SELECT DEP, AGE FROM EMP WHERE ID = ?", id)
	if err != nil {
		return "", 0, err
	}
	defer rows.Close()

	var (
		department string
		age        int
	)
	for rows.Next() {
		err := rows.Scan(&department, &age)
		if err != nil {
			return "", 0, err
		}
	}
	if err := rows.Err(); err != nil { // <-- always check this
		return "", 0, err
	}

	return department, age, nil
}
```

## Rule of thumb

After every `for rows.Next()` loop, check `rows.Err()`. `rows.Next()` returning `false` means either "no more rows" or "an error occurred during iteration" — `rows.Err()` tells you which.

## Running the example

This is a standalone `package main`, but it references a real database, so it's intended mainly as a reference snippet:

```sh
go vet ./...
go build .
```

## Reference

* Teiva Harsanyi, *100 Go Mistakes and How to Avoid Them* — Mistake 78: SQL rows iteration errors.
* [`database/sql` package documentation](https://pkg.go.dev/database/sql#Rows)
