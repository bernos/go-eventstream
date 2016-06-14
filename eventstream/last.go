package eventstream

func Last() Transformer {
	return TransformerFunc(func(in Stream) Stream {
		var (
			ch  = make(chan Event)
			out = in.CreateChild(ch)
		)

		go func() {
			defer close(ch)

			var e Event
			for e = range in.Events() {
			}

			out.Send(e.Value(), e.Error())
		}()

		return out
	})
}
