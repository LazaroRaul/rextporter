package config

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/simelo/rextporter/src/common"
	"github.com/spf13/viper"
)

// Host is a concept to grab information about a running server, for example:
// where is it http://localhost:1234 (Location + : + Port), what auth kind you need to use?
// what is the header key you in which you need to send the token, and so on.
type Host struct {
	Ref                  string
	Location             string `json:"location"`
	Port                 int    `json:"port"`
	AuthType             string `json:"auth_type"`
	TokenHeaderKey       string `json:"token_header_key"`
	GenTokenEndpoint     string `json:"gen_token_endpoint"`
	TokenKeyFromEndpoint string `json:"token_key_from_endpoint"`
}

// isValidUrl tests a string to determine if it is a valid URL or not.
func isValidURL(toTest string) bool {
	if _, err := url.ParseRequestURI(toTest); err != nil {
		return false
	}
	return true
}

func (host Host) validate() (errs []error) {
	if len(host.Ref) == 0 {
		errs = append(errs, errors.New("ref is required in host"))
	}
	if len(host.Location) == 0 {
		errs = append(errs, errors.New("location is required in host"))
	}
	if !isValidURL(host.Location) {
		errs = append(errs, errors.New("location is not a valid url in host"))
	}
	if !isValidURL(host.URIToGetToken()) {
		errs = append(errs, errors.New("location + port can not form a valid uri in host"))
	}
	if host.Port < 0 || host.Port > 65535 {
		errs = append(errs, errors.New("port number should be between 0 and 65535 in host"))
	}
	if strings.Compare(host.AuthType, "CSRF") == 0 && len(host.TokenHeaderKey) == 0 {
		errs = append(errs, errors.New("TokenHeaderKey is required if you are using CSRF"))
	}
	if strings.Compare(host.AuthType, "CSRF") == 0 && len(host.GenTokenEndpoint) == 0 {
		errs = append(errs, errors.New("GenTokenEndpoint is required if you are using CSRF"))
	}
	if strings.Compare(host.AuthType, "CSRF") == 0 && len(host.TokenKeyFromEndpoint) == 0 {
		errs = append(errs, errors.New("TokenKeyFromEndpoint is required if you are using CSRF"))
	}
	return errs
}

// MetricOptions keep information you about the metric, mostly the type(Counter, Gauge, Summary, and Histogram)
type MetricOptions struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (mo MetricOptions) validate() (errs []error) {
	if len(mo.Type) == 0 {
		errs = append(errs, errors.New("type is required in metric"))
	}
	return errs
}

// Metric keep the metric name as an instance of MetricOptions
type Metric struct {
	Name    string        `json:"name"`
	Options MetricOptions `json:"options"`
}

func (metric Metric) validate() (errs []error) {
	if len(metric.Name) == 0 {
		errs = append(errs, errors.New("name is required in metric"))
	}
	errs = append(errs, metric.Options.validate()...)
	return errs
}

// Link is a concept who map properties of a Metric in a Host, for example, you can define
// some hosts some metrics and in Link your specific the properties of a giving metric in
// a giving host, for example, the Url and the json path(Path) from where you can scrap the information.
type Link struct {
	HostRef    string `json:"host_ref"`
	MetricRef  string `json:"metric_ref"`
	URL        string `json:"url"`
	HTTPMethod string `json:"http_method"`
	Path       string `json:"path,omitempty"`
}

func (link Link) validate() (errs []error) {
	if len(link.HostRef) == 0 {
		errs = append(errs, errors.New("HostRef is required in Link(metric fo host)"))
	}
	if len(link.MetricRef) == 0 {
		errs = append(errs, errors.New("HostRef is required in Link(metric fo host)"))
	}
	if len(link.URL) == 0 {
		errs = append(errs, errors.New("url is required"))
	}
	if len(link.HTTPMethod) == 0 {
		errs = append(errs, errors.New("HttpMethod is required in Link(metric fo host)"))
	}
	if len(link.Path) == 0 {
		errs = append(errs, errors.New("path is required in Link(metric fo host)"))
	}
	host, hostNotFound := Config().FindHostByRef(link.HostRef)
	if hostNotFound != nil {
		errs = append(errs, hostNotFound)
	} else {
		if !isValidURL(host.URIToGetMetric(link)) {
			errs = append(errs, errors.New("can not create a valid uri under link"))
		}
		errs = append(errs, host.validate()...)
	}
	metric, metricNotFound := Config().findMetricByRef(link.MetricRef)
	if metricNotFound != nil {
		errs = append(errs, metricNotFound)
	} else {
		errs = append(errs, metric.validate()...)
	}
	return errs
}

// MetricName will return a name for a metric in a host
func (link Link) MetricName() string {
	return link.HostRef + "_" + link.MetricRef
}

// MetricDescription will look for the metric associated trough ref and return the description
func (link Link) MetricDescription() (string, error) {
	const generalScopeErr = "error getting metric description"
	var metric Metric
	var err error
	if metric, err = Config().findMetricByRef(link.MetricRef); err != nil {
		errCause := fmt.Sprintln("can not find the metric", err.Error())
		return "", common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return metric.Options.Description, err
}

// RootConfig is the top level node for the config tree, it has a list of hosts, a list of metrics
// and a list of links(MetricsForHost, says how a metric is mapped in a host).
type RootConfig struct {
	Hosts          []Host   `json:"hosts"`
	Metrics        []Metric `json:"metrics"`
	MetricsForHost []Link   `json:"metrics_for_host"`
}

var rootConfig RootConfig

func (conf RootConfig) validate() {
	var errs []error
	for _, host := range conf.Hosts {
		errs = append(errs, host.validate()...)
	}
	for _, metric := range conf.Metrics {
		errs = append(errs, metric.validate()...)
	}
	for _, mhost := range conf.MetricsForHost {
		errs = append(errs, mhost.validate()...)
	}
	if len(errs) != 0 {
		log.Println("some errors found")
		for _, err := range errs {
			log.Println(err.Error())
		}
		log.Panicln()
	}
}

// Config TODO(denisacostaq@gmail.com): make a singleton
func Config() RootConfig {
	//if b, err := json.MarshalIndent(rootConfig, "", " "); err != nil {
	//	log.Println("Error marshaling:", err)
	//} else {
	//	os.Stdout.Write(b)
	//	log.Println("\n\n\n\n\n")
	//}
	// TODO(denisacostaq@gmail.com): Make it a singleton
	return rootConfig
}

// NewConfigFromRawString allow you to define a `.toml` config in the fly, a raw string with the "config content"
func NewConfigFromRawString(strConf string) error {
	const generalScopeErr = "error creating a config instance"
	viper.SetConfigType("toml")
	buff := bytes.NewBuffer([]byte(strConf))
	if err := viper.ReadConfig(buff); err != nil {
		errCause := fmt.Sprintln("can not read the buffer", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	rootConfig = RootConfig{}
	if err := viper.Unmarshal(&rootConfig); err != nil {
		errCause := fmt.Sprintln("can not decode the config data", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	rootConfig.validate()
	return nil
}

// NewConfigFromFilePath TODO(denisacostaq@gmail.com): Fill some data structures for efficient lookup from ref to host for example
func NewConfigFromFilePath(path string) error {
	const generalScopeErr = "error creating a config instance"
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		errCause := fmt.Sprintln("error reading config file:", path, err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err := viper.Unmarshal(&rootConfig); err != nil {
		errCause := fmt.Sprintln("can not decode the config data", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	rootConfig.validate()
	return nil
}

// FindHostByRef will return a host where you can match the host.Ref with the ref parameter
// or an error if not found.
func (conf RootConfig) FindHostByRef(ref string) (host Host, err error) {
	found := false
	for _, host = range conf.Hosts {
		found = strings.Compare(host.Ref, ref) == 0
		if found {
			return
		}
	}
	if !found {
		errCause := fmt.Sprintln("can not find a host for Ref:", ref)
		err = errors.New(errCause)
	}
	return Host{}, err
}

// findMetricByRef will return a metric where you can match the metric.Ref with the ref parameter
// or an error if not found.
func (conf RootConfig) findMetricByRef(ref string) (metric Metric, err error) {
	found := false
	for _, metric = range conf.Metrics {
		found = strings.Compare(metric.Name, ref) == 0
		if found {
			return
		}
	}
	if !found {
		errCause := fmt.Sprintln("can not find a host for Ref:", ref)
		err = errors.New(errCause)
	}
	return Metric{}, err
}

// FindMetricType will return the metric type through the metric related with ref
func (link Link) FindMetricType() (metricType string, err error) {
	const generalScopeErr = "error looking for metric type"
	var metric Metric
	if metric, err = Config().findMetricByRef(link.MetricRef); err != nil {
		errCause := fmt.Sprintln("can not find metric by ref:", link.MetricRef, err.Error())
		return metricType, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metricType = metric.Options.Type
	return metricType, err
}

// URIToGetMetric build the URI from where you will to get metric information
func (host Host) URIToGetMetric(metricInHost Link) string {
	return host.Location + ":" + strconv.Itoa(host.Port) + metricInHost.URL
}

// URIToGetToken build the URI from where you will to get the token
func (host Host) URIToGetToken() string {
	return host.Location + ":" + strconv.Itoa(host.Port) + host.TokenKeyFromEndpoint
}

// FilterLinksByHost will return all links where you can match the host.Ref with link.HostRef
func (conf RootConfig) FilterLinksByHost(host Host) []Link {
	var links []Link
	for _, link := range conf.MetricsForHost {
		if strings.Compare(host.Ref, link.HostRef) == 0 {
			links = append(links, link)
		}
	}
	return links
}
