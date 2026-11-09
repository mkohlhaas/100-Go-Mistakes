# Utility Packages

This directory contains examples from **Mistake #13** of *100 Go Mistakes*: **utility packages** — catch-all packages that aggregate loosely-related helpers under vague names like `utils`, `helpers`, `common`, `util`, or `misc`.

The core idea is simple: **organize Go packages around domain responsibility, not around technical categories.** A package named after its job ("utility") rather than its purpose tells callers nothing about *what* the code does or *which feature* it serves.

---

## The example

**`stringset.go`**

```go
// Package stringset
package stringset

type Set map[string]struct{}

func New(...string) Set { return nil }

func (s Set) Sort() []string { return nil }
```

This is a stub of a string-set "utility" package. The method bodies return `nil`, so this is purely illustrating package **structure and naming**, not a working data structure.

### What's wrong here

- **No clear domain responsibility.** `stringset` doesn't belong to any business concept; it's just a grab-bag type (`Set`) plus generic operations (`New`, `Sort`) that could apply anywhere.
- **A dumping ground.** This is exactly how `utils`/`common` packages grow: anything that "doesn't have a better home" gets tossed in, so cohesion erodes and the package becomes hard to reason about.
- **Uninformative imports.** When you `import "utility-packages/stringset"`, the import path tells you nothing about the feature the code serves — only that it's a generic utility.
- **Poor package organization.** Go packages should be organized around domain concepts or features, not around technical categories like "string utilities."

### Better alternatives

Instead of a standalone utility package:

- Put a string set in the package/domain that actually uses it.
- If several features share the string set, give it a well-scoped, purposeful package name that reflects that shared domain, not its technical category.

**Rule of Thumb:** avoid package names like `util`, `utils`, `helpers`, `common`, or `misc`. Name packages after what they do for the business, not after the fact that you couldn't think of a better home for the code.
