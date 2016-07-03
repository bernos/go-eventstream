package eventstream

import (
	"reflect"
	"testing"

	"github.com/bernos/go-eventstream/eventstream/event"
)

func TestFilter(t *testing.T) {
	var (
		received []interface{}
		input    = FromSlice([]interface{}{1, 2, 3, 4, 5})
		expect   = []interface{}{1, 2, 3}
	)

	fn := func(e event.Event) bool {
		return e.Value().(int) < 4
	}

	out := Filter(fn).Transform(input)

	for event := range out.Events() {
		received = append(received, event.Value())
	}

	if !reflect.DeepEqual(received, expect) {
		t.Errorf("Want %v, got %v", expect, received)
	}
}
