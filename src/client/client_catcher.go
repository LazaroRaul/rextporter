package client

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/cache"
)

// CatcherCreator have info to create catcher client
type CatcherCreator struct {
	Cache         cache.Cache
	ClientFactory CacheableFactory
}

// CreateClient create a catcher client
func (cc CatcherCreator) CreateClient() (cl Client, err error) {
	var ccl CacheableClient
	if ccl, err = cc.ClientFactory.CreateClient(); err != nil {
		return ccl, err
	}
	return Catcher{cache: cc.Cache, dataKey: ccl.DataPath(), clientFactory: cc.ClientFactory}, nil
}

// Catcher have a client and a cache to save time loading data
type Catcher struct {
	cache         cache.Cache
	dataKey       string
	clientFactory CacheableFactory
}

// GetData return the data, can be from local cache or making the original request
func (cl Catcher) GetData(metricsCollector chan<- prometheus.Metric) (body []byte, err error) {
	if body, err = cl.cache.Get(cl.dataKey); err == nil {
		return body, err
	}
	cl.cache.Lock()
	defer cl.cache.Unlock()
	if body, err = cl.cache.Get(cl.dataKey); err == nil {
		return body, err
	}
	var ccl CacheableClient
	if ccl, err = cl.clientFactory.CreateClient(); err != nil {
		return nil, err
	}
	if body, err = ccl.GetData(metricsCollector); err == nil {
		cl.cache.Set(cl.dataKey, body)
	}
	return body, err
}
