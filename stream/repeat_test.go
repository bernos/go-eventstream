package stream

import (
	"math"
	"testing"
	"testing/quick"
)

func TestRepeat(t *testing.T) {

	assert := func(n int) bool {
		numRepeats := int(math.Abs(float64(n % 1000)))
		count := 0

		transformer := Map(Add(2))

		out, cancel := Repeat(1, transformer)

		for _ = range out.Events() {
			count++
			if count == numRepeats {
				cancel()
			}
		}

		if count != numRepeats {
			t.Errorf("Expected %d, got %d", numRepeats, count)
		}

		return true
	}

	config := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(assert, config); err != nil {
		t.Error(err)
	}
}
