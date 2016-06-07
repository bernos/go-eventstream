package stream

import (
	"sync"
)

type Mapper interface {
	Map(value interface{}) (interface{}, error)
}

// MapperFunc is a func that implements Mapper
type MapperFunc func(interface{}) (interface{}, error)

// Map satisfies the Mapper interface
func (fn MapperFunc) Map(x interface{}) (interface{}, error) {
	return fn(x)
}

func Map(m Mapper) Transformer {
	return PMap(m, 1)
}

// PMap is a parallel implementation of Map. It produces a Pipeline that maps
// all values from its input stream to its output stream via n concurrent instances
// of the Mapper m
func PMap(m Mapper, n int) Transformer {
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
						out.SendEvent(NewEvent(m.Map(event.Value())))
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
