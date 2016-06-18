package retry

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	count := 0
	messageCount := 0

	fn := func() (interface{}, error) {
		if count == 10 {
			return "done", nil
		}
		count++
		return nil, fmt.Errorf("Test retry")
	}

	r := Retry(fn, Log(func(format string, v ...interface{}) {
		messageCount++
	}))

	v, err := r()

	if v != "done" {
		t.Errorf("Want %s, got %s", "done", v)
	}

	if err != nil {
		t.Errorf("Unexpected error. %s", err.Error())
	}

	if messageCount < count {
		t.Errorf("Want %d messages, got %d", count, messageCount)
	}
}
