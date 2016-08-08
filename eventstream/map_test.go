package eventstream

import (
	"fmt"
	"testing"
	"time"

	"github.com/bernos/go-eventstream/eventstream/event"
)

func TestMap(t *testing.T) {
	in := IntStream(0, 5)
	transform := Map(Add(2))

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

	in.Cancel()
}

func TestPMap(t *testing.T) {
	n := 5

	mapper := MapperFunc(func(x interface{}) (interface{}, error) {
		time.Sleep(time.Second)
		return x, nil
	})

	in := IntStream(1, n)
	transform := PMap(mapper, n)

	out := transform.Transform(in)
	start := time.Now().UTC()

	var result []interface{}

	for event := range out.Events() {
		for i := range result {
			if result[i] == event.Value() {
				t.Errorf("Already saw %v", event.Value())
			}
		}
		result = append(result, event.Value())
	}

	duration := time.Since(start)

	if duration > time.Second*time.Duration(n) {
		t.Errorf("Expected < %d, got %s", n, duration)
	}

	in.Cancel()
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
