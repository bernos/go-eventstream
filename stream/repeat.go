package stream

import "sync"

func Repeat(initialValue interface{}, transformer Transformer) (Stream, CancelFunc) {
	var (
		done               = make(chan struct{})
		ack                = make(chan struct{})
		chInput            = make(chan Event)
		chEcho             = make(chan Event)
		input, cancelInput = FromChannel(chInput)
		echo, cancelEcho   = FromChannel(chEcho)
		wg                 sync.WaitGroup
	)

	output := transformer.Transform(input)

	cancel := func() {
		defer func() {
			cancelInput()
			cancelEcho()
		}()
		close(done)
		<-ack
	}

	go func() {
		defer close(ack)

		for {
			select {
			case event := <-output.Events():
				wg.Add(1)

				go func(event Event) {
					defer wg.Done()

					select {
					case <-done:
						return
					case chEcho <- event:
						if event.Error() == nil {
							input.SendEvent(event)
						}
					}
				}(event)
			case <-done:
				wg.Wait()
				return
			}
		}
	}()

	input.SendValue(initialValue)

	return echo, cancel
}
