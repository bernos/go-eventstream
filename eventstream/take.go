package eventstream

type Predicate func(Event) bool

func Take(n int) Transformer {
	return TransformerFunc(func(in Stream) Stream {
		var (
			ch  = make(chan Event)
			out = in.CreateChild(ch)
		)

		go func() {
			defer close(ch)

			count := 0

			for e := range in.Events() {
				out.Send(e.Value(), e.Error())
				count++
				if count == n {
					in.Cancel()
				}
			}
		}()

		return out
	})
}

func TakeWhile(fn Predicate) Transformer {
	return TransformerFunc(func(in Stream) Stream {
		var (
			ch  = make(chan Event)
			out = in.CreateChild(ch)
		)

		go func() {
			defer close(ch)
			done := false

			for e := range in.Events() {
				if !done && fn(e) {
					out.Send(e.Value(), e.Error())
				} else if !done {
					done = true
					in.Cancel()
				}
			}
		}()

		return out
	})

}

func TakeUntil(fn Predicate) Transformer {
	return TransformerFunc(func(in Stream) Stream {
		var (
			ch  = make(chan Event)
			out = in.CreateChild(ch)
		)

		go func() {
			defer close(ch)
			done := false

			for e := range in.Events() {
				if !done && fn(e) {
					done = true
					in.Cancel()
				} else if !done {
					out.Send(e.Value(), e.Error())
				}
			}
		}()

		return out
	})

}
