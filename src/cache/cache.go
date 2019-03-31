package cache

import "sync"

// Cache mechanism for caching strings
type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, content []byte)
	Lock()
	Unlock()
	Reset()
}

type baseCache struct {
	extLocker *sync.Mutex
}

// Lock allow some external resources that uses cache if required, this is not mandatory
func (c *baseCache) Lock() {
	c.extLocker.Lock()
}

// Unlock is a required call if you use Lock
func (c *baseCache) Unlock() {
	c.extLocker.Unlock()
}

// NewCache create a new cache mechanism
func NewCache() Cache {
	return newMemCache()
}
