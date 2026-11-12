# Function input: accept `io.Reader`, not filenames or `*os.File`

This example accompanies chapter **46 — Function input** of *100 Go Mistakes and How to Avoid Them*.

## The mistake

Don't make your function open files itself. Accepting a filename (`string`) or a concrete `*os.File` couples your function to the filesystem and makes it hard to test:

```go
// BUG: tightly coupled to the filesystem, awkward to test
func countEmptyLinesInFile(filename string) (int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return 0, err
    }
    defer func() { _ = file.Close() }()

    scanner := bufio.NewScanner(file)
    count := 0
    for scanner.Scan() {
        if scanner.Text() == "" {
            count++
        }
    }
    if err := scanner.Err(); err != nil {
        return 0, err
    }
    return count, nil
}
```

A caller with data already in memory (a `strings.Reader`, a network response, an in-memory buffer) has to write to a temp file just to use your function. And unit tests need the filesystem.

## The fix

Accept an `io.Reader` — the most abstract input type. Let the caller decide where the data comes from and who opens/closes the file:

```go
// GOOD: works with any readable source — files, strings, networks, etc.
func countEmptyLines(reader io.Reader) (int, error) {
    scanner := bufio.NewScanner(reader)
    count := 0
    for scanner.Scan() {
        if scanner.Text() == "" {
            count++
        }
    }
    if err := scanner.Err(); err != nil {
        return 0, err
    }
    return count, nil
}
```

Usage:

```go
// From a file
file, err := os.Open("data.txt")
if err != nil { ... }
defer func() { _ = file.Close() }()
count, err := countEmptyLines(file)

// From a string — no filesystem needed
count, err := countEmptyLines(strings.NewReader("foo\n\nbar\n"))
```

## Why this matters

| Approach | Flexible | Testable | Caller controls lifecycle |
|----------|----------|----------|---------------------------|
| `filename string` | ❌ | ❌ needs temp files | ❌ function opens/closes |
| `*os.File` | ⚠️ only files | ⚠️ needs `os.Open` | ⚠️ caller opens, function closes |
| `io.Reader` | ✅ anything readable | ✅ `strings.NewReader`, `bytes.Buffer`, etc. | ✅ caller opens/closes |

This is the well-known Go mantra: **accept interfaces, return structs**. Accepting `io.Reader` makes your function flexible and testable; returning a concrete type (not an interface) gives callers a specific value they can use directly.

## Scanner errors

Both functions check `scanner.Err()` after the `Scan()` loop. `Scan()` returns `false` both at EOF and on I/O error, so you must call `scanner.Err()` to distinguish them — otherwise read errors are silently swallowed.

## Running

```bash
go vet ./...
go test -v ./...
```
