package stream

import (
	"fmt"
	"testing"
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
	in, cls := NewStream()
	max := 5
	transformer := Map(Add(2))

	go func() {
		defer cls()

		for i := 0; i < max; i++ {
			in.SendEvent(NewEvent(i, nil))
		}
	}()

	out := transformer.Transform(in)
	count := 0

	for event := range out.Events() {
		if x, ok := event.Value().(int); ok {
			if x != count+2 {
				t.Errorf("Expected %d, but got %d", count+2, x)
			}
		} else {
			t.Errorf("Expected int but got %v", event.Value())
		}
		count++
	}
}

func TestMapViaTransform(t *testing.T) {
	in, cls := NewStream()
	max := 5
	out := Map(Add(2)).Transform(in)

	go func() {
		defer cls()

		for i := 0; i < max; i++ {
			in.SendEvent(NewEvent(i, nil))
		}
	}()

	count := 0

	for event := range out.Events() {
		if x, ok := event.Value().(int); ok {
			if x != count+2 {
				t.Errorf("Expected %d, but got %d", count+2, x)
			}
		} else {
			t.Errorf("Expected int but got %v", event.Value())
		}
		count++
	}
}

func TestMapViaStream(t *testing.T) {
	in := Once(1)
	out := Map(Add(2)).Transform(in)

	got := <-out.Events()

	if x, ok := got.Value().(int); ok {
		if x != 1+2 {
			t.Errorf("Want %d, got %d", 1+2, x)
		}
	} else {
		t.Errorf("Expected int but got %v", got.Value())
	}
}
