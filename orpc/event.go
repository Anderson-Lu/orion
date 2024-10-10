package orpc

const (
	EventTypeLog int16 = 1
)

type Event struct {
	typ  int16
	data interface{}
}

func (e Event) Type() int16 {
	return e.typ
}

func (e Event) Data() interface{} {
	return e.data
}
