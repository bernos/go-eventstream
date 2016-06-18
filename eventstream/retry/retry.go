package retry

import (
	"fmt"
	"math"
	"time"
)

const (
	DefaultMaxRetries = 10
	DefaultBaseDelay  = time.Millisecond
	DefaultMaxDelay   = time.Minute
	Infinity          = -1
)

type Retrier struct {
	MaxRetries  int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	ShouldRetry func(error) bool
}

func BaseDelay(d time.Duration) func(*Retrier) {
	return func(r *Retrier) {
		r.BaseDelay = d
	}
}

func Forever() func(*Retrier) {
	return func(r *Retrier) {
		r.MaxRetries = Infinity
	}
}

func MaxRetries(n int) func(*Retrier) {
	return func(r *Retrier) {
		r.MaxRetries = n
	}
}

func MaxDelay(d time.Duration) func(*Retrier) {
	return func(r *Retrier) {
		r.MaxDelay = d
	}
}

func While(fn func(error) bool) func(*Retrier) {
	return func(r *Retrier) {
		r.ShouldRetry = fn
	}
}

func Retry(fn func() (interface{}, error), options ...func(*Retrier)) func() (interface{}, error) {
	r := Retrier{
		BaseDelay:   DefaultBaseDelay,
		MaxDelay:    DefaultMaxDelay,
		MaxRetries:  DefaultMaxRetries,
		ShouldRetry: func(err error) bool { return true },
	}

	for _, o := range options {
		o(&r)
	}

	return func() (interface{}, error) {
		var count int

		for {
			value, err := fn()

			if err == nil {
				return value, err
			}

			if !r.ShouldRetry(err) {
				return nil, fmt.Errorf("Retrier aborted due to user supplied ShouldRetry func. Cause: %s", err.Error())
			}

			if count == r.MaxRetries {
				return nil, fmt.Errorf("Retrier exceeded max retry count of %d. Cause: %s", r.MaxRetries, err.Error())
			}

			time.Sleep(calculateDelayBinary(uint(count), r.BaseDelay, r.MaxDelay))
			count++
		}
	}
}

func binaryRaise(exponent uint) int64 {
	// Clamp exponent to avoid 64 bit overflow
	if exponent > 62 {
		exponent = 62
	}

	return 1 << exponent
}

func calculateDelayBinary(iteration uint, baseDelay, maxDelay time.Duration) time.Duration {
	m := (binaryRaise(iteration) - 1) >> 1

	// If multiplier is greater than maxDelay, then we we don't need to
	// calculate
	if m > int64(maxDelay) {
		return maxDelay
	}

	d := (time.Duration(m) * baseDelay) + baseDelay

	if d < maxDelay && d > 0 {
		return d
	}

	return maxDelay
}

func calculateDelay(iteration uint, baseDelay, maxDelay time.Duration) time.Duration {
	m := ((math.Pow(2, float64(iteration))) - 1) / 2

	// Check for wrapping, or multiplier greater than max duration in
	// nanosecs
	if int64(m) < 0 || int64(m) > int64(maxDelay) {
		return maxDelay
	}

	d := (time.Duration(m) * baseDelay) + baseDelay

	if d < maxDelay && d > 0 {
		return d
	}

	return maxDelay
}
