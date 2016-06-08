package stream

func FromChannel(ch chan Event) (Stream, CancelFunc) {
	ack := make(chan struct{})

	s := &stream{
		done:   make(chan struct{}),
		events: ch,
	}

	go func() {
		defer func() {
			close(ack)
			close(s.events)
		}()
		<-s.done
	}()

	return s, doneSignalerWithCallback(s, func() {
		<-ack
	})
}

func NewStream() (Stream, CancelFunc) {
	var (
		ack       = make(chan struct{})
		s, cancel = newStream(ack)
	)

	go func() {
		defer func() {
			close(ack)
			close(s.events)
		}()
		<-s.done
	}()

	return s, cancel
}

func Once(value interface{}) Stream {
	s, cancel := NewStream()

	go func() {
		defer cancel()
		s.SendValue(value)
	}()

	return s
}

func FromPoll(fn func() (interface{}, error)) (Stream, CancelFunc) {
	var (
		ack        = make(chan struct{})
		in, cancel = newStream(ack)
	)

	go func() {
		defer func() {
			close(ack)
			close(in.events)
		}()

		for {
			select {
			case <-in.done:
				return
			case in.events <- NewEvent(fn()):
			}
		}
	}()

	return in, cancel
}

func FromSlice(xs []interface{}) (Stream, CancelFunc) {
	var (
		ack        = make(chan struct{})
		in, cancel = newStream(ack)
	)

	go func() {
		defer func() {
			close(ack)
			close(in.events)
		}()

		for i := range xs {
			select {
			case <-in.done:
				return
			case in.events <- NewEvent(xs[i], nil):
			}
		}
	}()

	return in, cancel
}

func RepeatSlice(xs []interface{}) (Stream, CancelFunc) {
	var (
		ack        = make(chan struct{})
		in, cancel = newStream(ack)
	)

	go func() {
		defer func() {
			close(ack)
			close(in.events)
		}()

		i := 0

		for {
			select {
			case <-in.done:
				return
			case in.events <- NewEvent(xs[i], nil):
				i = (i + 1) % len(xs)
			}
		}
	}()

	return in, cancel
}

func newStream(ack chan struct{}) (*stream, CancelFunc) {
	in := &stream{
		done:   make(chan struct{}),
		events: make(chan Event),
	}

	return in, doneSignalerWithCallback(in, func() {
		<-ack
	})
}

func doneSignalerWithCallback(s *stream, fn func()) CancelFunc {
	isDone := false

	return func() {
		if !isDone {
			isDone = true
			close(s.done)
			fn()
		}
	}
}
