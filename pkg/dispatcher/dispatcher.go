package dispatcher

import (
	"errors"
	"time"
)

type IDispatcherHandler interface {
	Invoke(something interface{})
	Stop()
}

type DispatcherOption func(d *Dispatcher)

func WithPendingCap(cap int) DispatcherOption {
	return func(d *Dispatcher) {
		if cap <= 0 {
			return
		}
		d.pendingCap = cap
	}
}

func WithConcurrent(con int) DispatcherOption {
	return func(d *Dispatcher) {
		if con <= 0 {
			return
		}
		d.concurrentCnt = con
	}
}

func WithHandlerPoolFunc(c func() IDispatcherHandler) DispatcherOption {
	return func(d *Dispatcher) {
		d.handlerCreator = c
	}
}

type Dispatcher struct {
	pendingCap     int
	concurrentCnt  int
	chPending      chan interface{}
	chHandlers     []IDispatcherHandler
	handlerCreator func() IDispatcherHandler
	chDone         chan struct{}
	readyStack     *RandomStack
	name           string
	D              chan time.Time
}

func NewDispatcher(name string, opts ...DispatcherOption) (*Dispatcher, error) {
	d := &Dispatcher{pendingCap: 10, concurrentCnt: 1}
	for _, v := range opts {
		v(d)
	}
	d.chPending = make(chan interface{}, d.pendingCap*10)
	if d.handlerCreator == nil {
		return nil, errors.New("none handler creator specified")
	}
	for i := 0; i < d.concurrentCnt; i++ {
		d.chHandlers = append(d.chHandlers, d.handlerCreator())
	}
	d.chDone = make(chan struct{}, 1)
	d.readyStack = NewRandomStack(d.concurrentCnt)
	d.name = name
	d.D = make(chan time.Time, 1)
	return d, nil
}

func (d *Dispatcher) Commit(something interface{}) {
	d.chPending <- something
}

func (d *Dispatcher) Stop() {
	d.chDone <- struct{}{}
	for _, v := range d.chHandlers {
		v.Stop()
	}
}

func (d *Dispatcher) Start() {
	go func() {
		t := time.NewTicker(time.Second * 1)
		defer t.Stop()
		tick := time.Now().Unix()
		for {
			select {
			case <-d.chDone:
				goto DONE
			case data := <-d.chPending:
				idx := d.readyStack.Get()
				tick = time.Now().Unix()
				go func(id int, anything interface{}) {
					defer d.readyStack.Back(id)
					d.chHandlers[id%d.concurrentCnt].Invoke(anything)
				}(idx, data)
			case <-t.C:
				if len(d.chPending) == 0 && time.Now().Unix()-tick > 60 {
					goto DONE
				}
			}
		}
	DONE:
		d.D <- time.Now()
	}()
}

func NewRandomStack(cap int) *RandomStack {
	r := &RandomStack{}
	r.data = make(chan int, cap*2)
	for i := 0; i < cap; i++ {
		r.data <- i
	}
	r.cap = cap
	return r
}

type RandomStack struct {
	data chan int
	cap  int
}

func (r *RandomStack) Get() int {
	c := <-r.data
	return c
}

func (r *RandomStack) Back(k int) {
	r.data <- k
}
