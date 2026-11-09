# 10 — Type Embedding

This sample demonstrates **Go type embedding** and the pitfalls of embedding types whose methods you don't want to leak into your public API.

## What's in the file

`main.go` contains three examples:

### 1. Struct embedding (promoted fields)

```go
type Foo struct {
	Bar
}

type Bar struct {
	Baz int
}

func fooBar() {
	foo := Foo{}
	foo.Baz = 42   // Bar's field is "promoted" to Foo
}
```

`Foo` embeds `Bar` anonymously, so `Bar`'s `Baz` field is promoted onto `Foo`. You can read/write `foo.Baz` directly as if it belonged to `Foo`.

### 2. The mistake: embedding `sync.Mutex`

```go
type InMem struct {
	sync.Mutex   // embedded anonymously
	m map[string]int
}

func (i *InMem) Get(key string) (int, bool) {
	i.Lock()
	v, contains := i.m[key]
	i.Unlock()
	return v, contains
}
```

`InMem` embeds `sync.Mutex` anonymously. Because it's an embedded (unnamed) field, `sync.Mutex`'s `Lock()` and `Unlock()` methods are **promoted** onto `InMem`'s public method set. Any caller outside the package can now do:

```go
inMem.Lock()   // oops — I can grab your internal mutex
```

This leaks your internal synchronization primitive into your public API. Callers can deadlock your type, unlock a mutex they never locked, or otherwise break the invariants your methods rely on (`Get` assumes it alone holds the lock).

**The fix:** use an unexported *named* field instead:

```go
type InMem struct {
	mu sync.Mutex
	m  map[string]int
}
```

Then only methods inside the package can call `mu.Lock()`/`mu.Unlock()`.

A secondary subtlety: `sync.Mutex` must not be copied after its first use, and embedding it by value means copying an `InMem` copies the mutex — another reason to prefer a named (pointer or value) field rather than embedding.

### 3. Explicit delegation vs. implicit promotion

```go
type Logger struct {
	writeCloser io.WriteCloser   // NOTE: named field, NOT embedded
}

func (l Logger) Write(p []byte) (int, error) {
	return l.writeCloser.Write(p)
}

func (l Logger) Close() error {
	return l.writeCloser.Close()
}
```

Here `writeCloser` is a **named** field, so `io.WriteCloser`'s methods are *not* promoted onto `Logger`. `Logger` explicitly implements `Write` and `Close` and delegates to the wrapped value. This keeps full control over the API — nothing leaks implicitly.

## Bottom line

- Embedding promotes fields and methods (handy for reuse, but implicit).
- Embedding `sync.Mutex` (or any type with methods you want to keep private) exposes those methods publicly — a common encapsulation bug.
- Prefer unexported named fields for mutexes, and explicit delegation when you want to control exactly what your type exposes.
