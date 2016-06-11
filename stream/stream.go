package stream

// CancelFunc is a function that closes a stream
type CancelFunc func()

// Stream represents a continuous stream of Events
type Stream interface {
	InputStream
	OutputStream
}

type InputStream interface {
	Send(interface{}, error)
	SendError(error)
	SendEvent(Event)
	SendValue(interface{})
}

type OutputStream interface {
	Events() <-chan Event
	// FlatMap(FlatMapper) Stream
	// PFlatMap(FlatMapper, int) Stream
	// Map(Mapper) Stream
	// PMap(Mapper, int) Stream
	Transform(Transformer) Stream
}

type stream struct {
	done   chan struct{}
	events chan Event
}

func (s *stream) Events() <-chan Event {
	return s.events
}

func (s *stream) FlatMap(m FlatMapper) Stream {
	return FlatMap(m).Transform(s)
}

func (s *stream) PFlatMap(m FlatMapper, n int) Stream {
	return PFlatMap(m, n).Transform(s)
}

func (s *stream) Map(m Mapper) Stream {
	return Map(m).Transform(s)
}

func (s *stream) PMap(m Mapper, n int) Stream {
	return PMap(m, n).Transform(s)
}

func (s *stream) Send(value interface{}, err error) {
	s.SendEvent(NewEvent(value, err))
}

func (s *stream) SendError(err error) {
	s.SendEvent(NewEvent(nil, err))
}

func (s *stream) SendEvent(e Event) {
	s.events <- e
}

func (s *stream) SendValue(value interface{}) {
	s.SendEvent(NewEvent(value, nil))
}

func (s *stream) Transform(t Transformer) Stream {
	return t.Transform(s)
}
