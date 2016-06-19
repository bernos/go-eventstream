package eventstream

import (
	"fmt"
	"reflect"
	"testing"
)

func TestScan(t *testing.T) {
	var (
		received []interface{}
		input    = FromSlice([]interface{}{1, 1, 1, 1, 1})
		expect   = []interface{}{1, 2, 3, 4, 5}
	)

	fn := ReducerFunc(func(acc Accumulator, value interface{}) (Accumulator, error) {
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

func TestScanInputError(t *testing.T) {
	var (
		input = NewStream()
	)

	fn := ReducerFunc(func(acc Accumulator, value interface{}) (Accumulator, error) {
		return acc.(int) + value.(int), nil
	})

	go func() {
		defer input.Cancel()
		input.Send(nil, fmt.Errorf("test error"))
	}()

	out := Scan(fn).Transform(input)

	for event := range out.Events() {
		if event.Error() == nil {
			t.Errorf("Expected error, got %v", event)
		}
	}
}

var result Event

func BenchmarkScan(b *testing.B) {
	var (
		input = FromSlice([]interface{}{1, 1, 1, 1, 1})
	)

	fn := ReducerFunc(func(acc Accumulator, value interface{}) (Accumulator, error) {
		return acc.(int) + value.(int), nil
	})

	for n := 0; n < b.N; n++ {
		out := Scan(fn).Transform(input)

		for event := range out.Events() {
			result = event
		}
	}
}
