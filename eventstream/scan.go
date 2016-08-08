package eventstream

import "github.com/bernos/go-eventstream/eventstream/event"

type Accumulator interface{}

type Reducer interface {
	Reduce(Accumulator, interface{}) (Accumulator, error)
}

type ReducerFunc func(Accumulator, interface{}) (Accumulator, error)

func (fn ReducerFunc) Reduce(acc Accumulator, value interface{}) (Accumulator, error) {
	return fn(acc, value)
}

func Scan(r Reducer, acc Accumulator) Transformer {
	return TransformerFunc(func(in Stream) Stream {
		var (
			ch  = make(chan event.Event)
			out = in.CreateChild(ch)
		)

		go func() {
			defer close(ch)

			var (
				// acc interface{}
				err error
			)

			for e := range in.Events() {
				if e.Error() != nil {
					err = e.Error()
					// } else if acc == nil {
					// 	acc = e.Value()
				} else {
					acc, err = r.Reduce(acc, e.Value())
				}

				out.Send(event.New(acc, err))
			}
		}()

		return out
	})
}
