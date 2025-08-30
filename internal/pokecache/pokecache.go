// Package pokecache is used to cache results from pokeapi
package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}


type Cache struct {
	cache           map[string]cacheEntry
	mu	            sync.Mutex
	interval        time.Duration
}

func NewCache(numSeconds int) *Cache {
	var cache Cache
	cache.interval = time.Duration(numSeconds) * time.Second
	cache.cache = make(map[string]cacheEntry)
	go cache.reapLoop()
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	var entry cacheEntry
	entry.createdAt = time.Now()
	entry.val = val
	c.mu.Lock()
	c.cache[key] = entry
	c.mu.Unlock()	
}

func (c *Cache) Get(key string) ([]byte, bool) {	
	c.mu.Lock()
	entry, exists := c.cache[key]
	c.mu.Unlock()	
	if !exists {
		return nil, false
	}	
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mu.Lock()
		for key, value := range c.cache {
			if time.Since(value.createdAt) < c.interval {
				continue;
			}
			delete(c.cache, key)
		}
		c.mu.Unlock()
	}
}
