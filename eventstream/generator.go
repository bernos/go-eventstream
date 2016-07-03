package eventstream

type Generator interface {
	Generate(out chan Event, done chan struct{}) error
}

type GeneratorFunc func(chan Event, chan struct{}) error

func (fn GeneratorFunc) Generate(out chan Event, done chan struct{}) error {
	return fn(out, done)
}

// FromGenerator creates a new Stream from a Generator
func FromGenerator(g Generator) Stream {
	var (
		isDone = false
		ack    = make(chan struct{})
		done   = make(chan struct{})
	)

	cancel := func() {
		if !isDone {
			isDone = true
			close(done)
			<-ack
		}
	}

	out := newStream(make(chan Event), nil, cancel)

	go func() {
		defer func() {
			close(out.events)
			close(ack)
		}()

		err := g.Generate(out.events, done)

		if err != nil {
			out.Send(nil, err)
		}
	}()

	return out
}

// FromPoll creates a new Stream by polling the provided function
func FromPoll(fn func() (interface{}, error)) Stream {
	return FromGenerator(GeneratorFunc(func(out chan Event, done chan struct{}) error {
		for {
			select {
			case out <- NewEvent(fn()):
			case <-done:
				return nil
			}
		}
	}))
}

// FromSlice creates a new Stream that will emit an event for each item in a slice.
// Once all the events have been sent the Stream will automatically close
func FromSlice(xs []interface{}) Stream {
	return FromGenerator(GeneratorFunc(func(out chan Event, done chan struct{}) error {
		for i := range xs {
			select {
			case <-done:
				return nil
			case out <- NewEvent(xs[i], nil):
			}
		}
		return nil
	}))
}

// Once creates a new Stream that will emit a single value and then automatically close
func Once(value interface{}) Stream {
	return FromGenerator(GeneratorFunc(func(out chan Event, done chan struct{}) error {
		select {
		case <-done:
		case out <- NewEvent(value, nil):
		}
		return nil
	}))
}

// FromValues creates a new Stream that will emit an event for each of the arguments
// and then close
func FromValues(xs ...interface{}) Stream {
	return FromSlice(xs)
}

func newStream(events chan Event, parent Stream, cancel CancelFunc) *stream {
	return &stream{
		events: events,
		parent: parent,
		cancel: cancel,
	}
}
