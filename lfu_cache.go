package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

/*
LFU 按照反问频率由高到低排序
同频率按照访问时间排序
*/

type lfuCache struct {
	lock sync.RWMutex
	cap  int
	list list.List
	data map[string]*list.Element
}

type lfuContent struct {
	key       string
	value     any
	frequency int
	time      time.Time
}

func (cache *lfuCache) GetCache(key string) (any, error) {
	// 获取数据
	e, ok := cache.data[key]
	if !ok {
		return nil, errors.New(ErrorNotFound)
	}
	content, ok := e.Value.(*lfuContent)
	if !ok {
		return nil, errors.New(ErrorAssertFailure)
	}
	cache.lock.Lock()
	defer cache.lock.Unlock()
	// 成功frequency++,对链表重新排序
	content.frequency++
	content.time = time.Now()
	// 排序
	err := cache.lfuSort(e, content.frequency)
	if err != nil {
		return nil, err
	}
	return content.value, nil
}

func (cache *lfuCache) SetCache(key string, value any) error {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	//判断key是否已经存在
	e, ok := cache.data[key]
	if ok {
		content, ok2 := e.Value.(*lfuContent)
		if !ok2 {
			return errors.New(ErrorAssertFailure)
		}
		content.value = value
		if e.Prev() == nil && e.Next() == nil {
			cache.list.Remove(e)
			cache.list.PushFront(content)
		} else if e.Prev() != nil {
			prev := e.Prev()
			cache.list.Remove(e)
			cache.list.InsertAfter(content, prev)
		} else if e.Next() != nil {
			next := e.Next()
			cache.list.Remove(e)
			cache.list.InsertBefore(content, next)
		}
	}
	// 容量已经满
	content := &lfuContent{
		key:       key,
		value:     value,
		frequency: 0,
		time:      time.Now(),
	}
	if cache.list.Len() > cache.cap {
		// 删除最后一个元素
		last := cache.list.Back()
		cache.list.Remove(last)
		c, ok := last.Value.(*lfuContent)
		if !ok {
			return errors.New(ErrorAssertFailure)
		}
		delete(cache.data, c.key)
	}
	// 插入一个元素
	e, err := cache.lfuInsert(content)
	if err != nil {
		return err
	}
	cache.data[key] = e
	return nil
}

// lfuSort 前面元素最后一个不频率大于e的元素a,a与e交换
func (cache *lfuCache) lfuSort(e *list.Element, frequency int) error {
	ok := true
	content := new(lfuContent)
	targetElement := e.Prev()
	if targetElement == nil {
		return nil
	}
	for {
		content, ok = targetElement.Value.(*lfuContent)
		if !ok {
			return errors.New(ErrorAssertFailure)
		}
		if content.frequency > frequency {
			break
		}
		targetElement = targetElement.Prev()
		if targetElement == nil {
			break
		}
	}
	if targetElement == e.Prev() {
		// 前面的元素都比它大
		return nil
	}
	if targetElement == nil {
		// 前面的元素都比它小
		v := cache.list.Remove(e)
		cache.list.PushFront(v)
	}
	// 交换元素
	eBefore := e.Prev()
	tBefore := targetElement.Prev()
	tAfter := new(list.Element)
	if tBefore == nil {
		tAfter = targetElement.Next()
	}
	eContent := cache.list.Remove(e)
	tContent := cache.list.Remove(targetElement)
	cache.list.InsertAfter(eContent, tBefore)
	if tBefore == nil {
		cache.list.InsertBefore(eBefore, tAfter)
	}
	cache.list.InsertBefore(tContent, eBefore)
	return nil

}

func (cache *lfuCache) lfuInsert(content *lfuContent) (*list.Element, error) {
	//找出最往前的一个frequent为1的元素
	mark := cache.list.Back()
	for {
		c, ok := mark.Value.(*lfuContent)
		if !ok {
			return nil, errors.New(ErrorAssertFailure)
		}
		if c.frequency > 0 {
			break
		}
		mark = mark.Prev()
		if mark == nil {
			break
		}
	}
	if mark == nil {

		return cache.list.PushFront(content), nil
	}

	return cache.list.InsertAfter(content, mark), nil
}
