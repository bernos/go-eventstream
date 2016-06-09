package stream

type Transformer interface {
	Transform(Stream) Stream
	Repeat(interface{}) (Stream, CancelFunc)
}

type TransformFunc func(Stream) Stream

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
