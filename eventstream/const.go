package eventstream

func Const(x interface{}) Transformer {
	return Map(MapperFunc(func(v interface{}) (interface{}, error) {
		return x, nil
	}))
}
