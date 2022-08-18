package cache

import "time"

const (
	ErrorNotFound         = "key not found"
	ErrorAssertFailure    = "assert failure"
	DeFaultExpireTime     = 10 * time.Second
	DeFaultDeepExpireTime = 24 * time.Hour
)
