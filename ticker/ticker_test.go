package ticker_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/goaux/iter/ticker"
)

func ExampleAfter() {
	unit := 100 * time.Millisecond
	start := time.Now()
	for i, now := range ticker.After(context.TODO(), 3*unit) {
		elapse := now.Sub(start)
		fmt.Println(i, ((elapse / unit) * unit).String())
		if elapse >= 9*unit {
			break
		}
	}
	// Output:
	// 0 300ms
	// 1 600ms
	// 2 900ms
}

func ExampleBefore() {
	unit := 100 * time.Millisecond
	start := time.Now()
	for i, now := range ticker.Before(context.TODO(), 3*unit) {
		elapse := now.Sub(start)
		fmt.Println(i, ((elapse / unit) * unit).String())
		if elapse >= 9*unit {
			break
		}
	}
	// Output:
	// 0 0s
	// 1 300ms
	// 2 600ms
	// 3 900ms
}

func TestBefore(t *testing.T) {
	t.Run("break immediately", func(t *testing.T) {
		i := 0
		for range ticker.Before(context.TODO(), time.Second) {
			i++
			break
		}
		if i != 1 {
			t.Errorf("i must be 1, but %d", i)
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		cancel()
		i := 0
		for range ticker.Before(ctx, time.Second) {
			i++
		}
		if i != 0 {
			t.Errorf("i must be 0, but %d", i)
		}
	})

	t.Run("cancel context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()
		i := 0
		for range ticker.Before(ctx, time.Second) {
			i++
			cancel()
		}
		if i != 1 {
			t.Errorf("i must be 1, but %d", i)
		}
	})
}
