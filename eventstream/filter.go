package eventstream

import "github.com/bernos/go-eventstream/eventstream/event"

func Filter(fn Predicate) Transformer {
	return TransformerFunc(func(in Stream) Stream {
		var (
			ch  = make(chan event.Event)
			out = in.CreateChild(ch)
		)

		go func() {
			defer close(ch)

			for e := range in.Events() {
				if fn(e) {
					out.Send(e)
				}
			}
		}()

		return out
	})
}
