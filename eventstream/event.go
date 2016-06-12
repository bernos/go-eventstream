package eventstream

type Event interface {
	Value() interface{}
	Error() error
}

type event struct {
	value interface{}
	err   error
}

func NewEvent(value interface{}, err error) Event {
	return event{value, err}
}

func (e event) Value() interface{} {
	return e.value
}

func (e event) Error() error {
	return e.err
}
