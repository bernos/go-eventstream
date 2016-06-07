package stream

type Transformer interface {
	Transform(Stream) Stream
}

type TransformFunc func(Stream) Stream

func (t TransformFunc) Transform(s Stream) Stream {
	return t(s)
}

func Compose(f, g Transformer) Transformer {
	return TransformFunc(func(s Stream) Stream {
		return f.Transform(g.Transform(s))
	})
}
