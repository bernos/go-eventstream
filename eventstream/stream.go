package eventstream

import (
	"time"

	"github.com/bernos/go-eventstream/eventstream/event"
)

// Stream represents a continuous source of events
type Stream interface {
	Events() <-chan event.Event
	Send(interface{}, error)
	Cancel()
	CreateChild(chan event.Event) Stream

	Const(interface{}) Stream
	Filter(Predicate) Stream
	FlatMap(FlatMapper) Stream
	PFlatMap(FlatMapper, int) Stream
	Id() Stream
	Last() Stream
	Map(Mapper) Stream
	PMap(Mapper, int) Stream
	Reduce(Reducer) Stream
	Scan(Reducer) Stream
	Take(int) Stream
	TakeUntil(Predicate) Stream
	TakeWhile(Predicate) Stream
	Throttle(time.Duration) Stream
}

// CancelFunc cancels a stream
type CancelFunc func()

// Base stream implementation
type stream struct {
	events chan event.Event
	parent Stream
	cancel CancelFunc
}

func (s *stream) Const(x interface{}) Stream {
	return Const(x).Transform(s)
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

func (s *stream) Id() Stream {
	return Id().Transform(s)
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

func (s *stream) Events() <-chan event.Event {
	return s.events
}

func (s *stream) Send(value interface{}, err error) {
	s.events <- event.New(value, err)
}

func (s *stream) Throttle(d time.Duration) Stream {
	return Throttle(d).Transform(s)
}

func (s *stream) Cancel() {
	if s.cancel != nil {
		s.cancel()
	} else if s.parent != nil {
		s.parent.Cancel()
	} else {
		close(s.events)
	}
}

// FromChannel creates a child stream. Calling Cancel() on the child
// stream will walk the stream heirarchy and close the done channel
// on the root stream
func (s *stream) CreateChild(in chan event.Event) Stream {
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
		events = make(chan event.Event)
	)

	cancel := func() {
		if !isDone {
			isDone = true

			go func() {
				defer func() {
					close(ack)
					close(events)
				}()
				<-done
			}()

			close(done)
			<-ack
		}
	}

	return newStream(events, nil, cancel)
}
