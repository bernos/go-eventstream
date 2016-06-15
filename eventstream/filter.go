package eventstream

func Filter(fn Predicate) Transformer {
	return TransformerFunc(func(in Stream) Stream {
		var (
			ch  = make(chan Event)
			out = in.CreateChild(ch)
		)

		go func() {
			defer close(ch)

			for e := range in.Events() {
				if fn(e) {
					out.Send(e.Value(), e.Error())
				}
			}
		}()

		return out
	})
}
