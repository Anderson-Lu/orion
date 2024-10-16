package ratelimit

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type IRateLimitor interface {
	Allow() bool
}

type Config struct {
	Key             string
	Cap             int64
	TokensPerSecond int64
}

type RateLimiters struct {
	mls map[string]IRateLimitor
	mu  sync.RWMutex
}

func NewRateLimiters() *RateLimiters {
	return &RateLimiters{mls: make(map[string]IRateLimitor)}
}

func (r *RateLimiters) Register(config *Config) error {

	if config.Key == "" || config.Cap <= 0 || config.TokensPerSecond < 0 {
		return errors.New("bad config")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	ra := time.Duration(int64(float64(1.0)/float64(config.TokensPerSecond))) * time.Second
	r.mls[config.Key] = rate.NewLimiter(rate.Every(ra), int(config.Cap))
	return nil
}

func (r *RateLimiters) Allow(key string) bool {
	if key == "" {
		return false
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	if lim, ok := r.mls[key]; ok && lim != nil {
		return lim.Allow()
	}

	return false
}

type TokenRateLimitor struct {
	r *rate.Limiter
}

func (t *TokenRateLimitor) Allow() bool {
	return t.r.Allow()
}
