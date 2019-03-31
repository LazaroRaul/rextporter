package cache

import (
	"fmt"
	"sync"
)

// MemCache implements a cache mechanism in memory bay using a map
type MemCache struct {
	baseCache
	vals      map[string][]byte
	dataMutex *sync.RWMutex
}

func newMemCache() *MemCache {
	return &MemCache{
		baseCache: baseCache{extLocker: &sync.Mutex{}},
		vals:      make(map[string][]byte),
		dataMutex: &sync.RWMutex{},
	}
}

// Get return the cached value by a giving key, error if this val not found
func (c *MemCache) Get(key string) (val []byte, err error) {
	c.dataMutex.RLock()
	defer c.dataMutex.RUnlock()
	if val, ok := c.vals[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("value not found for key %s", key)
}

// Set save data with a giving key in the cache sistem
func (c *MemCache) Set(key string, data []byte) {
	c.dataMutex.Lock()
	defer c.dataMutex.Unlock()
	c.vals[key] = data
}

// Reset clear al the cached data
func (c *MemCache) Reset() {
	c.dataMutex.Lock()
	defer c.dataMutex.Unlock()
	for k := range c.vals {
		delete(c.vals, k)
	}
}
