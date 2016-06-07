package stream

func FromChannel(ch chan Event) (Stream, CloseFunc) {
	ack := make(chan struct{})

	s := &stream{
		done:   make(chan struct{}),
		events: ch,
	}

	s.close = doneSignalerWithCallback(s, func() {
		<-ack
	})

	go func() {
		defer func() {
			close(ack)
			close(s.events)
		}()
		<-s.done
	}()

	return s, s.close
}

func NewStream() (Stream, CloseFunc) {
	var (
		ack = make(chan struct{})
		s   = newStream(ack)
	)

	go func() {
		defer func() {
			close(ack)
			close(s.events)
		}()
		<-s.done
	}()

	return s, s.close
}

func Once(value interface{}) Stream {
	s, cls := NewStream()

	go func() {
		defer cls()
		s.SendValue(value)
	}()

	return s
}

func FromPoll(fn func() (interface{}, error)) (Stream, CloseFunc) {
	var (
		ack = make(chan struct{})
		in  = newStream(ack)
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

	return in, in.close
}

func FromSlice(xs []interface{}) (Stream, CloseFunc) {
	var (
		ack = make(chan struct{})
		in  = newStream(ack)
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

	return in, in.close
}

func RepeatSlice(xs []interface{}) (Stream, CloseFunc) {
	var (
		ack = make(chan struct{})
		in  = newStream(ack)
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

	return in, in.close
}

func newStream(ack chan struct{}) *stream {
	in := &stream{
		done:   make(chan struct{}),
		events: make(chan Event),
	}

	in.close = doneSignalerWithCallback(in, func() {
		<-ack
	})

	return in
}

func doneSignalerWithCallback(s *stream, fn func()) CloseFunc {
	isDone := false

	return func() {
		if !isDone {
			isDone = true
			close(s.done)
			fn()
		}
	}
}
