package event

// Event
type Event interface {
	Type() int16
	Data() interface{}
}
