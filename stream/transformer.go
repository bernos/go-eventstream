package stream

type Transformer interface {
	Transform(Stream) Stream
	Compose(Transformer) Transformer
	Map(Mapper) Transformer
	PMap(Mapper, int) Transformer
	FlatMap(FlatMapper) Transformer
	PFlatMap(FlatMapper, int) Transformer
	Repeat(interface{}) (Stream, CancelFunc)
}

type TransformFunc func(Stream) Stream

func (t TransformFunc) Compose(next Transformer) Transformer {
	return Compose(next, t)
}

func (t TransformFunc) Map(m Mapper) Transformer {
	return t.Compose(Map(m))
}

func (t TransformFunc) PMap(m Mapper, n int) Transformer {
	return t.Compose(PMap(m, n))
}

func (t TransformFunc) FlatMap(m FlatMapper) Transformer {
	return t.Compose(FlatMap(m))
}

func (t TransformFunc) PFlatMap(m FlatMapper, n int) Transformer {
	return t.Compose(PFlatMap(m, n))
}

func (t TransformFunc) Transform(s Stream) Stream {
	return t(s)
}

func (t TransformFunc) Repeat(initialValue interface{}) (Stream, CancelFunc) {
	return Repeat(initialValue, t)
}

func Compose(f, g Transformer) Transformer {
	return TransformFunc(func(s Stream) Stream {
		return f.Transform(g.Transform(s))
	})
}
