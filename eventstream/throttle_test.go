package eventstream

import (
	"testing"
	"time"
)

func TestThrottle(t *testing.T) {
	values := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9}
	throttle := time.Millisecond * 100
	expect := time.Duration(len(values)) * throttle

	out := Throttle(throttle).Transform(FromSlice(values))

	start := time.Now()
	count := 0

	for _ = range out.Events() {
		count++
	}

	diff := time.Since(start) - expect

	if diff > (time.Millisecond * 5) {
		t.Errorf("want %s, got %s", expect, diff)
	}
}
