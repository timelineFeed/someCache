package cache

import "fmt"

func Init(model int) (iCache, error) {
	switch model {
	case OCacheModel:
		return NewCache(), nil
	case TTLCacheModel:
		cache := NewTTLCache(DeFaultExpireTime)
		go cache.findOverDel()
		return cache, nil
	case DTTLCacheModel:
		cache := NewDoubleTTLCache(DeFaultExpireTime, DeFaultDeepExpireTime)
		go cache.findOverDel()
		return cache, nil
	default:
		return nil, fmt.Errorf(ErrorUnknownCache)
	}
}
