package eventstream

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
	in := NewStream()
	max := 5
	transform := Map(Add(2))

	go func() {
		defer in.Cancel()

		for i := 0; i < max; i++ {
			in.Send(i, nil)
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
