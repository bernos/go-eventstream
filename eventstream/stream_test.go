package eventstream

import (
	"testing"
)

func TestOnce(t *testing.T) {
	s := Once("hello")

	count := 0

	for e := range s.Events() {
		count++
		if e.Value() != "hello" {
			t.Errorf("Expected 'hello', got '%s'", e.Value())
		}
	}

	if count != 1 {
		t.Errorf("Want 1 event, got %d", count)
	}

}
