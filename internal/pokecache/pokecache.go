package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	CacheEntries map[string]cacheEntry
	mux          *sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{map[string]cacheEntry{}, &sync.RWMutex{}}
	if interval == 0 {
		return Cache{}
	}

	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	c.CacheEntries[key] = cacheEntry{time.Now(), val}
	c.mux.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.RLock()
	entry, exists := c.CacheEntries[key]
	c.mux.RUnlock()

	if !exists {
		return []byte{}, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	reap := func() {
		now := time.Now()

		c.mux.Lock()
		for k, v := range c.CacheEntries {
			if v.createdAt.Before(now.Add(-interval)) {
				delete(c.CacheEntries, k)
			}
		}
		c.mux.Unlock()
	}

	for range ticker.C {
		reap()
	}
}
