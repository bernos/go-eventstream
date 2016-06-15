package eventstream

import (
	"reflect"
	"testing"
)

func TestTake(t *testing.T) {
	var (
		received []interface{}
		input    = FromSlice([]interface{}{1, 2, 3, 4, 5})
		expect   = []interface{}{1, 2, 3}
		out      = Take(3).Transform(input)
	)

	for event := range out.Events() {
		received = append(received, event.Value())
	}

	if !reflect.DeepEqual(received, expect) {
		t.Errorf("Want %v, got %v", expect, received)
	}
}

func TestTakeWhile(t *testing.T) {
	var (
		received []interface{}
		input    = FromSlice([]interface{}{1, 2, 3, 4, 5})
		expect   = []interface{}{1, 2, 3}
	)

	fn := func(e Event) bool {
		return e.Value().(int) < 4
	}

	out := TakeWhile(fn).Transform(input)

	for event := range out.Events() {
		received = append(received, event.Value())
	}

	if !reflect.DeepEqual(received, expect) {
		t.Errorf("Want %v, got %v", expect, received)
	}
}

func TestTakeUntil(t *testing.T) {
	var (
		received []interface{}
		input    = FromSlice([]interface{}{1, 2, 3, 4, 5})
		expect   = []interface{}{1, 2, 3}
	)

	fn := func(e Event) bool {
		return e.Value().(int) > 3
	}

	out := TakeUntil(fn).Transform(input)

	for event := range out.Events() {
		received = append(received, event.Value())
	}

	if !reflect.DeepEqual(received, expect) {
		t.Errorf("Want %v, got %v", expect, received)
	}
}
