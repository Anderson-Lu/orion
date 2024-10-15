package circuit_break

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"
)

type WindowConfig struct {

	// window size, eg: 1000ms
	Duration int

	// bucket size = Duration / Buckets
	Buckets int
}

func (w *WindowConfig) Revise() {

	if w.Duration <= 0 {
		w.Duration = 1000
	}

	if w.Buckets <= 0 {
		w.Buckets = 10
	}

	if w.Duration < w.Buckets {
		w.Buckets = w.Duration
	}
}

type Window struct {
	wc     *WindowConfig
	bkts   []*Bucket
	now    int64
	cursor int
	sg     singleflight.Group
	stu    CircuitBreakStatus
	chk    IRuleChecker
}

func NewWindow(wc *WindowConfig, chk IRuleChecker) *Window {

	wc.Revise()

	wd := &Window{}
	wd.wc = wc
	wd.chk = chk
	wd.initBucket()
	wd.daemon()

	return wd
}

func (w *Window) initBucket() {
	w.bkts = make([]*Bucket, 0)
	for i := 0; i < w.wc.Buckets; i++ {
		w.bkts = append(w.bkts, new(Bucket))
	}
}

func (w *Window) metrics() [][]float64 {
	r := [][]float64{}
	for _, v := range w.bkts {
		r = append(r, []float64{float64(v.succCnt), float64(v.failCnt), v.avgCost})
	}
	return r
}

func (w *Window) updateStatus() {
	if w.chk == nil {
		return
	}
	w.stu = w.chk.Turn(context.Background(), w.metrics())
}

func (w *Window) daemon() {
	f := func() {
		bucketDuration := int64(w.wc.Duration / w.wc.Buckets)
		ticker := time.NewTicker(time.Millisecond * time.Duration(bucketDuration))
		for range ticker.C {
			w.initCursor(time.Now().UnixMilli())
			w.updateStatus()
		}
	}
	go f()
}

func (w *Window) initCursor(now int64) {
	w.sg.Do("init_cursor", func() (interface{}, error) {

		if now == w.now {
			return nil, nil
		}
		bucketDuration := int64(w.wc.Duration / w.wc.Buckets)

		if w.now == 0 {
			w.now = now
		}

		gapBuckets := (now - w.now) / bucketDuration
		if gapBuckets > int64(w.wc.Buckets) {
			gapBuckets = int64(w.wc.Buckets)
		}
		for i := 0; i < int(gapBuckets); i++ {
			idx := (w.cursor + (i + 1)) % len(w.bkts)
			w.bkts[idx].Reset(w.now + int64(i+1)*bucketDuration)
		}
		w.now += gapBuckets * bucketDuration
		w.cursor = (w.cursor + int(gapBuckets)) % len(w.bkts)
		return nil, nil
	})
}

func (w *Window) getBucket(now int64) *Bucket {
	w.initCursor(now)
	return w.bkts[w.cursor]
}

func (w *Window) Succ(now, cost int64) {
	w.getBucket(now).Succ(cost)
}

func (w *Window) Fail(now, cost int64) {
	w.getBucket(now).Fail(cost)
}
