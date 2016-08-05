package eventstream

import (
	"testing"

	"github.com/bernos/go-eventstream/eventstream/event"
)

func TestFromGenerator(t *testing.T) {
	g := func(out chan event.Event, done chan struct{}) error {
		var i int64

		for {
			select {
			case <-done:
				return nil
			case out <- event.New(i, nil):
				i++
			}
		}
	}

	max := int64(10000)
	var got int64

	s := FromGenerator(GeneratorFunc(g))

	for e := range s.Events() {
		v := e.Value().(int64)

		if v <= max {
			got = v
		} else {
			s.Cancel()
		}
	}

	if got != max {
		t.Errorf("Want %d, got %d", max, got)
	}
}
