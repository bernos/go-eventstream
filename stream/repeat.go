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
		defer cancelEcho()
		close(done)
		<-ack
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

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
							select {
							case <-done:
							case chInput <- event:
							}
						}
					}
				}(event)
			case <-done:
				return
			}
		}
	}()

	go func() {
		defer close(ack)
		wg.Wait()
		cancelInput()
	}()

	input.SendValue(initialValue)

	return echo, cancel
}
