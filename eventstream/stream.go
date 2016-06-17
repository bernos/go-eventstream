package eventstream

import "github.com/bernos/go-frp/eventstream/retry"

// Stream represents a continuous source of events
type Stream interface {
	Events() <-chan Event
	Send(interface{}, error)
	Cancel()
	CreateChild(chan Event) Stream

	Filter(Predicate) Stream
	FlatMap(FlatMapper) Stream
	PFlatMap(FlatMapper, int) Stream
	Last() Stream
	Map(Mapper) Stream
	PMap(Mapper, int) Stream
	Reduce(Reducer) Stream
	Scan(Reducer) Stream
	Take(int) Stream
	TakeUntil(Predicate) Stream
	TakeWhile(Predicate) Stream
}

// CancelFunc cancels a stream
type CancelFunc func()

// Base stream implementation
type stream struct {
	events chan Event
	parent Stream
	cancel CancelFunc
}

func (s *stream) Filter(fn Predicate) Stream {
	return Filter(fn).Transform(s)
}

func (s *stream) FlatMap(m FlatMapper) Stream {
	return FlatMap(m).Transform(s)
}

func (s *stream) PFlatMap(m FlatMapper, n int) Stream {
	return PFlatMap(m, n).Transform(s)
}

func (s *stream) Last() Stream {
	return Last().Transform(s)
}

func (s *stream) Map(fn Mapper) Stream {
	return Map(fn).Transform(s)
}

func (s *stream) PMap(m Mapper, n int) Stream {
	return PMap(m, n).Transform(s)
}

func (s *stream) Reduce(r Reducer) Stream {
	return Reduce(r).Transform(s)
}

func (s *stream) Scan(r Reducer) Stream {
	return Scan(r).Transform(s)
}

func (s *stream) Take(n int) Stream {
	return Take(n).Transform(s)
}

func (s *stream) TakeUntil(fn Predicate) Stream {
	return TakeUntil(fn).Transform(s)
}

func (s *stream) TakeWhile(fn Predicate) Stream {
	return TakeWhile(fn).Transform(s)
}

func (s *stream) Events() <-chan Event {
	return s.events
}

func (s *stream) Send(value interface{}, err error) {
	s.events <- NewEvent(value, err)
}

func (s *stream) Cancel() {
	if s.parent != nil {
		s.parent.Cancel()
	} else if s.cancel != nil {
		s.cancel()
	} else {
		close(s.events)
	}
}

// FromChannel creates a child stream. Calling Cancel() on the child
// stream will walk the stream heirarchy and close the done channel
// on the root stream
func (s *stream) CreateChild(in chan Event) Stream {
	return &stream{
		parent: s,
		events: in,
	}
}

// NewStream creates a new stream
func NewStream() Stream {
	var (
		isDone = false
		ack    = make(chan struct{})
		done   = make(chan struct{})
	)

	cancel := func(s *stream) CancelFunc {
		return func() {
			if !isDone {
				isDone = true

				go func() {
					defer func() {
						close(ack)
						close(s.events)
					}()
					<-done
				}()

				close(done)
				<-ack
			}
		}
	}

	return newStream(make(chan Event), nil, cancel)
}

// Once creates a stream that will send a single value, then cancel automatically
func Once(value interface{}) Stream {
	s := NewStream()

	go func() {
		defer s.Cancel()
		s.Send(value, nil)
	}()

	return s
}

// FromPoll creates a stream by calling fn for each value to send
func FromPoll(fn func() (interface{}, error)) Stream {
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

		for {
			select {
			case in.events <- NewEvent(fn()):
			case <-done:
				return
			}
		}
	}()

	return in
}

func FromPollWithRetries(fn func() (interface{}, error), options ...func(*retry.Retrier)) Stream {
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
	fn = retry.Retry(fn, options...)

	go func() {
		defer func() {
			close(in.events)
			close(ack)
		}()

		for {
			select {
			case in.events <- NewEvent(fn()):
			case <-done:
				return
			}
		}
	}()

	return in
}

// FromSlice creates a stream that sends each item in xs and then cancels automatically
func FromSlice(xs []interface{}) Stream {
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

		for i := range xs {
			select {
			case <-done:
				return
			case in.events <- NewEvent(xs[i], nil):
			}
		}
	}()

	return in
}

func newStream(events chan Event, parent Stream, cancel func(*stream) CancelFunc) *stream {
	s := &stream{
		events: events,
		parent: parent,
	}

	s.cancel = cancel(s)

	return s
}
