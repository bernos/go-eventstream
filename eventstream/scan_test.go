package eventstream

import (
	"reflect"
	"testing"
)

func TestScan(t *testing.T) {
	var (
		received []interface{}
		input    = FromSlice([]interface{}{1, 1, 1, 1, 1})
		expect   = []interface{}{1, 2, 3, 4, 5}
	)

	fn := ReducerFunc(func(acc, value interface{}) (interface{}, error) {
		return acc.(int) + value.(int), nil
	})

	out := Scan(fn).Transform(input)

	for event := range out.Events() {
		received = append(received, event.Value())
	}

	if !reflect.DeepEqual(received, expect) {
		t.Errorf("Want %v, got %v", expect, received)
	}
}
