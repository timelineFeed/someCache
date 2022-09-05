package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type ttlCache struct {
	cache
	expire time.Duration
}

type ttlContent struct {
	value      any
	expireTime time.Time
}

func NewTTLCache(size uint64, expire time.Duration) *ttlCache {
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
		value:      value,
		expireTime: time.Now().Add(c.expire),
	}
	return nil
}

func (c *ttlCache) handleExpire() {
	for {
		time.Sleep(1 * time.Second)
		err := c.FindOverDel()
		if err != nil {
			fmt.Printf("in handle expire err=%+v", err)
		}
	}

}

// FindOverDel 找出超时的key并删除
func (c *ttlCache) FindOverDel() error {

	for k, v := range c.contentMap {
		if o, ok := v.(ttlContent); !ok {
			return errors.New(ErrorAssertFailure)
		} else {
			if o.expireTime.Before(time.Now()) {
				delete(c.contentMap, k)
			}
		}
	}
	return nil
}
