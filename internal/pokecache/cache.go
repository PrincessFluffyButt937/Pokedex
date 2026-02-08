package pokecache

import (
	"sync"
	"time"
)

func NewCache(interval time.Duration) *Cache {
	c := Cache{
		data:     make(map[string]cacheEntry),
		interval: interval,
	}
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			c.reapLoop()
		}
	}()
	return &c
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	data     map[string]cacheEntry
	interval time.Duration
	mu       sync.RWMutex
}

func (s *Cache) Add(key string, val []byte) {
	t := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := cacheEntry{
		createdAt: t,
		val:       val,
	}
	s.data[key] = entry
}

func (s *Cache) Get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, exists := s.data[key]
	if !exists {
		return nil, false
	}
	return data.val, true
}

func (s *Cache) reapLoop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, val := range s.data {
		age := time.Since(val.createdAt)
		if age > s.interval {
			delete(s.data, key)
		}

	}
}
