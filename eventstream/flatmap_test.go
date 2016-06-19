package eventstream

import (
	"fmt"
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

func TestFlatMapInputError(t *testing.T) {
	var (
		input = NewStream()
	)

	go func() {
		defer input.Cancel()
		input.Send(nil, fmt.Errorf("test error"))
	}()

	fn := FlatMapperFunc(func(x interface{}) ([]interface{}, error) {
		return []interface{}{x, x}, nil
	})

	out := FlatMap(fn).Transform(input)

	for event := range out.Events() {
		if event.Error() == nil {
			t.Errorf("Expected error, but got %v", event)
		}
	}
}

func TestFlatMapError(t *testing.T) {
	var (
		input = NewStream()
	)

	go func() {
		defer input.Cancel()
		input.Send("foo", nil)
	}()

	fn := FlatMapperFunc(func(x interface{}) ([]interface{}, error) {
		return nil, fmt.Errorf("test error")
	})

	out := FlatMap(fn).Transform(input)

	for event := range out.Events() {
		if event.Error() == nil {
			t.Errorf("Expected error, but got %v", event)
		}
	}
}
