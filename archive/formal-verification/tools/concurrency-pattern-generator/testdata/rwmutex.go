package main

import "sync"

type Cache struct {
	mu sync.RWMutex
	data map[string]interface{}
}

func (c *Cache) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}
