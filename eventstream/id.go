package eventstream

func Id() Transformer {
	return Map(MapperFunc(func(x interface{}) (interface{}, error) {
		return x, nil
	}))
}
