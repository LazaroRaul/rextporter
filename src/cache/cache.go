package cache

// Storage mechanism for caching strings
type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, content []byte)
	Reset()
}

type CacheCleaner interface {
	ClearCache()
}

func NewCache() Cache {
	return newMemCache()
}
