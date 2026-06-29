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
