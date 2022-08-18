package cache

type iCache interface {
	GetCache(string) (any, error)
	SetCache(string, any) error
}
