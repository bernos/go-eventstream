package main

import (
	"fmt"

	"github.com/bernos/go-eventstream/eventstream"
)

func main() {
	out := eventstream.
		FromValues(1, 2, 3, 4, 5).
		PMap(square(), 9)

	for event := range out.Events() {
		fmt.Printf("%d\n", event.Value())
	}
}

func square() eventstream.Mapper {
	return eventstream.MapperFunc(func(x interface{}) (interface{}, error) {
		return x.(int) * x.(int), nil
	})
}
