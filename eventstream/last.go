package eventstream

import "github.com/bernos/go-eventstream/eventstream/event"

func Last() Transformer {
	return TransformerFunc(func(in Stream) Stream {
		var (
			ch  = make(chan event.Event)
			out = in.CreateChild(ch)
		)

		go func() {
			defer close(ch)

			var e event.Event
			for e = range in.Events() {
			}

			out.Send(e)
		}()

		return out
	})
}
