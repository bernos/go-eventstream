package eventstream

import "github.com/bernos/go-eventstream/eventstream/event"

type Predicate func(event.Event) bool

func Take(n int) Transformer {
	return TakeWhile(func(x int) Predicate {
		return func(e event.Event) bool {
			if x < n {
				x++
				return true
			}
			return false
		}
	}(0))
}

func TakeUntil(fn Predicate) Transformer {
	return TakeWhile(func(e event.Event) bool {
		return !fn(e)
	})
}

func TakeWhile(fn Predicate) Transformer {
	return TransformerFunc(func(in Stream) Stream {
		var (
			ch  = make(chan event.Event)
			out = in.CreateChild(ch)
		)

		go func() {
			defer close(ch)

			done := false

			for e := range in.Events() {
				if !done {
					if fn(e) {
						out.Send(e.Value(), e.Error())
					} else {
						done = true
						in.Cancel()
					}
				}
			}
		}()

		return out
	})
}
