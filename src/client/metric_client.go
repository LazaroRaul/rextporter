package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/oliveagle/jsonpath"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
)

// BaseClient have common data to be shared through embedded struct in those type who implement the
// client.Client interface
type BaseClient struct {
	req     *http.Request
	service config.Service
}

// HistogramClientOptions hold the necessary reference bucket to create an histogram
type HistogramClientOptions struct {
	Buckets []float64
}

// MetricClient implements the getRemoteInfo method from `client.Client` interface by using some `.toml` config parameters
// like for example: where is the host? it should be a GET, a POST or some other? ...
// sa NewMetricClient method.
type MetricClient struct {
	BaseClient
	token                  string
	metricJPath            string
	histogramClientOptions HistogramClientOptions
}

// NewMetricClient will put all the required info to be able to do http requests to get the remote data.
func NewMetricClient(metric config.Metric, service config.Service) (client *MetricClient, err error) {
	const generalScopeErr = "error creating a client to get a metric from remote endpoint"
	if strings.Compare(service.Mode, config.ServiceTypeAPIRest) != 0 {
		return client, errors.New("can not create an api rest metric client from a service of type " + service.Mode)
	}
	client = new(MetricClient)
	client.BaseClient.service = service
	client.metricJPath = metric.Path
	if strings.Compare(metric.Options.Type, config.KeyTypeHistogram) == 0 {
		client.histogramClientOptions = HistogramClientOptions{
			Buckets: metric.HistogramOptions.Buckets,
		}
	}
	client.BaseClient.req, err = http.NewRequest(metric.HTTPMethod, client.service.URIToGetMetric(metric), nil)
	if err != nil {
		errCause := fmt.Sprintln("can not create the request: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return client, nil
}

func (client *MetricClient) resetToken() (err error) {
	const generalScopeErr = "error making resetting the token"
	client.token = ""
	var clientToken *TokenClient
	if clientToken, err = newTokenClient(client.service); err != nil {
		errCause := fmt.Sprintln("can not find a host: ", err.Error())
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var data []byte
	if data, err = clientToken.getRemoteInfo(); err != nil {
		errCause := fmt.Sprintln("can make the request to get a token: ", err.Error())
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var jsonData interface{}
	if err = json.Unmarshal(data, &jsonData); err != nil {
		errCause := fmt.Sprintln("can not decode the body: ", string(data), " ", err.Error())
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var val interface{}
	jPath := "$" + strings.Replace(client.service.TokenKeyFromEndpoint, "/", ".", -1)
	if val, err = jsonpath.JsonPathLookup(jsonData, jPath); err != nil {
		errCause := fmt.Sprintln("can not locate the path: ", err.Error())
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	tk, ok := val.(string)
	if !ok {
		errCause := fmt.Sprintln("unable the get the token as a string value")
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	client.token = tk
	if len(client.token) == 0 {
		errCause := fmt.Sprintln("unable the get a not null(empty) token")
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return nil
}

func (client *MetricClient) getRemoteInfo() (data []byte, err error) {
	const generalScopeErr = "error making a server request to get metric from remote endpoint"
	client.req.Header.Set(client.service.TokenHeaderKey, client.token)
	getData := func() (data []byte, err error) {
		httpClient := &http.Client{}
		var resp *http.Response
		if resp, err = httpClient.Do(client.req); err != nil {
			errCause := fmt.Sprintln("can not do the request: ", err.Error())
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			errCause := fmt.Sprintf("no success response, status %s", resp.Status)
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			errCause := fmt.Sprintln("can not read the body: ", err.Error())
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		return data, nil
	}
	if data, err = getData(); err != nil {
		// log.Println("can not do the request:", err.Error(), "trying with a new token...")
		if err = client.resetToken(); err != nil {
			errCause := fmt.Sprintln("can not reset the token: ", err.Error())
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		if data, err = getData(); err != nil {
			errCause := fmt.Sprintln("can not do the request after a token reset neither: ", err.Error())
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
	}
	return data, nil
}

// GetMetric returns the metric previously bound through config parameters like:
// url(endpoint), json path, type and so on.
func (client *MetricClient) GetMetric() (val interface{}, err error) {
	const generalScopeErr = "error getting metric data"
	var data []byte
	if data, err = client.getRemoteInfo(); err != nil {
		return nil, util.ErrorFromThisScope(err.Error(), generalScopeErr)
	}
	var jsonData interface{}
	if err = json.Unmarshal(data, &jsonData); err != nil {
		errCause := fmt.Sprintln("can not decode the body: ", string(data), " ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	jPath := "$" + strings.Replace(client.metricJPath, "/", ".", -1)
	if val, err = jsonpath.JsonPathLookup(jsonData, jPath); err != nil {
		errCause := fmt.Sprintln("can not locate the path: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return val, nil
}

// HistogramValue hold the required values to create a histogram metric, the Count, Sum and buckets.
type HistogramValue struct {
	Count   uint64
	Sum     float64
	Buckets map[float64]uint64
}

func newHistogram(buckets []float64) HistogramValue {
	val := HistogramValue{
		Count:   0,
		Sum:     0,
		Buckets: make(map[float64]uint64, len(buckets)),
	}
	for _, bucket := range buckets {
		val.Buckets[bucket] = 0
	}
	return val
}

// GetHistogramValue will return a histogram data structure from the remote endpoint.
func (client *MetricClient) GetHistogramValue() (val HistogramValue, err error) {
	generalScopeErr := "error getting histogram values"
	var metric interface{}
	if metric, err = client.GetMetric(); err != nil {
		errCause := fmt.Sprintln("can not get the metric data: ", err.Error())
		return val, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metricCollection, okMetricCollection := metric.([]interface{})
	if !okMetricCollection {
		errCause := fmt.Sprintln("can not assert the metric type as slice")
		return val, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	val = newHistogram(client.histogramClientOptions.Buckets)
	for _, metricItem := range metricCollection {
		val.Count++
		metricVal, okMetricVal := metricItem.(float64)
		if !okMetricVal {
			errCause := fmt.Sprintln("can not assert the metric value to type float")
			return val, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		val.Sum += metricVal
		for _, bucket := range client.histogramClientOptions.Buckets {
			if bucket <= metricVal {
				val.Buckets[bucket]++
			}
		}
	}
	return val, err
}
