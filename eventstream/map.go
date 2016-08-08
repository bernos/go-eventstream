package eventstream

type Mapper interface {
	Map(value interface{}) (interface{}, error)
}

type MapperFunc func(interface{}) (interface{}, error)

func (fn MapperFunc) Map(x interface{}) (interface{}, error) {
	return fn(x)
}

func Map(m Mapper) Transformer {
	return Scan(ReducerFunc(func(acc Accumulator, x interface{}) (Accumulator, error) {
		return m.Map(x)
	}), nil)
}

func PMap(m Mapper, n int) Transformer {
	return TransformerFunc(func(in Stream) Stream {
		return Merge(FanOut(Map(m), n)(in)...)
	})
}

func FanOut(t Transformer, n int) func(Stream) []Stream {
	return func(in Stream) []Stream {
		streams := make([]Stream, n)

		for i := 0; i < n; i++ {
			streams[i] = t.Transform(in)
		}

		return streams
	}
}
