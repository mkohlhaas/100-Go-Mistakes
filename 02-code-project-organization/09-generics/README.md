# 09 — Generics

This sample is from Chapter 2 ("Code Project Organization"), Mistake #9 of *100 Go Mistakes* by Teiva Harsanyi.

## What it demonstrates

It contrasts the pre-generics approach (`any` / `interface{}` with a type switch) with Go 1.18+ generics for writing type-safe algorithms and data structures that work across multiple types.

### The mistake — `getKeys(m any)`

Before generics, extracting map keys required an `any` parameter and a manual type switch:

```go
func getKeys(m any) ([]any, error) {
	switch t := m.(type) {
	default:
		return nil, fmt.Errorf("unknown type: %T", t)
	case map[string]int:
		var keys []any
		for k := range t {
			keys = append(keys, k)
		}
		return keys, nil
	case map[int]string:
		// ...
	}
	return nil, nil
}
```

Problems:
- **Not scalable** — you must add a new `case` for every map key/value type.
- **Loss of type safety** — returns `[]any`, forcing callers to type-assert again.
- **Boilerplate** — the same loop logic is repeated per type.

### The generics fix — `getKeysGenerics`

```go
func getKeysGenerics[K comparable, V any](m map[K]V) []K {
	var keys []K
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
```

One function handles **any** map type:
- `K comparable` — keys must be comparable (a requirement for map keys).
- `V any` — values are unconstrained.
- Returns `[]K`, preserving type safety — no type assertions needed at the call site.

### Custom constraints — `customConstraint`

```go
type customConstraint interface {
	~int | ~string
}
```

Shows Go's union (`|`) and underlying-type (`~`) operators to constrain a type parameter to a specific set of types, including named types whose underlying type is `int` or `string`.

### Generic structs and methods — `Node[T]`

```go
type Node[T any] struct {
	Val  T
	next *Node[T]
}

func (n *Node[T]) Add(next *Node[T]) { n.next = next }
```

A generic linked-list node where the value type is a type parameter, including a generic method.

### Generics + interfaces — `SliceFn[T]`

```go
type SliceFn[T any] struct {
	S       []T
	Compare func(T, T) bool
}

func (s SliceFn[T]) Len() int           { return len(s.S) }
func (s SliceFn[T]) Less(i, j int) bool { return s.Compare(s.S[i], s.S[j]) }
func (s SliceFn[T]) Swap(i, j int)      { s.S[i], s.S[j] = s.S[j], s.S[i] }
```

A generic type satisfying `sort.Interface` by injecting a comparison function. In `main`, it sorts a `[]int` via `sort.Sort`.

## Bottom line

Prefer generics over `interface{}` / `any` when you write algorithms or data structures that should work across multiple types. Generics give you type safety, less boilerplate, and no runtime type assertions — while still allowing custom constraints and composition with interfaces.

## Run it

```sh
go run main.go
# Output: [1 2 3]
```