package stream

import "sync"

type Mapper interface {
	Map(value interface{}) (interface{}, error)
}

// MapperFunc is a func that implements Mapper
type MapperFunc func(interface{}) (interface{}, error)

// Map satisfies the Mapper interface
func (fn MapperFunc) Map(x interface{}) (interface{}, error) {
	return fn(x)
}

func Map(m Mapper) Transformer {
	return PMap(m, 1)
}

// PMap is a parallel implementation of Map. It produces a Pipeline that maps
// all values from its input stream to its output stream via n concurrent instances
// of the Mapper m
func PMap(m Mapper, n int) Transformer {
	return TransformFunc(func(in Stream) Stream {
		var (
			wg       sync.WaitGroup
			out, cls = NewStream()
		)

		wg.Add(n)

		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()

				for event := range in.Events() {
					if event.Error() != nil {
						out.SendEvent(event)
					} else {
						out.SendEvent(NewEvent(m.Map(event.Value())))
					}
				}
			}()
		}

		go func() {
			defer cls()
			wg.Wait()
		}()

		return out
	})
}

func Loop(start Stream, cancel CancelFunc, shouldContinue func() bool) Transformer {
	return TransformFunc(func(in Stream) Stream {
		var (
			wg       sync.WaitGroup
			out, cls = NewStream()
			feedback = make(chan Event)
			done     = make(chan struct{})
		)

		// Read from in and send to out. Also send to feedback chan
		// if there was no error
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				cls()
				close(done)
			}()

			for e := range in.Events() {

				if e.Error() == nil {
					wg.Add(1)
					go func(e Event) {
						defer wg.Done()
						select {
						case <-done:
							return
						case feedback <- e:
						}
					}(e)
				}
			}
		}()

		// Read from feedback ch, and send back to
		// start stream
		go func() {
			defer cancel()

			for e := range feedback {
				if !shouldContinue() {
					return
				}

				select {
				case <-done:
					return
				default:
					start.SendEvent(e)
					out.SendEvent(e)
				}
			}
		}()

		go func() {
			defer close(feedback)
			wg.Wait()
		}()

		// Read from feedback ch, and send back to
		// start stream
		// go func() {

		// 	for {
		// 		time.Sleep(time.Nanosecond)

		// 		select {
		// 		case <-done:
		// 			return
		// 		case e := <-feedback:
		// 		}
		// 	}
		// }()

		// Read from feedback ch, and send back to
		// start stream
		// go func() {

		// 	for {
		// 		select {
		// 		case <-done:
		// 			return
		// 		case e := <-feedback:
		// 			go func(e Event) {
		// 				select {
		// 				case <-done:
		// 					return
		// 				default:
		// 					start.SendEvent(e)
		// 				}
		// 			}(e)
		// 		}
		// 	}
		// }()

		return out
	})
}
