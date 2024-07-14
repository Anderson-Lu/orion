package event

import (
	"errors"
	"sync/atomic"
	"time"
)

var (
	errPublishTimeout = errors.New("Event hub publish timeout")
)

const (
	DefaultEventHubCap    = 10
	DefaultPublishTimeout = time.Millisecond * 3
)

type EventPublishOption func(e *EventHub)
type EventConsumerFunc func(msg Event) error

// IEventHub
type IEventHub interface {
	Publish(evt ...Event) error
}

type EventHub struct {
	cs          map[uint32]chan Event
	csCap       int
	disp        uint32
	cursor      atomic.Int32
	cursorCap   uint32
	state       atomic.Int32
	pubTimeout  time.Duration
	consumers   map[uint32]EventConsumerFunc
	reconsumers map[uint32]EventConsumerFunc
}

func WithPublishTimeout(timeout time.Duration) EventPublishOption {
	return func(e *EventHub) {
		e.pubTimeout = timeout
	}
}

func WithMultiDiapatcher(num uint32) EventPublishOption {
	return func(e *EventHub) {
		e.disp = num
	}
}

func WithConsumer(dataType uint32, consumer EventConsumerFunc) EventPublishOption {
	return func(e *EventHub) {
		if e.consumers == nil {
			e.consumers = map[uint32]EventConsumerFunc{}
		}
		e.consumers[dataType] = consumer
	}
}

func WithConsumerRetry(dataType uint32, retry EventConsumerFunc) EventPublishOption {
	return func(e *EventHub) {
		if e.reconsumers == nil {
			e.reconsumers = make(map[uint32]EventConsumerFunc)
		}
		e.reconsumers[dataType] = retry
	}
}

func NewEventHub(cap int, opts ...EventPublishOption) *EventHub {

	if cap <= 0 {
		cap = DefaultEventHubCap
	}

	e := &EventHub{csCap: cap, pubTimeout: DefaultPublishTimeout}
	for _, option := range opts {
		option(e)
	}

	e.init()
	return e
}

func (e *EventHub) init() {
	if ok := e.state.CompareAndSwap(0, 1); !ok {
		return
	}
	if e.disp == 0 {
		e.disp = 1
	}
	e.cs = make(map[uint32]chan Event)
	for i := 0; i < int(e.disp); i++ {
		e.cs[uint32(i)] = make(chan Event, e.csCap)
	}
	for i := 0; i < int(e.disp); i++ {
		e.consume(uint32(i))
	}
	e.cursor.Store(0)
}

func (e *EventHub) dispatcher() chan Event {
	n := e.cursor.Add(1)
	if n > int32(e.cursorCap) {
		e.tryResetCursor()
	}
	return e.cs[uint32(n%int32(e.disp))]
}

func (e *EventHub) tryResetCursor() {
	if loaded := e.cursor.Load(); loaded < int32(e.cursorCap) {
		return
	}
	e.cursor.Store(0)
}

func (e *EventHub) consume(i uint32) {
	go func() {
		if int(i) > len(e.cs)-1 {
			return
		}
		for msg := range e.cs[i] {
			consumer, ok := e.consumers[uint32(msg.Type())]
			if ok && consumer != nil {
				if err := consumer(msg); err != nil && e.reconsumers[uint32(msg.Type())] != nil {
					e.reconsumers[uint32(msg.Type())](msg)
				}
			}
		}
	}()
}

func (e *EventHub) Publish(evts ...Event) error {
	for _, et := range evts {
		select {
		case e.dispatcher() <- et:
			return nil
		case <-time.After(e.pubTimeout):
			return errPublishTimeout
		}

	}
	return nil
}
