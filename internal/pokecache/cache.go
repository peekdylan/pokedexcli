package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu    *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// NewCache creates a new cache with a reap loop that runs at the given interval
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache: make(map[string]cacheEntry),
		mu:    &sync.Mutex{},
	}

	go c.reapLoop(interval)

	return c
}

// Add adds a new entry to the cache
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get retrieves an entry from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

// reapLoop runs in the background and removes old entries
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.reap(interval)
	}
}

// reap removes entries older than the given duration
func (c *Cache) reap(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.cache {
		if now.Sub(entry.createdAt) > interval {
			delete(c.cache, key)
		}
	}
}
