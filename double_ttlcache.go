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

var _ iTTLCache = &doubleTTLCache{}

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
		// 没有new key,只设置 new key
		d.contentMap[nKey] = &doubleContent{
			value:      value,
			expireTime: time.Now().Add(d.newExpire),
		}
		return nil
	} else {
		if content, o := v.(*doubleContent); !o {
			return fmt.Errorf(ErrorAssertFailure)
		} else {
			// 将new的值更新到old，用插入的值更新new的
			newContent := &doubleContent{
				value:      value,
				expireTime: time.Now().Add(d.newExpire),
			}
			d.contentMap[nKey] = newContent
			d.contentMap[oKey] = &doubleContent{
				value:      content.value,
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
			if content, can := o.(*doubleContent); !can {
				return nil, errors.New(ErrorAssertFailure)
			} else {
				return content.value, nil
			}
		}

	} else {
		// new key found
		if content, can := v.(*doubleContent); !can {
			return nil, errors.New(ErrorAssertFailure)
		} else {
			return content.value, nil
		}
	}
}

func (d *doubleTTLCache) GetCacheExpire(key string) (any, time.Time, error) {
	cache, err := d.GetCache(key)
	if err != nil {
		return nil, time.Time{}, err
	}
	content := cache.(*doubleContent)
	return content.value, content.expireTime, nil
}

func (d *doubleTTLCache) handleExpire() {
	for {
		time.Sleep(1 * time.Second)
		err := d.findOverDel()
		if err != nil {
			fmt.Printf("in handle expire err=%+v", err)
		}
	}
}

func (d *doubleTTLCache) findOverDel() error {
	for k, v := range d.contentMap {
		if vale, ok := v.(doubleContent); !ok {
			fmt.Println(ErrorAssertFailure)
			return errors.New(ErrorAssertFailure)
		} else {
			if vale.expireTime.Before(time.Now()) {
				delete(d.contentMap, k)
			}
		}
	}
	return nil
}
