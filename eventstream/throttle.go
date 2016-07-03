package eventstream

import (
	"time"

	"github.com/bernos/go-eventstream/eventstream/event"
)

func Throttle(d time.Duration) Transformer {
	return TransformerFunc(func(in Stream) Stream {
		var (
			ch  = make(chan event.Event)
			out = in.CreateChild(ch)
		)

		go func() {
			defer close(ch)
			ticker := time.NewTicker(d)
			defer ticker.Stop()

			for event := range in.Events() {
				out.Send(event.Value(), event.Error())
				<-ticker.C
			}
		}()

		return out
	})
}
