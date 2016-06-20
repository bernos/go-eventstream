package eventstream

type Generator interface {
	Generate(chan Event, chan struct{}) error
}

type GeneratorFunc func(chan Event, chan struct{}) error

func (fn GeneratorFunc) Generate(out chan Event, done chan struct{}) error {
	return fn(out, done)
}

func FromGenerator(g Generator) Stream {
	var (
		isDone = false
		ack    = make(chan struct{})
		done   = make(chan struct{})
	)

	cancel := func(s *stream) CancelFunc {
		return func() {
			if !isDone {
				isDone = true
				close(done)
				<-ack
			}
		}
	}

	in := newStream(make(chan Event), nil, cancel)

	go func() {
		defer func() {
			close(in.events)
			close(ack)
		}()

		err := g.Generate(in.events, done)

		if err != nil {
			in.Send(nil, err)
		}
	}()

	return in
}

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

func Once(value interface{}) Stream {
	return FromGenerator(GeneratorFunc(func(out chan Event, done chan struct{}) error {
		select {
		case <-done:
		case out <- NewEvent(value, nil):
		}
		return nil
	}))
}

func FromValues(xs ...interface{}) Stream {
	return FromSlice(xs)
}

func newStream(events chan Event, parent Stream, cancel func(*stream) CancelFunc) *stream {
	s := &stream{
		events: events,
		parent: parent,
	}

	s.cancel = cancel(s)

	return s
}
