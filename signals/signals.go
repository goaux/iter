// Package signals provides utility functions for handling OS signals.
package signals

import (
	"context"
	"iter"
	"os"
	"os/signal"
	"sync/atomic"
)

// Wait returns an iterator that waits for the specified signals and passes the
// received signal to the loop body along with the loop index.
//
// The iterator yields the loop index and the received signal.
func Wait(ctx context.Context, signals ...os.Signal) iter.Seq2[int, os.Signal] {
	return func(yield func(int, os.Signal) bool) {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, signals...)
		defer signal.Stop(ch)
		for i := 0; ; i++ {
			select {
			case s := <-ch:
				if !yield(i, s) {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

// Context returns an iterator that creates a new context from the parent
// context each time in the loop, then passes it to the loop body. The context
// passed to the loop body will be cancelled when one of the specified signals
// is received.
//
// The iterator yields a new context and the loop index.
//
// Control passes immediately into the loop body.
// This is different from [Wait].
//
// After receiving a signal, [Get] returns it.
func Context(parent context.Context, signals ...os.Signal) iter.Seq2[context.Context, int] {
	return func(yield func(context.Context, int) bool) {
		parent, cancel := context.WithCancel(parent)
		defer cancel()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, signals...)
		defer signal.Stop(ch)
		loop := func(i int) bool {
			ctx, cancel := context.WithCancel(parent)
			var receive atomic.Value
			ctx = context.WithValue(ctx, getKey{}, &receive)
			go func() {
				select {
				case sig := <-ch:
					receive.Store(sig)
					cancel()
				case <-parent.Done():
				}
			}()
			ok := yield(ctx, i)
			return ok && !isDone(parent)
		}
		for i := 0; loop(i); i++ {
		}
	}
}

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

type getKey struct{}

// Get retrieves the signal stored in the context, if any.
// It returns nil if no signal is stored.
func Get(ctx context.Context) os.Signal {
	if va := ctx.Value(getKey{}); va != nil {
		if si := va.(*atomic.Value).Load(); si != nil {
			return si.(os.Signal)
		}
	}
	return nil
}
