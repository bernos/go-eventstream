package eventstream

import (
	"reflect"
	"testing"
)

func TestReduce(t *testing.T) {
	var (
		received int
		count    int
		input    = FromSlice([]interface{}{1, 2, 3, 4, 5})
		expect   = 15
	)

	fn := ReducerFunc(func(acc Accumulator, value interface{}) (Accumulator, error) {
		return acc.(int) + value.(int), nil
	})

	out := Reduce(fn).Transform(input)

	for event := range out.Events() {
		count++
		received = event.Value().(int)
	}

	if count != 1 {
		t.Errorf("Expected 1 value, go %d", count)
	}

	if !reflect.DeepEqual(received, expect) {
		t.Errorf("Want %v, got %v", expect, received)
	}
}
