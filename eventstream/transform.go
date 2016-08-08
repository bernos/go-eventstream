package eventstream

import "time"

type Transformer interface {
	Transform(Stream) Stream
	Compose(Transformer) Transformer
	Const(interface{}) Transformer
	Filter(Predicate) Transformer
	FlatMap(FlatMapper) Transformer
	PFlatMap(FlatMapper, int) Transformer
	Id() Transformer
	Last() Transformer
	Map(Mapper) Transformer
	PMap(Mapper, int) Transformer
	Reduce(Reducer, Accumulator) Transformer
	Scan(Reducer, Accumulator) Transformer
	Take(int) Transformer
	TakeWhile(Predicate) Transformer
	TakeUntil(Predicate) Transformer
	Throttle(time.Duration) Transformer
}

type TransformerFunc func(Stream) Stream

func (t TransformerFunc) Transform(s Stream) Stream {
	return t(s)
}

func (t TransformerFunc) Const(x interface{}) Transformer {
	return t.Compose(Const(x))
}

func (t TransformerFunc) Filter(fn Predicate) Transformer {
	return t.Compose(Filter(fn))
}

func (t TransformerFunc) FlatMap(m FlatMapper) Transformer {
	return t.Compose(FlatMap(m))
}

func (t TransformerFunc) Id() Transformer {
	return t.Compose(Id())
}

func (t TransformerFunc) Last() Transformer {
	return t.Compose(Last())
}

func (t TransformerFunc) PFlatMap(m FlatMapper, n int) Transformer {
	return t.Compose(PFlatMap(m, n))
}

func (t TransformerFunc) Map(m Mapper) Transformer {
	return t.Compose(Map(m))
}

func (t TransformerFunc) PMap(m Mapper, n int) Transformer {
	return t.Compose(PMap(m, n))
}

func (t TransformerFunc) Reduce(r Reducer, acc Accumulator) Transformer {
	return t.Compose(Reduce(r, acc))
}

func (t TransformerFunc) Scan(r Reducer, acc Accumulator) Transformer {
	return t.Compose(Scan(r, acc))
}

func (t TransformerFunc) Take(n int) Transformer {
	return t.Compose(Take(n))
}

func (t TransformerFunc) TakeUntil(fn Predicate) Transformer {
	return t.Compose(TakeUntil(fn))
}

func (t TransformerFunc) TakeWhile(fn Predicate) Transformer {
	return t.Compose(TakeWhile(fn))
}

func (t TransformerFunc) Throttle(d time.Duration) Transformer {
	return t.Compose(Throttle(d))
}

func (t TransformerFunc) Compose(next Transformer) Transformer {
	return Compose(next, t)
}

func Compose(f, g Transformer) Transformer {
	return TransformerFunc(func(s Stream) Stream {
		return f.Transform(g.Transform(s))
	})
}
