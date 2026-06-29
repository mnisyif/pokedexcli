// Package pokecache implements a simple caching system that adds, gets entries
// and removes older ones
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
	cacheMap map[string]cacheEntry
	mu       *sync.Mutex
}

func NewCache(interval time.Duration) (*Cache, error) {
	newCache := &Cache{
		cacheMap: make(map[string]cacheEntry),
		mu:       &sync.Mutex{},
	}
	go newCache.reapLoop(interval)

	return newCache, nil
}
func (c *Cache) reapLoop(interval time.Duration) {
	// ticker := time.Ticker(interval)
	for {
		time.Sleep(interval)
		c.mu.Lock()
		for key, entry := range c.cacheMap {
			if time.Since(entry.createdAt) > interval {
				delete(c.cacheMap, key)
			}
		}
		c.mu.Unlock()
	}
}
