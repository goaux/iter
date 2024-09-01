// Package ticker provides iterators for timed loops using Go's [iter] package.
//
// It offers two main functions: [After] and [Before], which allow you to create
// timed loops with different behaviors.
//
// Features
//
//   - [After]: Creates an iterator that executes the loop body after waiting for a specified duration.
//   - [Before]: Creates an iterator that executes the loop body before waiting for a specified duration.
//   - Both functions support context cancellation for immediate loop termination.
package ticker

import (
	"context"
	"iter"
	"time"
)

// After returns an iterator that passes control to the loop body
// after waiting for duration d.
//
// This iterator waits first and then executes the loop body.
//
// The iterator passes the loop index and the current time returned by internal
// [time.Ticker] to the loop body.
//
// The loop and wait can be interrupted immediately by canceling the context.
func After(ctx context.Context, d time.Duration) iter.Seq2[int, time.Time] {
	return after(ctx, d, 0)
}

// Before returns an iterator that passes control to the loop body
// before waiting for duration d.
//
// This iterator executes the loop body first and then waits.
//
// The iterator passes the loop index and the current time returned by internal
// [time.Ticker] to the loop body.
//
// The loop and wait can be interrupted immediately by canceling the context.
func Before(ctx context.Context, d time.Duration) iter.Seq2[int, time.Time] {
	return func(yield func(int, time.Time) bool) {
		select {
		case <-ctx.Done():
			return
		default:
			if !yield(0, time.Now()) {
				return
			}
		}
		after(ctx, d, 1)(yield)
	}
}

func after(ctx context.Context, d time.Duration, i int) iter.Seq2[int, time.Time] {
	return func(yield func(int, time.Time) bool) {
		t := time.NewTicker(d)
		defer t.Stop()
		for ; ; i++ {
			select {
			case v := <-t.C:
				if !yield(i, v) {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}
}
