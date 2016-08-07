package eventstream

import (
	"fmt"
	"math/rand"
)

func MakeRandomIntStream(seed int64) (int, Stream) {
	r := rand.New(rand.NewSource(seed))

	xs := make([]interface{}, r.Intn(10000)+1000)

	for i := 0; i < len(xs); i++ {
		xs[i] = i
	}

	return len(xs), FromSlice(xs)
}

// IntStream creates a stream that will emit ints between
// min and max inclusive. The stream will be automatically
// cancelled once all values have been sent
func IntStream(min, max int) Stream {
	xs := make([]interface{}, max-min+1)

	for i := 0; i <= max-min; i++ {
		xs[i] = min + i
	}
	return FromSlice(xs)
}

type IntMapper func(int) int

func (m IntMapper) Map(value interface{}) (interface{}, error) {
	if x, ok := value.(int); ok {
		return m(x), nil
	}
	return nil, fmt.Errorf("IntMapper expects int value from stream")
}

func Add(x int) Mapper {
	return IntMapper(func(y int) int {
		return x + y
	})
}
