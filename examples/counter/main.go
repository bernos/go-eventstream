package main

import (
	"fmt"
	"math/rand"

	"github.com/bernos/go-eventstream/eventstream"
)

var (
	directions   = []string{"up", "down"}
	directionMap = map[string]int{
		"up":   1,
		"down": -1,
	}
)

func main() {
	out := eventstream.
		FromPoll(randomDirection).
		Map(toInt()).
		Scan(sum())

	for event := range out.Events() {
		fmt.Printf("%d\n", event.Value())

		if event.Value().(int) > 500 {
			out.Cancel()
		}
	}
}

func randomDirection() (interface{}, error) {
	return directions[rand.Intn(len(directions))], nil
}

func toInt() eventstream.Mapper {
	return eventstream.MapperFunc(func(x interface{}) (interface{}, error) {
		return directionMap[x.(string)], nil
	})
}

func sum() eventstream.Reducer {
	return eventstream.ReducerFunc(func(acc eventstream.Accumulator, x interface{}) (eventstream.Accumulator, error) {
		return acc.(int) + x.(int), nil
	})
}
