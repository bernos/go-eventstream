package stream

// CloseFunc is a function that closes a stream
type CloseFunc func()

// Stream represents a continuous stream of Events
type Stream interface {
	Events() <-chan Event
	// Done() <-chan struct{}
	// Map(Mapper) (Stream, CloseFunc)
	// PMap(Mapper, int) (Stream, CloseFunc)
	Send(interface{}, error)
	SendError(error)
	SendEvent(Event)
	SendValue(interface{})
	// Transform(Transformer) (Stream, CloseFunc)
}

type stream struct {
	done   chan struct{}
	events chan Event
	close  CloseFunc
}

// func (s *stream) Done() <-chan struct{} {
// 	return s.done
// }

func (s *stream) Events() <-chan Event {
	return s.events
}

func (s *stream) Map(m Mapper) (Stream, CloseFunc) {
	return Map(m).Transform(s), s.close
}

func (s *stream) PMap(m Mapper, n int) (Stream, CloseFunc) {
	return PMap(m, n).Transform(s), s.close
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

func (s *stream) Transform(t Transformer) (Stream, CloseFunc) {
	return t.Transform(s), s.close
}
