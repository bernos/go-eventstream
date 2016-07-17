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

			for e := range in.Events() {
				out.Send(e)
				<-ticker.C
			}
		}()

		return out
	})
}
