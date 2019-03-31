package scrapper

// WorkPool is use to run scrapper task in goroutines
var WorkPool *Pool

func init() {
	WorkPool = NewPool(6)
	WorkPool.StartDispatcher()
}
