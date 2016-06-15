package eventstream

func Reduce(r Reducer) Transformer {
	return Compose(Last(), Scan(r))
}
