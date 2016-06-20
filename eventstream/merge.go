package eventstream

import "sync"

// Merge merges events from Streams a and b and passes them to its
// output stream. Cancelling the output stream will cancel both a
// and b. The output stream will close automatically if both a and
// b are closed independently
func Merge(in ...Stream) Stream {
	var (
		wg     sync.WaitGroup
		isDone = false
		ch     = make(chan Event)
	)

	cancel := func(s *stream) CancelFunc {
		return func() {
			if !isDone {
				isDone = true

				for _, s := range in {
					s.Cancel()
				}
			}
		}
	}

	out := newStream(ch, nil, cancel)

	for _, s := range in {
		wg.Add(1)

		go func(s Stream) {
			defer wg.Done()

			for event := range s.Events() {
				out.Send(event.Value(), event.Error())
			}
		}(s)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	return out
}
