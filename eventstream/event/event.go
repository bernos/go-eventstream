package event

type Event interface {
	Value() interface{}
	Error() error
}

type event struct {
	value interface{}
	err   error
}

func New(value interface{}, err error) Event {
	return event{value, err}
}

func FromError(err error) Event {
	return event{nil, err}
}

func FromValue(value interface{}) Event {
	return event{value, nil}
}

func (e event) Value() interface{} {
	return e.value
}

func (e event) Error() error {
	return e.err
}
