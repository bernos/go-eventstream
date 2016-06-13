package eventstream

import "testing"
import "time"

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

func TestCancelChild(t *testing.T) {
	var (
		ch     = make(chan (Event))
		parent = NewStream()
		child  = parent.CreateChild(ch)
	)

	// Wait for parent event chan to close, then close
	// child event chan
	go func() {
		defer close(ch)
		<-parent.Events()
	}()

	// Cancel child stream, which will cascade back to the
	// parent stream, closing it's event channel
	child.Cancel()

	select {
	case <-child.Events():
	case <-time.After(time.Second):
		t.Error("Timed out waiting for child stream to close")
	}
}

func TestDefaultCancel(t *testing.T) {
	s := &stream{
		events: make(chan Event),
	}

	go func() {
		s.Cancel()
	}()

	select {
	case <-s.Events():
	case <-time.After(time.Second):
		t.Error("Timed out waiting for stream to close")
	}
}

func TestMultipleCancel(t *testing.T) {
	s := NewStream()

	s.Cancel()
	s.Cancel()
}