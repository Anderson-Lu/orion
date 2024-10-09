package cache

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type IFetcher interface {
	FetchFromOrigin(ctx context.Context, key string) (interface{}, error)
}

type TickerCacheItem struct {
	data interface{}
	ts   int64
}

func (tci *TickerCacheItem) IsExpire(expireSec int) bool {
	if expireSec <= 0 {
		return false
	}
	return time.Now().Unix()-tci.ts > int64(expireSec)
}

type TickerCache struct {
	fetcher   IFetcher
	expire    int
	cached    sync.Map
	sg        singleflight.Group
	cleanTick time.Duration
}

type TickerCacheOption func(*TickerCache)

func New(opts ...TickerCacheOption) *TickerCache {
	tc := &TickerCache{}
	for _, v := range opts {
		v(tc)
	}
	go tc.autoElimination()
	return tc
}

func (t *TickerCache) autoElimination() {
	if t.cleanTick <= 0 || t.expire <= 0 {
		return
	}
	ticker := time.NewTicker(t.cleanTick)
	for range ticker.C {
		t.cached.Range(func(key, value any) bool {
			if value == nil {
				t.cached.Delete(key)
				return true
			}
			c, ok := value.(*TickerCacheItem)
			if !ok || c == nil {
				t.cached.Delete(key)
				return true
			}
			if c.IsExpire(t.expire) {
				t.cached.Delete(key)
				return true
			}
			return true
		})
	}
}

func WithOrigin(f IFetcher) TickerCacheOption {
	return func(tc *TickerCache) {
		tc.fetcher = f
	}
}

func WithTTL(expireSeconds int) TickerCacheOption {
	return func(tc *TickerCache) {
		tc.expire = expireSeconds
	}
}

func WithElimination(tickSeconds int) TickerCacheOption {
	return func(tc *TickerCache) {
		tc.cleanTick = time.Second * time.Duration(tickSeconds)
	}
}

func (t *TickerCache) Get(ctx context.Context, key string) (interface{}, error) {
	c, ok := t.cached.Load(key)
	if !ok || c == nil {
		return t.Refresh(ctx, key)
	}
	ca, ok := c.(*TickerCacheItem)
	if !ok || ca == nil {
		return t.Refresh(ctx, key)
	}
	if ca.IsExpire(t.expire) {
		return t.Refresh(ctx, key)
	}
	return ca.data, nil
}

func (t *TickerCache) Refresh(ctx context.Context, key string) (interface{}, error) {
	data, err, _ := t.sg.Do(key, func() (interface{}, error) {
		info, err := t.fetcher.FetchFromOrigin(ctx, key)
		if err != nil {
			return nil, err
		}
		t.cached.Store(key, &TickerCacheItem{
			data: info,
			ts:   time.Now().Unix(),
		})
		return info, nil
	})
	return data, err
}
