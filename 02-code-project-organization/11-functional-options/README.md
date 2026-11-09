# 11 — Functional Options

This sample is from Chapter 2 ("Code Project Organization"), Mistake #11 of *100 Go Mistakes* by Teiva Harsanyi.

It shows the **Java-style mutable builder pattern applied to Go** — the mistake the chapter warns against — and contrasts it with the idiomatic Go alternative: **functional options**.

## What's in the directory

Three variants live side by side:

| Directory | Approach |
|---|---|
| `builder/` | Mutable builder with a separate `Build()` step (**the mistake**) |
| `config-struct/` | Passing a plain `Config` struct directly (too inflexible) |
| `functional-options/` | Variadic functional options (`...Option`) (**the fix**) |

## The mistake — `builder/main.go`

`builder/main.go` defines a `Config` with a single `Port` field, plus a `ConfigBuilder` that mutates state and exposes a separate `Build()` step:

```go
type ConfigBuilder struct {
	port *int
}

func (b *ConfigBuilder) Port(port int) *ConfigBuilder {
	b.port = &port
	return b
}

func (b *ConfigBuilder) Build() (Config, error) {
	cfg := Config{}

	if b.port == nil {
		cfg.Port = defaultHTTPPort
	} else {
		if *b.port == 0 {
			cfg.Port = randomPort()
		} else if *b.port < 0 {
			return Config{}, errors.New("port should be positive")
		} else {
			cfg.Port = *b.port
		}
	}

	return cfg, nil
}
```

A caller must use it in two steps:

```go
builder := ConfigBuilder{}
builder.Port(8080)            // mutates state
cfg, err := builder.Build()   // separate build step
```

### What's wrong with it

1. **Mutable state.** `ConfigBuilder` is a mutable object. `Port()` mutates it and `Build()` consumes it, so it isn't concurrency-safe and it's easy to reuse a builder incorrectly.
2. **Two-step construction.** Callers must remember to call `Build()` after setting options — an extra ceremony step with no benefit in Go.
3. **Duplicated logic.** Validation and defaulting live inside `Build()` (and are duplicated again in the functional-options fix), rather than being expressed close to where each option is supplied.
4. **Unidiomatic in Go.** Go doesn't use method-chaining builders as a convention. The standard library and most Go libraries (e.g., `http.Server`, gRPC `DialOption`) use **functional options** instead.

## The fix — `functional-options/main.go`

Configuration is expressed as variadic functions that mutate a private `options` struct:

```go
type Option func(options *options) error

func WithPort(port int) Option {
	return func(options *options) error {
		if port < 0 {
			return errors.New("port should be positive")
		}
		options.port = &port
		return nil
	}
}

func NewServer(addr string, opts ...Option) (*http.Server, error) {
	var options options
	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return nil, err
		}
	}
	// apply defaults / finalize config here
	...
}
```

Usage collapses to a single call:

```go
_, _ = NewServer("localhost", WithPort(8080))
```

This keeps validation and defaulting in one place, needs no mutable builder, lets callers skip options entirely (`NewServer("localhost")`), and matches how Go libraries express optional configuration.

## Bottom line

- Don't port Java-style mutable builders to Go.
- Use **functional options** (`Option func(...)` with variadic `...Option`) for optional configuration.
- Keep validation/defaulting where each option is applied; avoid a separate mutable `Build()` step.

## Run it

From this directory:

```sh
go run builder/main.go       # the mistake (builder)
go run config-struct/main.go # plain struct
go run functional-options/main.go # the fix (functional options)
```

Each `main.go` is a standalone `package main`, so run them individually.
