package eventstream

import (
	"sync"
	"time"

	"github.com/bernos/go-eventstream/eventstream/event"
)

// Stream represents a continuous source of events
type Stream interface {
	Events() <-chan event.Event
	Send(event.Event)
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
	Reduce(Reducer, Accumulator) Stream
	Scan(Reducer, Accumulator) Stream
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

func (s *stream) Reduce(r Reducer, acc Accumulator) Stream {
	return Reduce(r, acc).Transform(s)
}

func (s *stream) Scan(r Reducer, acc Accumulator) Stream {
	return Scan(r, acc).Transform(s)
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

func (s *stream) Send(e event.Event) {
	s.events <- e
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
		wg     sync.WaitGroup
		isDone = false
		done   = make(chan struct{})
		events = make(chan event.Event)
	)

	cancel := func() {
		if !isDone {
			isDone = true
			wg.Add(1)
			go func() {
				defer func() {
					wg.Done()
					close(events)
				}()
				<-done
			}()

			close(done)
			wg.Wait()
		}
	}

	return newStream(events, nil, cancel)
}

func newStream(events chan event.Event, parent Stream, cancel CancelFunc) *stream {
	return &stream{
		events: events,
		parent: parent,
		cancel: cancel,
	}
}
