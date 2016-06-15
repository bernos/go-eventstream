package eventstream

import (
	"reflect"
	"testing"
)

func TestFlatMap(t *testing.T) {
	var (
		received []interface{}
		input    = FromSlice([]interface{}{1, 2, 3})
		expect   = []interface{}{1, 1, 2, 2, 3, 3}
	)

	fn := FlatMapperFunc(func(x interface{}) ([]interface{}, error) {
		return []interface{}{x, x}, nil
	})

	out := FlatMap(fn).Transform(input)

	for event := range out.Events() {
		received = append(received, event.Value())
	}

	if !reflect.DeepEqual(received, expect) {
		t.Errorf("Want %v, got %v", expect, received)
	}
}
