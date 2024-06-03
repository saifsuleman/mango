package main

import (
	"sync"
	"time"
)

const (
    WriteThrough = 1 << iota
    Loading
    WriteBack
)

type LoadingFunction[K comparable, V any] func(key K) (V, error)
type WeighingFunction[K comparable, V any] func(key K, value V) int64
type WriteThroughFunction[K comparable, V any] func(key K, value V) error

type cacheEntry[T any] struct {
	value      T
	expiration int64
}

type LoadingCache[K comparable, V any] struct {
	maximumSize          int64
	maximumWeight        int64
	expirationSeconds    int64
	loadingFunction      LoadingFunction[K, V]
	writeThroughFunction WriteThroughFunction[K, V]
	weighingFunction     WeighingFunction[K, V]
	cache                map[K]*cacheEntry[V]
	currentWeight        int64
	mu                   sync.RWMutex
}


func (c *LoadingCache[K, V]) addEntry(key K, value V) {
	if c.weighingFunction != nil {
		c.currentWeight += c.weighingFunction(key, value)
	}

	c.evictIfNecessary()

	c.cache[key] = &cacheEntry[V]{
		value:      value,
		expiration: time.Now().Unix() + c.expirationSeconds,
	}
}

func (c *LoadingCache[K, V]) shouldEvict() bool {
    if c.maximumSize > 0 && int64(len(c.cache)) >= c.maximumSize {
        return true
    }

    if c.weighingFunction == nil || c.maximumWeight <= 0 {
        return false
    }

    return c.currentWeight > c.maximumWeight
}

func (c *LoadingCache[K, V]) evictIfNecessary() {
	c.mu.Lock()
	defer c.mu.Unlock()

    for c.shouldEvict() {
        // implement LRU eviction policy
    }
}

func (c *LoadingCache[K, V]) isExpired(entry *cacheEntry[V]) bool {
	return time.Now().Unix() > entry.expiration
}

func (c *LoadingCache[K, V]) Get(key K) (V, error) {
	c.mu.RLock()
	entry, exists := c.cache[key]
	if exists && !c.isExpired(entry) {
		c.mu.RUnlock()
		return entry.value, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists = c.cache[key]
	if exists && !c.isExpired(entry) {
		return entry.value, nil
	}

	value, err := c.loadingFunction(key)
	if err != nil {
		var zeroValue V
		return zeroValue, err
	}

	c.addEntry(key, value)
	return value, nil
}

func (c *LoadingCache[K, V]) Invalidate(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, exists := c.cache[key]; exists {
		if c.weighingFunction != nil {
			c.currentWeight -= c.weighingFunction(key, entry.value)
		}
		delete(c.cache, key)
	}
}

func (c *LoadingCache[K, V]) Refresh(key K) (V, error) {
	if _, exists := c.cache[key]; exists {
		c.mu.Lock()
		defer c.mu.Unlock()

		value, err := c.loadingFunction(key)
		if err != nil {
			var zeroValue V
			return zeroValue, err
		}

		c.cache[key].expiration = time.Now().Unix() + c.expirationSeconds
		c.cache[key].value = value
		return value, nil
	}

	return c.Get(key)
}

func (c *LoadingCache[K, V]) IsLoaded(key K) bool {
	if _, exists := c.cache[key]; exists {
		return true
	}
	return false
}
