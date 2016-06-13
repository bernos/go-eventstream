package eventstream

// Stream represents a continuous source of events
type Stream interface {
	Events() <-chan Event
	Send(interface{}, error)
	Cancel()
	CreateChild(chan Event) Stream
}

// CancelFunc cancels a stream
type CancelFunc func()

// Base stream implementation
type stream struct {
	events chan Event
	parent Stream
	cancel CancelFunc
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
