package libs

import (
	"sync"
)

type Cache struct {
	items map[string][]byte
	mu    sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string][]byte),
	}
}

func (c *Cache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.items[key]
	return value, exists
}
