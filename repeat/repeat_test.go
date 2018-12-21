package repeat_test

import (
	"context"
	"testing"
	"time"

	"github.com/m-lab/nodeinfo/repeat"
)

func TestRepeatOnce(t *testing.T) {
	count := 0
	f := func() { count++ }
	repeat.Forever(context.Background(), f, repeat.Config{Once: true})
	if count != 1 {
		t.Error("Once should mean once.")
	}
}

func TestRepeatForever(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	// We use count rather than a waitgroup because an extra call to f() shouldn't
	// cause the test to fail - cancel() races with the timer, and that's both
	// fundamental and okay. Contexts can be canceled() multiple times no problem,
	// but if you ever call .Done() on a WaitGroup more times than you .Add(), you
	// get a panic.
	count := 1000
	f := func() {
		if count < 0 {
			cancel()
		} else {
			count--
		}
	}
	wt := time.Duration(1 * time.Microsecond)
	go repeat.Forever(ctx, f, repeat.Config{Expected: wt, Min: wt, Max: wt})
	<-ctx.Done()
	// If this does not run forever, then f() was called at least 100 times and
	// then the context was canceled.
}
