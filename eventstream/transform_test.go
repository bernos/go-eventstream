package eventstream

import (
	"testing"
	"time"

	"github.com/bernos/go-eventstream/eventstream/event"
)

func TestCompose(t *testing.T) {
	makeTransformer := func(char string) Transformer {
		return TransformerFunc(func(in Stream) Stream {
			var (
				ch  = make(chan event.Event)
				out = in.CreateChild(ch)
			)

			go func() {
				defer close(ch)

				for e := range in.Events() {
					out.Send(e.Value().(string)+char, nil)
				}
			}()

			return out
		})
	}

	s := NewStream()

	go func() {
		defer s.Cancel()
		s.Send("a", nil)
	}()

	out := makeTransformer("b").
		Compose(makeTransformer("c")).
		Transform(s)

	select {
	case e := <-out.Events():
		if e.Value().(string) != "abc" {
			t.Errorf("Want 'abc', got '%v'", e.Value())
		}
	case <-time.After(time.Second):
		t.Error("Timed out waiting for transformer")
	}
}
