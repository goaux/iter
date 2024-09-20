# iter

`iter` is a Go package that provides utility functions and types for working with iterators in Go (since go1.23).

[![Go Reference](https://pkg.go.dev/badge/github.com/goaux/iter.svg)](https://pkg.go.dev/github.com/goaux/iter)
[![Go Report Card](https://goreportcard.com/badge/github.com/goaux/iter)](https://goreportcard.com/report/github.com/goaux/iter)

Iterator utilities are provided in several packages, each of which serves a different purpose.

Sub-packages:

- `bufioreader` and `bufioscanner`, which offer convenient ways to iterate over buffered I/O operations.
- `signals` provides iterators for the os.Signal event loop.
- `ticker` provides iterators for the time.Ticker event loop.
- `transform` provides functions transforming iterators.

## Sub-packages

### bufioreader

The `bufioreader` package provides a wrapper around `bufio.Reader` with additional iteration capabilities.

#### Features:
- Iterate over bytes, slices, or strings delimited by a specified byte
- Error handling through `Err()`, `GetErrorBuffer()` and `GetErrorBufferString()` methods

#### Example usage:

```go
import "github.com/goaux/iter/bufioreader"

r := bufioreader.NewReader(strings.NewReader("hello\nworld"))
for i, s := range r.ReadString('\n') {
    fmt.Printf("[%d] %q\n", i, s)
}
if err := r.Err(); err != nil {
    // r.Err never returns io.EOF.
    // Under normal circumstances, bufio.Reader will return io.EOF at the end,
    // but this is to indicate the end of the data and is not an error.
    fmt.Printf("error: %v (remain:%q)\n", err, bufioreader.GetErrorBufferString(err))
}
```

### bufioscanner

The `bufioscanner` package provides a wrapper around `bufio.Scanner` with iteration capabilities.

#### Features:
- Iterate over bytes or strings using a `bufio.Scanner`
- Compatible with different `Split` functions

#### Example usage:

```go
import "github.com/goaux/iter/bufioscanner"

s := bufioscanner.NewScanner(bytes.NewBufferString("Where are you going\nfor your next vacation?"))
s.Split(bufio.ScanWords)
for i, word := range s.Text() {
    fmt.Printf("[%d] %q\n", i, word)
}
if err := s.Err(); err != nil {
    fmt.Println("Error:", err)
}
```

### signals

The `signals` package provides an iterator for receiving signals.

#### Example usage:

```go
import "github.com/goaux/iter/signals"

for i, si := range signals.Wait(context.TODO(), syscall.SIGHUP, syscall.SIGINT) {
    // When a signal is received, control passes to the loop body.
    // In this example, the loop body will be entered upon receiving either a SIGHUP or a SIGINT.
    fmt.Println(i, si) // Outputs the loop index and the received signal.
    if si == syscall.SIGINT {
        break // The loop continues until break.
    }
}

for ctx, i := range signals.Context(ctx, syscall.SIGUSR1, syscall.SIGUSR2) {
    // Control passes immediately to the loop body.
    fmt.Printf("%d enter body.\n", i)
    <-ctx.Done()
    // When a signal is received, ctx is cancelled.
    // In this example, receiving either SIGUSR1 or SIGUSR2 will cancel the context.
    fmt.Printf("%d leave body. signal:%v\n", i, signals.Get(ctx))
    if signals.Get(ctx) == syscall.SIGUSR1 {
      break // The loop continues until break.
    }
}
```

### ticker

The `ticker` package provides iterators for timed loops using Go's iter package.

It offers two main functions: After and Before, which allow you to create
timed loops with different behaviors.

#### Features

- After: Creates an iterator that executes the loop body after waiting for a specified duration.
- Before: Creates an iterator that executes the loop body before waiting for a specified duration.
- Both functions support context cancellation for immediate loop termination.

#### Example usage:

```go
import "github.com/goaux/iter/ticker"

// This is a ticker event loop.
// The ticker event received, enter the loop body.
// Ticker events occur `after` a specified duration.
start := time.Now()
for i, now := range ticker.After(ctx, time.Second) {
    // Entering the loop body after waiting a specified duration.
    elapse := now.Sub(start)
    fmt.Println(i, elapse.String())
    if i >= 9 { // Since i starts from 0, this means to break after looping 10 times.
        break
    }
}
```

```go
import "github.com/goaux/iter/ticker"

// This is a ticker event loop.
// The ticker event received, enter the loop body.
// Ticker events occur `before` a specified duration.
start := time.Now()
for i, now := range ticker.Before(ctx, time.Second) {
    // Entering the loop body before waiting a specified duration at first.
    elapse := now.Sub(start)
    fmt.Println(i, elapse.String())
    if i >= 9 { // Since i starts from 0, this means to break after looping 10 times.
        break
    }
}
```

### transform

The `transform` package provides functions transforming iterators.

```go
func Concat[T any](iterators ...iter.Seq[T]) iter.Seq[T]
func Concat2[S, T any](iterators ...iter.Seq2[S, T]) iter.Seq2[S, T]
func Keys[K, V any](iterator iter.Seq2[K, V]) iter.Seq[K]
func Map[S, T any](iterator iter.Seq[S], f func(S) T) iter.Seq[T]
func Map2[S, T, U, V any](iterator iter.Seq2[S, T], f func(S, T) (U, V)) iter.Seq2[U, V]
func MapIn[S, T, U any](iterator iter.Seq2[S, T], f func(S, T) U) iter.Seq[U]
func MapOut[S, T, U any](iterator iter.Seq[S], f func(S) (T, U)) iter.Seq2[T, U]
func Resize[T any](iterator iter.Seq[T], size int) iter.Seq[T]
func Resize2[S, T any](iterator iter.Seq2[S, T], size int) iter.Seq2[S, T]
func Select[T any](iterator iter.Seq[T], f func(T) bool) iter.Seq[T]
func Select2[K, V any](iterator iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V]
func Skip[T any](iterator iter.Seq[T], skip int) iter.Seq[T]
func Skip2[S, T any](iterator iter.Seq2[S, T], skip int) iter.Seq2[S, T]
func Swap[S, T any](iterator iter.Seq2[S, T]) iter.Seq2[T, S]
func Values[K, V any](iterator iter.Seq2[K, V]) iter.Seq[V]
func Zip[S, T any](lhs iter.Seq[S], rhs iter.Seq[T]) iter.Seq2[S, T]
func ZipAll[S, T any](lhs iter.Seq[S], rhs iter.Seq[T]) iter.Seq2[S, T]
func ZipIndex[T any](iterator iter.Seq[T]) iter.Seq2[int, T]
func ZipLeft[S, T any](lhs iter.Seq[S], rhs iter.Seq[T]) iter.Seq2[S, T]
func ZipRight[S, T any](lhs iter.Seq[S], rhs iter.Seq[T]) iter.Seq2[S, T]
```
