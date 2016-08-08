package eventstream

func Reduce(r Reducer, acc Accumulator) Transformer {
	return Compose(Last(), Scan(r, acc))
}
