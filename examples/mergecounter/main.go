package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bernos/go-eventstream/eventstream"
)

func main() {
	up := eventstream.
		FromPoll(randomTimer("up")).
		Const(1)

	down := eventstream.
		FromPoll(randomTimer("down")).
		Const(-1)

	total := eventstream.
		Merge(up, down).
		Scan(sum())

	for event := range total.Events() {
		fmt.Printf("%d\n", event.Value())
	}
}

func randomTimer(message string) func() (interface{}, error) {
	return func() (interface{}, error) {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		return message, nil
	}
}

func sum() eventstream.Reducer {
	return eventstream.ReducerFunc(func(acc eventstream.Accumulator, x interface{}) (eventstream.Accumulator, error) {
		return acc.(int) + x.(int), nil
	})
}
