package cache

import (
	"fmt"
	"sync"
	"time"
)

type ttlCache struct {
	cache
	expire time.Time
}

type ttlContent struct {
	key    string
	value  any
	expire time.Time
}

func NewTTLCache(size uint64, expire time.Time) *ttlCache {
	return &ttlCache{
		cache: cache{
			lock:       sync.Mutex{},
			contentMap: make(map[string]any),
		},
		expire: expire,
	}
}

func (c *ttlCache) GetCache(key string) (any, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if value, ok := c.contentMap[key]; !ok {
		return nil, fmt.Errorf(ErrorNotFound)
	} else {
		if v, o := value.(ttlContent); !o {
			return nil, fmt.Errorf(ErrorAssertFailure)
		} else {
			return v.value, nil
		}

	}
}

func (c *ttlCache) SetCache(key string, value any) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.contentMap[key] = &ttlContent{
		key:    key,
		value:  value,
		expire: c.expire,
	}
	return nil
}

func (c *ttlCache) handleExpire() {
	//TODO 借鉴gingame util

}
