package eventstream

import (
	"reflect"
	"testing"
)

func TestLast(t *testing.T) {
	var (
		received int
		count    int
		input    = FromSlice([]interface{}{1, 2, 3, 4, 5})
		expect   = 5
		out      = Last().Transform(input)
	)

	for event := range out.Events() {
		count++
		received = event.Value().(int)
	}

	if count != 1 {
		t.Errorf("Expected 1 event, got %d", count)
	}

	if !reflect.DeepEqual(received, expect) {
		t.Errorf("Want %v, got %v", expect, received)
	}
}
