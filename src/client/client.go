package client

// Client to get remote data.
type Client interface {
	// GetData will get tha date based on a URL(but can be a cached value for example).
	GetData() (body []byte, err error)
}

// CacheableClient should return a key(DataPath) for catching resource values
type CacheableClient interface {
	Client
	DataPath() string
}

type baseCacheableClient struct {
	dataPath string
}

// DataPath is the endpoint in the case of http clients
func (cl baseCacheableClient) DataPath() string {
	return cl.dataPath
}

// TODO(denisacostaq@gmail.com): check out http://localhost:6060/pkg/github.com/prometheus/client_golang/api/#NewClient
