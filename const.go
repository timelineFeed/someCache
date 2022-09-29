package cache

import "time"

const (
	ErrorNotFound      = "在缓存中找不到该值"
	ErrorAssertFailure = "断言失败"
	ErrorUnknownCache  = "未知缓存类型"
	ErrorReflectFailed = "reflect failed"
)

const (
	DeFaultExpireTime     = 10 * time.Second
	DeFaultDeepExpireTime = 24 * time.Hour
)

const (
	OCacheModel    = iota //普通cache
	TTLCacheModel         // ttl
	DTTLCacheModel        // 双ttl
)
