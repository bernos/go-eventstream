package eventstream

import (
	"fmt"
	"sync"

	"github.com/bernos/go-eventstream/eventstream/event"
)

type Mapper interface {
	Map(value interface{}) (interface{}, error)
}

type MapperFunc func(interface{}) (interface{}, error)

func (fn MapperFunc) Map(x interface{}) (interface{}, error) {
	return fn(x)
}

func Map(m Mapper) Transformer {
	// return PMap(m, 1)
	return Scan(ReducerFunc(func(acc Accumulator, x interface{}) (Accumulator, error) {
		fmt.Print("REDUCE!")
		return m.Map(x)
	}), nil)
}

func PMap(m Mapper, n int) Transformer {
	return TransformerFunc(func(in Stream) Stream {

		streams := make([]Stream, n)

		for i := 0; i < n; i++ {
			streams[i] = Map(m).Transform(in)
		}

		return Merge(streams...)
	})
}

func OldPMap(m Mapper, n int) Transformer {
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
						out.Send(event.New(m.Map(e.Value())))
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
