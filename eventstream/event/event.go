package event

type Event interface {
	Value() interface{}
	Error() error
	WithValue(interface{}) Event
	WithError(error) Event
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

func (e event) WithValue(value interface{}) Event {
	return New(value, e.Error())
}

func (e event) WithError(err error) Event {
	return New(e.Value(), err)
}
