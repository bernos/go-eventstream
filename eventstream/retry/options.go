package retry

import (
	"math"
	"time"
)

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

func ExponentialBackoff() func(*Retrier) {
	return func(r *Retrier) {
		r.CalculateDelay = calculateDelayBinary
	}
}

func Log(fn func(format string, v ...interface{})) func(*Retrier) {
	return func(r *Retrier) {
		r.Log = fn
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
