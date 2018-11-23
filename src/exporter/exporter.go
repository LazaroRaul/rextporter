package exporter

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
	log "github.com/sirupsen/logrus"
)

func findMetricsName(metricsData string) (metricsNames []string) {
	rex := regexp.MustCompile(`# TYPE [a-zA-Z_:][a-zA-Z0-9_:]*`)
	metricsNameLines := rex.FindAllString(metricsData, -1)
	metricsNames = make([]string, len(metricsNameLines))
	for idx, metricsNameLine := range metricsNameLines {
		metricsNameLineColumns := strings.Split(metricsNameLine, " ")
		// FIXME(denissacostaq@gmail.com): be careful indexing here
		metricsNames[idx] = metricsNameLineColumns[2]
	}
	return metricsNames
}

func appendPrefixForMetrics(prefix string, metricsData string) ([]byte, error) {
	metricsName := findMetricsName(metricsData)
	for _, metricName := range metricsName {
		repl := strings.NewReplacer(
			"# HELP "+metricName+" ", "# HELP "+prefix+"_"+metricName+" ",
			"# TYPE "+metricName+" ", "# TYPE "+prefix+"_"+metricName+" ",
		)
		metricsData = repl.Replace(metricsData)
		metricsData = strings.Replace(metricsData, "\n"+metricName, "\n"+prefix+"_"+metricName, -1)
	}
	if len(metricsName) == 0 {
		err := fmt.Errorf("data from %s not appear to be from a metrics(trough prometheus instrumentation) endpoint", string(prefix))
		log.WithError(err).Errorln("append prefix error, content ignored")
	}
	return []byte(metricsData), nil
}

func exposedMetricsMiddleware(metricsMiddleware []MetricMiddleware, promHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getCustomData := func() (data []byte, err error) {
			recorder := httptest.NewRecorder()
			for _, cl := range metricsMiddleware {
				if exposedMetricsData, err := cl.client.GetExposedMetrics(); err != nil {
					log.WithError(err).Error("error getting metrics from service " + cl.client.Name)
				} else {
					if prefixed, err := appendPrefixForMetrics(cl.client.Name, string(exposedMetricsData)); err == nil {
						if count, err := recorder.Write(prefixed); err != nil || count != len(prefixed) {
							if err != nil {
								log.WithError(err).Errorln("error writing prefixed content")
							}
							if count != len(prefixed) {
								log.WithFields(log.Fields{
									"wrote":    count,
									"required": len(prefixed),
								}).Errorln("no enough content wrote")
							}
						}
					}
				}
			}
			if data, err = ioutil.ReadAll(recorder.Body); err != nil {
				log.WithError(err).Errorln("can not read recorded custom data")
				return nil, err
			}
			return data, nil
		}
		getDefaultData := func() (data []byte, err error) {
			generalScopeErr := "error reding default data"
			recorder := httptest.NewRecorder()
			promHandler.ServeHTTP(recorder, r)
			var reader io.ReadCloser
			switch recorder.Header().Get("Content-Encoding") {
			case "gzip":
				reader, err = gzip.NewReader(recorder.Body)
				if err != nil {
					errCause := fmt.Sprintln("can not create gzip reader.", err.Error())
					return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
				}
			default:
				reader = ioutil.NopCloser(bytes.NewReader(recorder.Body.Bytes()))
			}
			defer reader.Close()
			if data, err = ioutil.ReadAll(reader); err != nil {
				log.WithError(err).Errorln("can not read recorded default data")
				return nil, err
			}
			return data, nil
		}
		var allData []byte
		if defaultData, err := getDefaultData(); err != nil {
			log.WithError(err).Errorln("error getting default data")
		} else {
			allData = append(allData, defaultData...)
		}
		if customData, err := getCustomData(); err != nil {
			log.WithError(err).Errorln("error getting custom data")
		} else {
			allData = append(allData, customData...)
		}
		w.Header().Set("Content-Type", "text/plain")
		if allData == nil {
			allData = []byte("a")
		}
		if count, err := w.Write(allData); err != nil || count != len(allData) {
			if err != nil {
				log.WithError(err).Errorln("error writing data")
			}
			if count != len(allData) {
				log.WithFields(log.Fields{
					"wrote":    count,
					"required": len(allData),
				}).Errorln("no enough content wrote")
			}
		}
	})
}

// ExportMetrics will read the config from mainConfigFile if any or use a default one.
func ExportMetrics(handlerEndpoint string, listenPort uint16, conf config.RootConfig) (srv *http.Server) {
	if collector, err := newSkycoinCollector(conf); err != nil {
		log.WithError(err).Panicln("Can not create metrics")
	} else {
		prometheus.MustRegister(collector)
	}
	metricsMiddleware, err := createMetricsMiddleware(conf)
	if err != nil {
		log.WithError(err).Panicln("Can not create forward_metrics metrics")
	}
	port := fmt.Sprintf(":%d", listenPort)
	srv = &http.Server{Addr: port}
	http.Handle(
		handlerEndpoint,
		gziphandler.GzipHandler(exposedMetricsMiddleware(metricsMiddleware, promhttp.Handler())))
	go func() {
		log.Infoln(fmt.Sprintf("Starting server in port %d, path %s ...", listenPort, handlerEndpoint))
		log.WithError(srv.ListenAndServe()).Errorln("unable to start the server")
	}()
	return srv
}

// TODO(denisacostaq@gmail.com): you can use a NewProcessCollector, NewGoProcessCollector, make a blockchain collector sense?
