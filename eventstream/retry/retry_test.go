package retry

import (
	"fmt"
	"testing"
	"testing/quick"
	"time"
)

func SkipTestCalculateDelay(t *testing.T) {
	f := func(max int) bool {
		last := time.Duration(0)

		for i := 0; i < max; i++ {
			d := calculateDelay(i, time.Millisecond, time.Minute)
			if d < 0 || d < last || d > time.Minute {
				return false
			}
			last = d
		}
		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestMaxRetries(t *testing.T) {
	count := 0

	fn := func() (interface{}, error) {
		count++
		return nil, fmt.Errorf("foo")
	}

	r := Retry(fn, MaxRetries(2))
	r()

	if count != 3 {
		t.Errorf("Want %d, got %d", 3, count)
	}
}
