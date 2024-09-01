package signals_test

import (
	"context"
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/goaux/iter/signals"
)

func ExampleWait() {
	go func() {
		pid := syscall.Getpid()
		time.Sleep(200 * time.Millisecond)
		fmt.Println("SIGHUP")
		syscall.Kill(pid, syscall.SIGHUP)
		time.Sleep(100 * time.Millisecond)
		fmt.Println("SIGINT")
		syscall.Kill(pid, syscall.SIGINT)
	}()
	for i, si := range signals.Wait(context.TODO(), syscall.SIGHUP, syscall.SIGINT) {
		fmt.Println(i, si)
		if si == syscall.SIGINT {
			break
		}
	}
	// Output:
	// SIGHUP
	// 0 hangup
	// SIGINT
	// 1 interrupt
}

func ExampleContext() {
	go func() {
		pid := syscall.Getpid()
		time.Sleep(300 * time.Millisecond)
		fmt.Println("SIGUSR1")
		syscall.Kill(pid, syscall.SIGUSR1)
		time.Sleep(100 * time.Millisecond)
		fmt.Println("SIGUSR2")
		syscall.Kill(pid, syscall.SIGUSR2)
		time.Sleep(100 * time.Millisecond)
		fmt.Println("SIGINT")
		syscall.Kill(pid, syscall.SIGINT)
	}()
	for ctx := range signals.Context(context.TODO(), syscall.SIGINT) {
		fmt.Println("start loop")
		for ctx, i := range signals.Context(ctx, syscall.SIGUSR1, syscall.SIGUSR2) {
			fmt.Printf("%d enter body. signal:%v\n", i, signals.Get(ctx))
			<-ctx.Done()
			fmt.Printf("%d leave body. signal:%v\n", i, signals.Get(ctx))
		}
		fmt.Printf("finish loop. signal:%v\n", signals.Get(ctx))
		break
	}
	// Output:
	// start loop
	// 0 enter body. signal:<nil>
	// SIGUSR1
	// 0 leave body. signal:user defined signal 1
	// 1 enter body. signal:<nil>
	// SIGUSR2
	// 1 leave body. signal:user defined signal 2
	// 2 enter body. signal:<nil>
	// SIGINT
	// 2 leave body. signal:<nil>
	// finish loop. signal:interrupt
}

func TestWait(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()
	i := 0
	for range signals.Wait(ctx, syscall.SIGINT) {
		i++
	}
	if i != 0 {
		t.Errorf("i must be 0, but i is %d", i)
	}
}

func TestContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()
	i := 0
	for ctx := range signals.Context(ctx, syscall.SIGINT) {
		if si := signals.Get(ctx); si != nil {
			t.Errorf("si must be nil, but got %v", si)
		}
		i++
		<-ctx.Done()
		if si := signals.Get(ctx); si != nil {
			t.Errorf("si must be nil, but got %v", si)
		}
	}
	if i != 1 {
		t.Errorf("i must be 1, but i is %d", i)
	}
}

func TestGet(t *testing.T) {
	got := signals.Get(context.TODO())
	if got != nil {
		t.Errorf("signals.Get must return nil, but got=%v", got)
	}
}
