package eventstream

type Transformer interface {
	Transform(Stream) Stream
	Compose(Transformer) Transformer
}

type TransformerFunc func(Stream) Stream

func (t TransformerFunc) Transform(s Stream) Stream {
	return t(s)
}

func (t TransformerFunc) Compose(next Transformer) Transformer {
	return Compose(next, t)
}

func Compose(f, g Transformer) Transformer {
	return TransformerFunc(func(s Stream) Stream {
		return f.Transform(g.Transform(s))
	})
}
