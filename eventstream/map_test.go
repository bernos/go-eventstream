package eventstream

import (
	"fmt"
	"testing"
	"time"

	"github.com/bernos/go-eventstream/eventstream/event"
)

type IntMapper func(int) int

func (m IntMapper) Map(value interface{}) (interface{}, error) {
	if x, ok := value.(int); ok {
		return m(x), nil
	}
	return nil, fmt.Errorf("IntMapper expects int value from stream")
}

func Add(x int) Mapper {
	return IntMapper(func(y int) int {
		return x + y
	})
}

func TestMap(t *testing.T) {
	in := NewStream()
	max := 5
	transform := Map(Add(2))

	go func() {
		defer in.Cancel()

		for i := 0; i < max; i++ {
			in.Send(event.New(i, nil))
		}
	}()

	out := transform.Transform(in)
	count := 0

	for event := range out.Events() {
		if x, ok := event.Value().(int); ok {
			if x != count+2 {
				t.Errorf("Want %d, got %d", count+2, x)
			}
		} else {
			t.Errorf("Want int, got %v", event.Value())
		}
		count++
	}
}

func TestMapError(t *testing.T) {
	fn := MapperFunc(func(x interface{}) (interface{}, error) {
		return 1, nil
	})

	s := NewStream()
	m := Map(fn)

	go func() {
		defer s.Cancel()
		s.Send(event.New(2, fmt.Errorf("foo")))
	}()

	out := m.Transform(s)

	select {
	case e := <-out.Events():
		if e.Error() == nil {
			t.Errorf("Expected error, but got %v", e)
		}
		if e.Value() == nil {
			t.Errorf("Expected value to be forwarded on error")
		}
	case <-time.After(time.Second):
		t.Errorf("Timeout waiting for error from stream")
	}
}

func TestMapperFunc(t *testing.T) {
	fn := MapperFunc(func(x interface{}) (interface{}, error) {
		return x.(int) + 10, nil
	})

	if y, err := fn.Map(1); err != nil {
		t.Error(err)
	} else {
		if y.(int) != 11 {
			t.Errorf("Want %d, got %d", 11, y.(int))
		}
	}

}
