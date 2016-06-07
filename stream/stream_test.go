package stream

import (
	"log"
	"testing"
)

func TestFromSlice(t *testing.T) {
	values := []interface{}{1, 2, 3, 4, 5}
	s, _ := FromSlice(values)
	var out []int

	for event := range s.Events() {
		if x, ok := event.Value().(int); ok {
			out = append(out, x)
		} else {
			t.Errorf("Expected value to be int")
		}
	}

	if len(out) != len(values) {
		t.Errorf("Expected len to be %d, but got %d", len(values), len(out))
	}

	for i := range out {
		if out[i] != values[i] {
			t.Errorf("Expected value %d, got %d", values[i], out[i])
		}
	}
}

func TestFromSliceWithClose(t *testing.T) {
	values := []interface{}{1, 2, 3, 4, 5}
	s, cls := FromSlice(values)
	var out []int

	for event := range s.Events() {
		if x, ok := event.Value().(int); ok {
			if x == 4 {
				cls()
			} else {
				out = append(out, x)
			}
		} else {
			t.Errorf("Expected value to be int")
		}
	}

	if len(out) != 3 {
		t.Errorf("Expected len to be %d, but got %d", 3, len(out))
	}

	for i := range out {
		log.Println(out[i])
		if out[i] != values[i] {
			t.Errorf("Expected value %d, got %d", values[i], out[i])
		}
	}
}

func TestRepeatSlice(t *testing.T) {
	values := []interface{}{1, 2, 3, 4, 5}
	s, cls := RepeatSlice(values)
	var out []int

	for event := range s.Events() {
		if x, ok := event.Value().(int); ok {
			out = append(out, x)
			if len(out) == len(values)*2 {
				cls()
			}
		} else {
			t.Errorf("Expected value to be int")
		}
	}

	if len(out) != len(values)*2 {
		t.Errorf("Expected len to be %d, but got %d", len(values)*2, len(out))
	}

	for i := range out {
		if out[i] != values[i%len(values)] {
			t.Errorf("Expected value %d, got %d", values[i%len(values)], out[i])
		}
	}
}
