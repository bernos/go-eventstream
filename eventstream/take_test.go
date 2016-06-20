package eventstream

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func TestTake(t *testing.T) {
	var (
		received []interface{}
		input    = FromSlice([]interface{}{1, 2, 3, 4, 5})
		expect   = []interface{}{1, 2, 3}
		out      = Take(3).Transform(input)
	)

	for event := range out.Events() {
		received = append(received, event.Value())
	}

	if !reflect.DeepEqual(received, expect) {
		t.Errorf("Want %v, got %v", expect, received)
	}
}

func TestTakeWhile(t *testing.T) {
	var (
		received []interface{}
		input    = FromSlice([]interface{}{1, 2, 3, 4, 5})
		expect   = []interface{}{1, 2, 3}
	)

	fn := func(e Event) bool {
		return e.Value().(int) < 4
	}

	out := TakeWhile(fn).Transform(input)

	for event := range out.Events() {
		received = append(received, event.Value())
	}

	if !reflect.DeepEqual(received, expect) {
		t.Errorf("Want %v, got %v", expect, received)
	}
}

func TestTakeUntil(t *testing.T) {
	var (
		received []interface{}
		input    = FromSlice([]interface{}{1, 2, 3, 4, 5})
		expect   = []interface{}{1, 2, 3}
	)

	fn := func(e Event) bool {
		return e.Value().(int) > 3
	}

	out := TakeUntil(fn).Transform(input)

	for event := range out.Events() {
		received = append(received, event.Value())
	}

	if !reflect.DeepEqual(received, expect) {
		t.Errorf("Want %v, got %v", expect, received)
	}
}

func makeRandomIntStream(seed int64) (int, Stream) {
	r := rand.New(rand.NewSource(seed))

	xs := make([]interface{}, r.Intn(10000))

	for i := 0; i < len(xs); i++ {
		xs[i] = i
	}

	return len(xs), FromSlice(xs)
}

func TestTakeUntilOverflow(t *testing.T) {
	f := func(x int64) bool {
		n, s := makeRandomIntStream(x)
		max := 0
		want := int(n / 2)

		// Ensure take func is preceeded by some other transformers
		// force overflow
		out := s.Id().Id().TakeUntil(func(e Event) bool {
			return e.Value().(int) > want
		})

		for event := range out.Events() {
			max = event.Value().(int)
		}

		return max == want
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestTakeWhileOverflow(t *testing.T) {
	f := func(x int64) bool {
		n, s := makeRandomIntStream(x)
		max := 0
		want := int(n / 2)

		// Ensure take func is preceeded by some other transformers
		// force overflow
		out := s.Id().Id().TakeWhile(func(e Event) bool {
			return e.Value().(int) <= want
		})

		for event := range out.Events() {
			max = event.Value().(int)
		}

		return max == want
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestTakeOverflow(t *testing.T) {
	f := func(x int64) bool {
		n, s := makeRandomIntStream(x)
		max := 0
		want := int(n / 2)

		// Ensure take func is preceeded by some other transformers
		// force overflow
		out := s.Id().Id().Take(want)

		for event := range out.Events() {
			max = event.Value().(int)
		}

		return max == want-1
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
