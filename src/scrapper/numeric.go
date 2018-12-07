package scrapper

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/util"
)

// Numeric scrapper can get numeric(gauges or counters) metrics
type Numeric struct {
	baseScrapper
}

func newNumeric(cf client.Factory, p BodyParser, path, jobName, instanceName, datasource string) Scrapper {
	return Numeric{
		baseScrapper: baseScrapper{
			clientFactory: cf,
			parser:        p,
			jsonPath:      path,
			jobName:       jobName,
			instanceName:  instanceName,
			datasource:    datasource,
		},
	}
}

// GetMetric returns a single number with the metric value, is a counter or a gauge
func (n Numeric) GetMetric(metricsCollector chan<- prometheus.Metric) (val interface{}, err error) {
	const generalScopeErr = "error scrapping numeric(gauge|counter) metric"
	var iBody interface{}
	if iBody, err = getData(n.clientFactory, n.parser, metricsCollector); err != nil {
		errCause := "numeric client can not decode the body"
		return val, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if val, err = n.parser.pathLookup(n.jsonPath, iBody); err != nil {
		errCause := fmt.Sprintln("can not get node: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return val, err
}
