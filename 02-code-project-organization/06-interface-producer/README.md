# Interface Producer — Mistake #6

This directory demonstrates **Mistake #6 from Chapter 2 ("Code Project Organization") of *100 Go Mistakes*: Producing interfaces** — acting as an "interface producer" instead of letting consumers define their own interfaces.

## The Mistake

**`incorrect/`** shows the producer (`store`) defining and exporting a large, monolithic interface:

```go
type CustomerStorage interface {
    StoreCustomer(customer Customer) error
    GetCustomer(id string) (Customer, error)
    UpdateCustomer(customer Customer) error
    GetAllCustomers() ([]Customer, error)
    GetCustomersWithoutContract() ([]Customer, error)
    GetCustomersWithNegativeBalance() ([]Customer, error)
}
```

This violates the **Interface Segregation Principle**: consumers are forced to depend on a bloated interface even when they only need one or two methods. The producer shouldn't guess what interfaces consumers need.

## The Fix

**`correct/`** follows the Go maxim: **"Accept interfaces, return structs."**

- **Producer (`store`)** exports a concrete `Store` type with methods — no interface.
- **Consumer (`client`)** defines its own small, focused `customersGetter` interface with only `GetAllCustomers`, close to where it's used.
- **`NewClient`** wires them together via dependency injection.

### Key differences

| Aspect | `incorrect/` (mistake) | `correct/` (fix) |
|---|---|---|
| Who defines the interface? | Producer (`store`) | Consumer (`client`) |
| Interface size | 7 methods (bloated) | 1 method (minimal) |
| Producer exports | `CustomerStorage` interface | Concrete `Store` struct |
| Coupling | Tight to producer's interface | Loose via structural typing |
| Principle | Violates Interface Segregation | Follows "Accept interfaces, return structs" |

## Directory Structure

```
incorrect/   The mistake — producer defines a bloated CustomerStorage interface
correct/     The fix — producer returns a concrete Store; client defines its own small interface
```

## Running

Each folder is a standalone Go module:

```bash
# The mistake
cd incorrect
go build ./...
go vet ./...

# The fix
cd correct
go build ./...
go vet ./...
go run .      # demonstrates NewClient wiring Store into Client
```

## Takeaway

Don't define big interfaces in a package and expect other packages to implement or depend on them. Let each consumer define the minimal interface it requires, and have producers export concrete types rather than imposing a single fat interface on everyone.
