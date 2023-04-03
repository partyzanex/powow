package challenge

import (
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type cache[K comparable, V any] struct {
	store map[K]V
	ttl   map[K]time.Time
	mu    sync.RWMutex
}

func newCache[K comparable, V any]() *cache[K, V] {
	return &cache[K, V]{
		store: make(map[K]V),
		ttl:   make(map[K]time.Time),
		mu:    sync.RWMutex{},
	}
}

func (c *cache[K, V]) Get(key K) (value V, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.store[key]
	if !ok {
		return value, ErrKeyNotFound
	}

	return value, nil
}

func (c *cache[K, V]) Set(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = value
	c.ttl[key] = time.Now().Add(ttl)

	time.AfterFunc(ttl, func() {
		c.Remove(key)
	})
}

func (c *cache[K, V]) TTL(key K) (time.Time, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ttl, ok := c.ttl[key]
	if !ok {
		return time.Time{}, ErrKeyNotFound
	}

	return ttl, nil
}

func (c *cache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
	delete(c.ttl, key)
}
