package circuit_break

import (
	"sync"
	"sync/atomic"
)

type Bucket struct {
	cursor  int64
	succCnt int64
	failCnt int64

	avgCost float64

	mu sync.Mutex
}

func (b *Bucket) Succ(cost int64) {
	atomic.AddInt64(&b.succCnt, 1)
	b.update(cost)
}

func (b *Bucket) Fail(cost int64) {
	atomic.AddInt64(&b.failCnt, 1)
	b.update(cost)
}

func (b *Bucket) Stat() (cnt int64, succRate, avgCost float64) {
	return b.succCnt + b.failCnt, float64(b.succCnt) / float64(b.succCnt+b.failCnt), b.avgCost
}

func (b *Bucket) update(cost int64) {

	b.mu.Lock()
	defer b.mu.Unlock()

	cnt := b.succCnt + b.failCnt
	if cnt == 0 {
		return
	}
	if cnt == 1 {
		b.avgCost = float64(cnt)
		return
	}

	b.avgCost = (float64(cnt-1)/float64(cnt))*b.avgCost + (float64(cost) / float64(cnt))
}

func (b *Bucket) Reset(cursor int64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	atomic.StoreInt64(&b.cursor, cursor)
	atomic.StoreInt64(&b.failCnt, 0)
	atomic.StoreInt64(&b.succCnt, 0)
	b.avgCost = 0
}
