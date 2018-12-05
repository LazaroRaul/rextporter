package config

import (
	"container/list"
	"errors"
	"fmt"
)

const (
	// ServiceTypeAPIRest is the key you should define in the config file for a service who request remote data
	// and uses this to build the metrics.
	ServiceTypeAPIRest = "rest_api"
	// ServiceTypeProxy is the key you should define in the config file for a service to work like a middleware/forward_metrics.
	ServiceTypeProxy = "forward_metrics"
)

// Service is a concept to grab information about a running server, for example:
// where is it http://localhost:1234 (Location + : + Port), what auth kind you need to use?
// what is the header key you in which you need to send the token, and so on.
type Service struct {
	Name  string   `json:"name"`
	Modes []string `json:"modes"`
	// Scheme is http or https
	Scheme               string   `json:"scheme"`
	Port                 uint16   `json:"port"`
	BasePath             string   `json:"basePath"`
	MetricsToForwardPath string   `json:"metrics_to_forward"`
	AuthType             string   `json:"authType"`
	TokenHeaderKey       string   `json:"tokenHeaderKey"`
	GenTokenEndpoint     string   `json:"genTokenEndpoint"`
	TokenKeyFromEndpoint string   `json:"tokenKeyFromEndpoint"`
	Location             Server   `json:"location"`
	Metrics              []Metric `json:"metrics"`
}

// MetricName returns a promehteus style name for the giving metric name.
func (srv Service) MetricName(metricName string) string {
	// return prometheus.BuildFQName(SystemProgramName, srv.Name, metricName)
	return metricName
}

// URIToGetMetric build the URI from where you will to get metric information
func (srv Service) URIToGetMetric(metric Metric) string {
	return fmt.Sprintf("%s://%s:%d%s%s", srv.Scheme, srv.Location.Location, srv.Port, srv.BasePath, metric.URL)
}

// URIToGetExposedMetric build the URI from where you will to get the exposed metrics.
func (srv Service) URIToGetExposedMetric() string {
	return fmt.Sprintf("%s://%s:%d%s%s", srv.Scheme, srv.Location.Location, srv.Port, srv.BasePath, srv.MetricsToForwardPath)
}

// URIToGetToken build the URI from where you will to get the token
func (srv Service) URIToGetToken() string {
	return fmt.Sprintf("%s://%s:%d%s%s", srv.Scheme, srv.Location.Location, srv.Port, srv.BasePath, srv.GenTokenEndpoint)
}

// FilterMetricsByType will return all the metrics who match with the 't' parameter in this service.
func (srv Service) FilterMetricsByType(t string) (metrics []Metric) {
	tmpMetrics := list.New()
	for _, metric := range srv.Metrics {
		if metric.Options.Type == t {
			tmpMetrics.PushBack(metric)
		}
	}
	metrics = make([]Metric, tmpMetrics.Len())
	idxLink := 0
	for it := tmpMetrics.Front(); it != nil; it = it.Next() {
		metrics[idxLink] = it.Value.(Metric)
		idxLink++
	}
	return metrics
}

// CountMetricsByType will return the number of metrics who match with the 't' parameter in this service.
func (srv Service) CountMetricsByType(t string) (amount int) {
	for _, metric := range srv.Metrics {
		if metric.Options.Type == t {
			amount++
		}
	}
	return
}

func (srv Service) validateProxy() (errs []error) {
	if !isValidURL(srv.URIToGetExposedMetric()) {
		errs = append(errs, errors.New("can not create a valid url to get the exposed metric"))
	}
	if len(srv.MetricsToForwardPath) == 0 {
		errs = append(errs, errors.New("you need to define metricsToForwardPath if you enable proxy(forward_metrics) mode"))
	}
	return errs
}

func (srv Service) validateAPIRest() (errs []error) {
	if !isValidURL(srv.URIToGetToken()) {
		errs = append(errs, errors.New("can not create a valid url to get token: "+srv.URIToGetToken()))
	}
	for _, metric := range srv.Metrics {
		if !isValidURL(srv.URIToGetMetric(metric)) {
			errs = append(errs, errors.New("can not create a valid url to get metric: "+srv.URIToGetMetric(metric)))
		}
	}
	if srv.AuthType == "CSRF" && len(srv.TokenHeaderKey) == 0 {
		errs = append(errs, errors.New("TokenHeaderKey is required if you are using CSRF"))
	}
	if srv.AuthType == "CSRF" && len(srv.TokenKeyFromEndpoint) == 0 {
		errs = append(errs, errors.New("TokenKeyFromEndpoint is required if you are using CSRF"))
	}
	if srv.AuthType == "CSRF" && len(srv.GenTokenEndpoint) == 0 {
		errs = append(errs, errors.New("GenTokenEndpoint is required if you are using CSRF"))
	}
	return errs
}

func (srv Service) validate() (errs []error) {
	if len(srv.Name) == 0 {
		errs = append(errs, errors.New("name is required in service"))
	}
	if len(srv.Scheme) == 0 {
		errs = append(errs, errors.New("scheme is required in service"))
	}
	if srv.Port < 1 || srv.Port > 65535 {
		errs = append(errs, errors.New("port must be betwen 1 and 65535"))
	}
	if len(srv.Modes) == 0 {
		errs = append(errs, fmt.Errorf("you you have to define at least a service mode, possibles are: %s or %s", ServiceTypeAPIRest, ServiceTypeProxy))
	}
	for _, mode := range srv.Modes {
		switch mode {
		case ServiceTypeProxy:
			errs = append(errs, srv.validateProxy()...)
		case ServiceTypeAPIRest:
			errs = append(errs, srv.validateAPIRest()...)
		default:
			if len(mode) == 0 {
				errs = append(errs, fmt.Errorf("mode is required in service"))
			} else {
				errs = append(errs, fmt.Errorf("mode allow instances of %s or %s only", ServiceTypeAPIRest, ServiceTypeProxy))
			}
		}
	}
	for _, metric := range srv.Metrics {
		errs = append(errs, metric.validate()...)
	}
	errs = append(errs, srv.Location.validate()...)
	return errs
}
