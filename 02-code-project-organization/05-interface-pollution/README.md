# Interface Pollution

This directory contains examples from **Mistake #5** of *100 Go Mistakes*: **interface pollution**.

The core idea is simple: keep Go interfaces small, focused, and driven by the **consumer's** needs. Avoid large, "kitchen-sink" interfaces that force implementers to write unused methods and let callers reach into behavior they shouldn't touch.

## Contents

- [1. `copy/` — Narrow interfaces are not pollution](#1-copy--narrow-interfaces-are-not-pollution)
- [2. `decoupling/` — Interface vs. concrete type](#2-decoupling--interface-vs-concrete-type)
- [3. `restricting-behavior/` — Interface segregation](#3-restricting-behavior--interface-segregation)
- [Types in Go](#types-in-go)
- [Type Sets (Go 1.18+)](#type-sets-go-118)
- [Running the examples](#running-the-examples)

---

The directory is split into three examples:

| Directory | What it shows |
|-----------|---------------|
| `copy/` | A clean baseline using narrow standard-library interfaces |
| `decoupling/` | Coupling to a concrete type vs. decoupling via a small interface |
| `restricting-behavior/` | Interface segregation — exposing only what the consumer needs |

---

## 1. `copy/` — Narrow interfaces are not pollution

```go
func copySourceToDest(source io.Reader, dest io.Writer) error {
	b, err := io.ReadAll(source)
	if err != nil {
		return err
	}
	_, err = dest.Write(b)
	return err
}
```

This uses `io.Reader` and `io.Writer`, two tiny interfaces from the standard library. Each describes exactly one capability. Any type with a single matching method satisfies them. This is **good** interface design: the caller depends on behavior, not on a concrete data shape.

**Rule of Thumb: interfaces in, concrete types out.**
Accept interfaces as parameters (maximum flexibility), and return concrete types (callers know exactly what they're getting without type assertions).

The test in `main_test.go` simply verifies that data flows from a `strings.Reader` into a `bytes.Buffer`.

### Exceptions to the rule

"Concrete Types" is a heuristic, not a law. Returning an interface is appropriate when:

- **Hiding implementation details (factory pattern).** You don't want callers to depend on which concrete type you chose:
  ```go
  func NewLogger() Logger { return &fileLogger{} }
  ```
- **Returning one of several implementations.** The caller genuinely needs to treat multiple concrete types uniformly:
  ```go
  func Open(name string) (io.ReadWriteCloser, error)
  ```
- **The concrete type is unexported.** If the implementation type lives in an internal package, returning it would force callers to use an interface anyway.

When in doubt, prefer concrete return types, but don't twist the API to avoid an interface if an interface is the honest description of what you're returning.

---

## 2. `decoupling/` — Interface vs. concrete type

### Without an interface (`without.go`)

```go
type CustomerService struct {
	store Store
}

func (cs CustomerService) CreateNewCustomer(id string) error {
	customer := Customer{id: id}
	return cs.store.StoreCustomer(customer)
}

type Store struct{}
func (s Store) StoreCustomer(customer Customer) error { return nil }
```

`CustomerService` is locked to the concrete type `Store`. You cannot test it with a fake, or swap in a database, without changing the service itself.

### With an interface (`with.go`)

```go
type customerStorer interface {
	StoreCustomer(Customer) error
}

type CustomerService2 struct {
	storer customerStorer
}

func (cs CustomerService2) CreateNewCustomer(id string) error {
	customer := Customer{id: id}
	return cs.storer.StoreCustomer(customer)
}
```

Now `CustomerService2` depends on **behavior** (`customerStorer`), not a concrete type. You can inject a mock in tests, a SQL-backed store in production, or an in-memory store in a CLI tool — all without touching the service code.

---

## 3. `restricting-behavior/` — Interface segregation

`IntConfig` exposes both `Get()` and `Set()`, but `Foo` only ever needs to **read** the threshold:

```go
type IntConfig struct {
	value int
}

func (c *IntConfig) Get() int     { return c.value }
func (c *IntConfig) Set(value int) { c.value = value }

type intConfigGetter interface {
	Get() int
}

type Foo struct {
	threshold intConfigGetter
}

func NewFoo(threshold intConfigGetter) Foo {
	return Foo{threshold: threshold}
}
```

`Foo` receives an `intConfigGetter`, so through that interface it can **only** call `Get()`. Even though the underlying `*IntConfig` supports `Set()`, `Foo` physically cannot call it. The consumer sees only what it needs — nothing more.

---

## Types in Go

A **type** in Go is a set of values plus the operations that can be performed on them. Every value has a type, and the type determines what you can do with that value.

### Named vs. unnamed types

- **Named types** are declared with the `type` keyword:
  ```go
  type Age int
  type Customer struct { Name string }
  type Reader interface { Read([]byte) (int, error) }
  ```
- **Unnamed (composite) types** are written inline:
  ```go
  []int                // slice
  map[string]int       // map
  [3]int               // array
  struct{ name string } // anonymous struct
  func(int) int        // function type
  ```

Two named types are different types even if they share the same underlying representation:

```go
type Age int
type Year int

var a Age = 1
var y Year = 1
// a = y  // compile error: cannot use y (type Year) as type Age
```

### Basic ("primitive") types

| Category | Types |
|----------|-------|
| Booleans | `bool` |
| Integers | `int`, `int8`–`int64`, `uint`, `uint8`–`uint64`, `uintptr` |
| Floats | `float32`, `float64` |
| Complex | `complex64`, `complex128` |
| Strings | `string` |
| Rune / Byte | `rune` (alias for `int32`), `byte` (alias for `uint8`) |

There is no `char` type; `byte` is for raw bytes and `rune` is for Unicode code points.

### Composite types

- **Arrays**: `[N]T` — fixed-size sequence
- **Slices**: `[]T` — dynamic, flexible view into an array
- **Maps**: `map[K]V` — hash table
- **Structs**: `struct { ... }` — grouped fields; Go's replacement for classes
- **Pointers**: `*T`
- **Functions**: `func(...) (...)`
- **Channels**: `chan T`, `chan<- T`, `<-chan T`
- **Interfaces**: `interface { ... }`

### Interface types

An interface is an abstract type defined by a **method set**, not a data layout. Any type that implements the methods satisfies the interface implicitly — there is no `implements` keyword.

Good Go style keeps interfaces small and consumer-driven (see `io.Reader`, `customerStorer`, `intConfigGetter` above).

### Type aliases vs. new types

```go
type Duration = time.Duration   // alias: Duration IS time.Duration
type MyDuration time.Duration   // new named type (distinct from time.Duration)
```

- `=` creates an **alias** — both names refer to the exact same type.
- No `=` creates a **new defined type** — same underlying type, but a different type for assignment.

### Underlying type

Every type has an **underlying type**. For `type Age int`, the underlying type is `int`. For structs, slices, maps, and interfaces, the underlying type is itself. This matters for type assertions, comparisons, and generic constraints.

### Zero values

Every type has a zero value: `0` for numbers, `""` for strings, `false` for booleans, `nil` for pointers/slices/maps/channels/functions/interfaces, and all-zero fields for structs. Variables are initialized to their zero value if you don't assign one.

---

## Type Sets (Go 1.18+)

A **type set** is the set of all Go types that satisfy an interface.

Before Go 1.18, an interface's type set was defined only by its method set (e.g. any type with a `Read` method satisfies `io.Reader`). Since Go 1.18, an interface can also restrict which **underlying types** are allowed using:

- `~` — "underlying type is"
- `|` — union ("or")

### Example

```go
interface {
	~int | ~int64
	String() string
}
```

| Piece | Meaning |
|-------|---------|
| `~int` | Any type whose underlying type is `int` (e.g. `int`, `type MyInt int`) |
| `~int64` | Any type whose underlying type is `int64` |
| `\|` | Either one |
| `String() string` | The type must also implement this method |

So the type set is: **all types whose underlying type is `int` or `int64`, and that have a `String() string` method.**

### Concrete examples

```go
type MyInt int
type MyInt64 int64
type MyFloat float64

func (m MyInt) String() string       { return "my int" }
func (m MyInt64) String() string     { return "my int64" }
func (m MyFloat) String() string     { return "my float" }
```

| Type | Underlying type | Has `String()`? | Satisfies the interface? |
|------|----------------|-----------------|--------------------------|
| `MyInt` | `int` | Yes | ✅ Yes |
| `MyInt64` | `int64` | Yes | ✅ Yes |
| `MyFloat` | `float64` | Yes | ❌ No (wrong underlying type) |
| `int` | `int` | No | ❌ No (missing `String()`) |

### The `~` tilde matters

Without `~`, `int` means **exactly** the named type `int`:

```go
type MyInt int

var _ interface{ int } = MyInt(0)   // ❌ compile error
var _ interface{ ~int } = MyInt(0)  // ✅ OK
```

### Real-world example: `constraints.Ordered`

The standard library's `golang.org/x/exp/constraints` defines:

```go
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
	~float32 | ~float64 | ~string
}
```

This enables generic functions over ordered types:

```go
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

Min(3, 5)        // T = int
Min(3.5, 2.1)    // T = float64
Min("b", "a")    // T = string
```

### Is this like a Haskell type class?

In spirit, yes: both provide **ad-hoc polymorphism** (functions that work over many types supporting certain operations).

| | Haskell type class (`Ord`) | Go interface constraint (`constraints.Ordered`) |
|---|---|---|
| Satisfaction | Nominal / explicit (`instance Ord MyType where ...`) | Structural / implicit (based on underlying type + methods) |
| Extensibility | Open — anyone can add instances | Closed — fixed list of underlying types |
| What is constrained | Methods in a dictionary | Built-in operators on underlying types |
| Higher kinds | Yes (`Functor`, `Monad`) | No |

You **cannot** add your own type to `constraints.Ordered`; there is no `instance` declaration in Go. Think of Go's generic constraints as "lightweight, structural type classes" — useful and powerful, but more limited than Haskell's.

---

## Running the examples

```bash
# From the copy/ directory
go test ./...

# From restricting-behavior/
go run main.go
```

The other directories (`decoupling/`) are illustrative snippets; they require a shared `Customer`/`Store` definition to compile as a standalone program.
