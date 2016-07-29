package eventstream

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/bernos/go-eventstream/eventstream/event"
)

type SQSAPI interface {
	ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
}

func FromSQSQueue(svc SQSAPI, queueURL string, options ...func(*sqs.ReceiveMessageInput)) Stream {
	poller := sqsPoller(svc, queueURL, options...)

	return FromGenerator(GeneratorFunc(func(out chan event.Event, done chan struct{}) error {
		for {
			select {
			case <-done:
				return nil
			default:
				for _, e := range poller() {
					out <- e
				}
			}
		}
	}))
}

func sqsPoller(svc SQSAPI, queueURL string, options ...func(*sqs.ReceiveMessageInput)) func() []event.Event {
	return func() []event.Event {
		var (
			output []event.Event
		)

		input := &sqs.ReceiveMessageInput{
			QueueUrl: aws.String(queueURL),
		}

		for _, o := range options {
			o(input)
		}

		resp, err := svc.ReceiveMessage(input)

		if err != nil {
			output = append(output, event.New(nil, err))
		} else if resp != nil && len(resp.Messages) > 0 {
			for _, m := range resp.Messages {
				output = append(output, event.New(m, nil))
			}
		}

		return output
	}
}
