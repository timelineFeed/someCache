package cache

import (
	"fmt"
	"sync"
)

type cache struct {
	lock       sync.Mutex
	contentMap map[string]any
}

func NewCache(size int) *cache {
	return &cache{
		lock:       sync.Mutex{},
		contentMap: make(map[string]any),
	}
}

func (c *cache) GetCache(key string) (any, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if value, ok := c.contentMap[key]; !ok {
		return nil, fmt.Errorf(ErrorNotFound)
	} else {
		return value, nil
	}
}

func (c *cache) SetCache(key string, value any) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.contentMap[key] = value
	return nil
}
