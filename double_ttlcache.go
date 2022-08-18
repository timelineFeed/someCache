package cache

import (
	"fmt"
	"sync"
	"time"
)

type doubleTTLCache struct {
	cache
	newExpire time.Duration
	oldExpire time.Duration
}

type doubleContent struct {
	key           string
	newValue      any
	oldValue      any
	newExpireTime time.Time
	oldExpireTime time.Time
}

func NewDoubleTTLCache(new, old time.Duration) *doubleTTLCache {
	return &doubleTTLCache{
		cache: cache{
			lock:       sync.Mutex{},
			contentMap: make(map[string]any),
		},
		newExpire: new,
		oldExpire: old,
	}
}

func (d *doubleTTLCache) SetCache(key string, value any) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if v, ok := d.contentMap[key]; !ok {
		d.contentMap[key] = &doubleContent{
			key:           key,
			newValue:      value,
			oldValue:      value,
			newExpireTime: time.Now().Add(d.newExpire),
			oldExpireTime: time.Now().Add(d.oldExpire),
		}
		return nil
	} else {
		if content, o := v.(doubleContent); !o {
			return fmt.Errorf(ErrorAssertFailure)
		} else {
			newContent := &doubleContent{
				key:           key,
				newValue:      value,
				oldValue:      content.newValue,
				newExpireTime: time.Now().Add(d.newExpire),
				oldExpireTime: time.Now().Add(d.oldExpire),
			}
			d.contentMap[key] = newContent
			return nil
		}
	}
}

func (d *doubleTTLCache) GetCache(key string) (any, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if v, ok := d.contentMap[key]; !ok {
		return nil, fmt.Errorf(ErrorNotFound)
	} else {
		if content, o := v.(*doubleContent); !o {
			return nil, fmt.Errorf(ErrorAssertFailure)
		} else {
			if content.newValue != nil {
				return content.newValue, nil
			} else {
				return content.oldValue, nil
			}
		}
	}
}

func (d *doubleTTLCache) handleExpire() {
	//TODO
}
