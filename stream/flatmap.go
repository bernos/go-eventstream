package stream

import "sync"

type FlatMapper interface {
	FlatMap(interface{}) ([]interface{}, error)
}

type FlatMapperFunc func(interface{}) ([]interface{}, error)

func (fn FlatMapperFunc) FlatMap(value interface{}) ([]interface{}, error) {
	return fn(value)
}

func FlatMap(m FlatMapper) Transformer {
	return PFlatMap(m, 1)
}

func PFlatMap(m FlatMapper, n int) Transformer {
	return TransformFunc(func(in Stream) Stream {
		var (
			wg       sync.WaitGroup
			out, cls = NewStream()
		)

		wg.Add(n)

		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()

				for event := range in.Events() {
					if event.Error() != nil {
						out.SendEvent(event)
					} else {
						values, err := m.FlatMap(event.Value())

						if err == nil {
							for v := range values {
								out.SendValue(values[v])
							}
						} else {
							out.SendError(err)
						}
					}
				}
			}()
		}

		go func() {
			defer cls()
			wg.Wait()
		}()

		return out
	})
}
