package eventstream

type Reducer interface {
	Reduce(acc interface{}, value interface{}) (interface{}, error)
}

type ReducerFunc func(acc interface{}, value interface{}) (interface{}, error)

func (fn ReducerFunc) Reduce(acc interface{}, value interface{}) (interface{}, error) {
	return fn(acc, value)
}

func Scan(r Reducer) Transformer {
	return TransformerFunc(func(in Stream) Stream {
		var (
			ch  = make(chan Event)
			out = in.CreateChild(ch)
		)

		go func() {
			defer close(ch)

			var (
				acc interface{}
				err error
			)

			for e := range in.Events() {
				if e.Error() != nil {
					err = e.Error()
				} else if acc == nil {
					acc = e.Value()
				} else {
					acc, err = r.Reduce(acc, e.Value())
				}

				out.Send(acc, err)
			}
		}()

		return out
	})
}
