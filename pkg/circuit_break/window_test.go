package circuit_break

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkWindow$ github.com/Anderson-Lu/orion/pkg/circuit_break
func BenchmarkWindow(b *testing.B) {
	wd := NewWindow(&WindowConfig{Duration: 50, Buckets: 5}, nil)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if idx%2 == 0 {
				wd.Succ(time.Now().UnixMilli(), int64(idx))
			} else {
				wd.Fail(time.Now().UnixMilli(), int64(idx))
			}
		}(i)
	}
	wg.Wait()
}

func TestWindow(t *testing.T) {
	wd := NewWindow(&WindowConfig{Duration: 50, Buckets: 5}, nil)
	assert.Equal(t, 5, len(wd.bkts))

	// [x, 0, 0, 0, 0]
	wd.Succ(1, 1)
	assert.Equal(t, 0, int(wd.cursor))
	assert.Equal(t, 0, int(wd.bkts[0].cursor))
	assert.Equal(t, 1, int(wd.bkts[0].avgCost))
	assert.Equal(t, 1, int(wd.bkts[0].succCnt))
	assert.Equal(t, 0, int(wd.bkts[0].failCnt))

	wd.Fail(2, 2)
	assert.Equal(t, 0, int(wd.cursor))
	assert.Equal(t, 0, int(wd.bkts[0].cursor))
	assert.Equal(t, float64(1.5), (wd.bkts[0].avgCost))
	assert.Equal(t, 1, int(wd.bkts[0].succCnt))
	assert.Equal(t, 1, int(wd.bkts[0].failCnt))

	// // [0, x, 0, 0, 0]
	wd.Succ(11, 1)
	assert.Equal(t, 1, wd.cursor)
	assert.Equal(t, 10, int(wd.bkts[1].cursor))
	assert.Equal(t, float64(1.0), (wd.bkts[1].avgCost))
	assert.Equal(t, 1, int(wd.bkts[1].succCnt))
	assert.Equal(t, 0, int(wd.bkts[1].failCnt))

	wd.Succ(21, 1)
	assert.Equal(t, 2, wd.cursor)
	assert.Equal(t, 20, int(wd.bkts[2].cursor))
	assert.Equal(t, float64(1.0), (wd.bkts[2].avgCost))
	assert.Equal(t, 1, int(wd.bkts[2].succCnt))
	assert.Equal(t, 0, int(wd.bkts[2].failCnt))

	wd.Succ(31, 1)
	assert.Equal(t, 3, wd.cursor)
	assert.Equal(t, 30, int(wd.bkts[3].cursor))
	assert.Equal(t, float64(1.0), (wd.bkts[3].avgCost))
	assert.Equal(t, 1, int(wd.bkts[3].succCnt))
	assert.Equal(t, 0, int(wd.bkts[3].failCnt))

	wd.Succ(41, 1)
	assert.Equal(t, 4, wd.cursor)
	assert.Equal(t, 40, int(wd.bkts[4].cursor))
	assert.Equal(t, float64(1.0), (wd.bkts[4].avgCost))
	assert.Equal(t, 1, int(wd.bkts[4].succCnt))
	assert.Equal(t, 0, int(wd.bkts[4].failCnt))

	wd.Succ(51, 1)
	assert.Equal(t, 0, wd.cursor)
	assert.Equal(t, 50, int(wd.bkts[0].cursor))
	assert.Equal(t, float64(1.0), (wd.bkts[0].avgCost))
	assert.Equal(t, 1, int(wd.bkts[0].succCnt))
	assert.Equal(t, 0, int(wd.bkts[0].failCnt))

	wd.Succ(61, 1)
	assert.Equal(t, 1, wd.cursor)
	assert.Equal(t, 60, int(wd.bkts[1].cursor))
	assert.Equal(t, float64(1.0), (wd.bkts[1].avgCost))
	assert.Equal(t, 1, int(wd.bkts[1].succCnt))
	assert.Equal(t, 0, int(wd.bkts[1].failCnt))

	wd.Succ(69, 1)
	assert.Equal(t, 1, wd.cursor)
	assert.Equal(t, 60, int(wd.bkts[1].cursor))
	assert.Equal(t, float64(1.0), (wd.bkts[1].avgCost))
	assert.Equal(t, 2, int(wd.bkts[1].succCnt))
	assert.Equal(t, 0, int(wd.bkts[1].failCnt))
}
