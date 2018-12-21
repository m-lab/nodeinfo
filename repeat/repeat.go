package repeat

import (
	"context"
	"math/rand"
	"time"
)

// Config represents the time we should wait. "Once" is provided as a helper,
// because frequently for unit testing and integration testing, you only want
// the "Forever" loop to run once.
//
// The zero value of this struct has Once set to false, which means the value
// only needs to be set explicitly in codepaths where it might be true.
type Config struct {
	Expected, Min, Max time.Duration
	Once               bool
}

func (c Config) waittime() time.Duration {
	wt := time.Duration(rand.ExpFloat64() * float64(c.Expected))
	if wt < c.Min {
		wt = c.Min
	}
	if c.Max != 0 && wt > c.Max {
		wt = c.Max
	}
	return wt
}

// Forever calls the given function repeatedly, waiting a c.Expected amount of
// time between calls on average. The wait time is actually random and will
// generate a memoryless (Poisson) distribution of f() calls in time, ensuring
// that f() has the PASTA property (Poisson Arrivals See Time Averages). This
// statistical guarantee is subject to two caveats.
//
// Caveat 1 is that, in a nod to the realities of systems needing to have
// guarantees, we allow the random wait time to be clamped both above and below.
// This means that calls to f() should be at least c.Min and at most c.Max
// apart in time. This clamping causes bias in the timing. For use of this
// function to be statistically sensible, the clamping should not be too
// extreme. The exact mathematical meaning of "too extreme" depends on your
// situation, but a nice rule of thumb is c.Min should be at most 10% of expected
// and c.Max should be at least 250% of expected. These values mean that less
// than 10% of time you will be waiting c.Min and less than 10% of the time you
// will be waiting c.Max.
//
// Caveat 2 is that this assumes that the function f() takes negligible time to
// run when compared to the expected wait time. Technically memoryless events
// have the property that the times between successive event starts has the
// exponential distribution, and this function makes it so that the time between
// one event ending and the next event starting has the exponential
// distribution.
func Forever(ctx context.Context, f func(), c Config) {
	if c.Once {
		f()
	} else {
		for {
			f()
			select {
			case <-ctx.Done():
				return
			case <-time.After(c.waittime()):
				continue
			}
		}
	}
}
