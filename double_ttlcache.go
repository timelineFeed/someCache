package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// 双TTL 应该是分两个key去映射内容
const (
	newKey = "%s_new"
	oldKey = "%s_old"
)

type doubleTTLCache struct {
	cache
	newExpire time.Duration
	oldExpire time.Duration
}

type doubleContent struct {
	value      any
	expireTime time.Time
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
	nKey := fmt.Sprintf(newKey, key)
	oKey := fmt.Sprintf(oldKey, key)
	// 先找一下是否有new key
	if v, ok := d.contentMap[nKey]; !ok {
		// 没有new key,那么 new old key都设置为改值
		d.contentMap[nKey] = &doubleContent{
			value:      value,
			expireTime: time.Now().Add(d.newExpire),
		}
		d.contentMap[oKey] = &doubleContent{
			value:      value,
			expireTime: time.Now().Add(d.oldExpire),
		}
		return nil
	} else {
		if content, o := v.(doubleContent); !o {
			return fmt.Errorf(ErrorAssertFailure)
		} else {
			// 将new的值更新到old，用插入的值更新new的
			newContent := &doubleContent{
				value:      value,
				expireTime: time.Now().Add(d.newExpire),
			}
			d.contentMap[nKey] = newContent
			d.contentMap[oKey] = &doubleContent{
				value:      content,
				expireTime: time.Now().Add(d.oldExpire),
			}
			return nil
		}
	}
}

func (d *doubleTTLCache) GetCache(key string) (any, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	nKey := fmt.Sprintf(newKey, key)
	oKey := fmt.Sprintf(oldKey, key)
	if v, ok := d.contentMap[nKey]; !ok {
		// new key not found
		if o, ok := d.contentMap[oKey]; !ok {
			return nil, fmt.Errorf(ErrorNotFound)
		} else {
			// old key found
			if value, can := o.(doubleContent); !can {
				return nil, errors.New(ErrorAssertFailure)
			} else {
				return value, nil
			}
		}

	} else {
		// new key found
		if value, can := v.(doubleContent); !can {
			return nil, errors.New(ErrorAssertFailure)
		} else {
			return value, nil
		}
	}
}

func (d *doubleTTLCache) handleExpire() {
	for {
		time.Sleep(1 * time.Second)
		err := d.foundExpireDel()
		if err != nil {
			fmt.Printf("in handle expire err=%+v", err)
		}
	}
}

func (d *doubleTTLCache) foundExpireDel() error {
	for k, v := range d.contentMap {
		if vale, ok := v.(doubleContent); !ok {
			return errors.New(ErrorAssertFailure)
		} else {
			if vale.expireTime.Before(time.Now()) {
				delete(d.contentMap, k)
			}
		}
	}
	return nil
}
