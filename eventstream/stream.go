package eventstream

type Stream interface {
	Events() <-chan Event
	Send(interface{}, error)
	Cancel()
	FromChannel(chan Event) Stream
}

type stream struct {
	events chan Event
}

func (s *stream) Events() <-chan Event {
	return s.events
}

func (s *stream) Send(value interface{}, err error) {
	s.events <- NewEvent(value, err)
}

func (s *stream) Cancel() {
	close(s.events)
}

func (s *stream) FromChannel(in chan Event) Stream {
	child := &childStream{
		parent: s,
	}
	child.stream.events = in

	return child
}

type sourceStream struct {
	stream
	ack    chan struct{}
	done   chan struct{}
	isDone bool
}

func (s *sourceStream) Cancel() {
	if !s.isDone {
		s.isDone = true
		close(s.done)
		<-s.ack
	}
}

func NewStream() Stream {
	s := &sourceStream{
		ack:    make(chan struct{}),
		done:   make(chan struct{}),
		isDone: false,
	}

	go func() {
		defer func() {
			close(s.ack)
			close(s.events)
		}()
		<-s.done
	}()

	s.stream.events = make(chan Event)

	return s
}

type childStream struct {
	stream
	parent Stream
}

func (s *childStream) Cancel() {
	s.parent.Cancel()
}
