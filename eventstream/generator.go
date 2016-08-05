package eventstream

import (
	"sync"

	"github.com/bernos/go-eventstream/eventstream/event"
)

type Generator interface {
	Generate(out chan event.Event, done chan struct{}) error
}

type GeneratorFunc func(chan event.Event, chan struct{}) error

func (fn GeneratorFunc) Generate(out chan event.Event, done chan struct{}) error {
	return fn(out, done)
}

// FromGenerator creates a new Stream from a Generator
func FromGenerator(g Generator) Stream {
	var (
		wg     sync.WaitGroup
		isDone = false
		done   = make(chan struct{})
	)

	cancel := func() {
		if !isDone {
			isDone = true
			close(done)
			wg.Wait()
		}
	}

	out := newStream(make(chan event.Event), nil, cancel)

	wg.Add(1)

	go func() {
		defer wg.Done()

		err := g.Generate(out.events, done)

		if err != nil {
			out.Send(event.New(nil, err))
		}

		close(out.events)
	}()

	return out
}

// FromPoll creates a new Stream by polling the provided function
func FromPoll(fn func() (interface{}, error)) Stream {
	return FromGenerator(GeneratorFunc(func(out chan event.Event, done chan struct{}) error {
		for {
			select {
			case out <- event.New(fn()):
			case <-done:
				return nil
			}
		}
	}))
}

// FromSlice creates a new Stream that will emit an event for each item in a slice.
// Once all the events have been sent the Stream will automatically close
func FromSlice(xs []interface{}) Stream {
	return FromGenerator(GeneratorFunc(func(out chan event.Event, done chan struct{}) error {
		for i := range xs {
			select {
			case <-done:
				return nil
			case out <- event.FromValue(xs[i]):
			}
		}
		return nil
	}))
}

// Once creates a new Stream that will emit a single value and then automatically close
func Once(value interface{}) Stream {
	return FromGenerator(GeneratorFunc(func(out chan event.Event, done chan struct{}) error {
		select {
		case <-done:
		case out <- event.FromValue(value):
		}
		return nil
	}))
}

// FromValues creates a new Stream that will emit an event for each of the arguments
// and then close
func FromValues(xs ...interface{}) Stream {
	return FromSlice(xs)
}
