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

func TestOnce(t *testing.T) {
	s := Once("hello")

	count := 0

	for e := range s.Events() {
		count++
		if e.Value() != "hello" {
			t.Errorf("Expected 'hello', got '%s'", e.Value())
		}
	}

	if count != 1 {
		t.Errorf("Want 1 event, got %d", count)
	}

}

func TestFromPoll(t *testing.T) {
	x := 0

	poll := func() (interface{}, error) {
		y := x
		x++
		return y, nil
	}

	s := FromPoll(poll)

	max := 3
	count := 0

	for e := range s.Events() {

		if e.Value().(int) != count {
			t.Errorf("Want %d, got %d", count, e.Value().(int))
		}

		count++

		if count == max {
			s.Cancel()
		}
	}

	if count != max {
		t.Errorf("Want 5 values, got %d", count)
	}
}

func TestFromSlice(t *testing.T) {
	s := FromSlice([]interface{}{1, 2, 3, 4, 5})

	count := 0

	for e := range s.Events() {
		count++

		if e.Value().(int) != count {
			t.Errorf("Want %d, got %d", count, e.Value().(int))
		}
	}

	if count != 5 {
		t.Errorf("Want 5 values, got %d", count)
	}
}

func TestFromSliceCancel(t *testing.T) {
	s := FromSlice([]interface{}{1, 2, 3, 4, 5})

	max := 3
	count := 0

	for e := range s.Events() {
		count++

		if e.Value().(int) != count {
			t.Errorf("Want %d, got %d", count, e.Value().(int))
		}

		if count == max {
			s.Cancel()
		}
	}

	if count != max {
		t.Errorf("Want 5 values, got %d", count)
	}

}
