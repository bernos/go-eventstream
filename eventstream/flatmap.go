package eventstream

import "sync"
import "github.com/bernos/go-eventstream/eventstream/event"

type FlatMapper interface {
	FlatMap(interface{}) ([]interface{}, error)
}

type FlatMapperFunc func(interface{}) ([]interface{}, error)

func (fn FlatMapperFunc) FlatMap(x interface{}) ([]interface{}, error) {
	return fn(x)
}

func FlatMap(m FlatMapper) Transformer {
	return PFlatMap(m, 1)
}

func PFlatMap(m FlatMapper, n int) Transformer {
	return TransformerFunc(func(in Stream) Stream {
		var (
			wg  sync.WaitGroup
			ch  = make(chan event.Event)
			out = in.CreateChild(ch)
		)

		wg.Add(n)

		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()

				for e := range in.Events() {
					if e.Error() != nil {
						out.Send(e)
					} else {
						if xs, err := m.FlatMap(e.Value()); err == nil {
							for x := range xs {
								out.Send(event.New(xs[x], nil))
							}
						} else {
							out.Send(e.WithError(err))
						}
					}
				}
			}()
		}

		go func() {
			defer close(ch)
			wg.Wait()
		}()

		return out
	})
}
