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

func TestFromSlice(t *testing.T) {
	s := FromSlice([]interface{}{1, 2, 3, 4, 5})

	count := 0

	for e := range s.Events() {
		count++

		if e.Value().(int) != count {
			t.Errorf("Want %d, got %d", count, e.Value().(int))
		}
	}

	if count != 5 {
		t.Errorf("Want 5 values, got %d", count)
	}
}
