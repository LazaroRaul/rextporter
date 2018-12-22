package toml2config

import (
	"fmt"

	"github.com/simelo/rextporter/src/core"
	"github.com/simelo/rextporter/src/memconfig"
	"github.com/simelo/rextporter/src/tomlconfig"
	log "github.com/sirupsen/logrus"
)

type metricName2Metric map[string]tomlconfig.Metric
type serviceName2MetricName2Metric map[string]metricName2Metric

func buildMetricsMapping(conf tomlconfig.RootConfig) (metricsMapping serviceName2MetricName2Metric) {
	metricsMapping = make(serviceName2MetricName2Metric)
	for _, srv := range conf.Services {
		mtrName2Metric := make(metricName2Metric)
		for _, mtr := range srv.Metrics {
			mtrName2Metric[mtr.Name] = mtr
		}
		metricsMapping[srv.Name] = mtrName2Metric
	}
	return metricsMapping
}

func createService(srv tomlconfig.Service, metricsMapping serviceName2MetricName2Metric) (service core.RextServiceDef, err error) {
	mtrN2Metric := metricsMapping[srv.Name]
	service = &memconfig.Service{}
	basePath := fmt.Sprintf("%s://%s:%d", srv.Protocol, srv.Location.Location, srv.Port)
	service.SetBasePath(basePath)
	srvOpts := service.GetOptions()
	srvOpts.SetString(core.OptKeyRextServiceDefJobName, srv.Name)
	srvOpts.SetString(core.OptKeyRextServiceDefInstanceName, fmt.Sprintf("%s:%d", srv.Location.Location, srv.Port))
	auth := &memconfig.HTTPAuth{}
	auth.SetAuthType(srv.AuthType)
	authOpts := auth.GetOptions()
	authOpts.SetString(core.OptKeyRextAuthDefTokenHeaderKey, srv.TokenHeaderKey)
	authOpts.SetString(core.OptKeyRextAuthDefTokenKeyFromEndpoint, srv.TokenKeyFromEndpoint)
	authOpts.SetString(core.OptKeyRextAuthDefTokenGenEndpoint, srv.GenTokenEndpoint)
	service.SetAuthForBaseURL(auth)
	for _, resPath := range srv.ResourcePaths {
		var resDef core.RextResourceDef
		switch resPath.PathType {
		case "rest_api":
			resDef = createResourceFrom4API(mtrN2Metric, resPath)
			resDef.SetType(resPath.PathType)
			resDef.SetResourceURI(resPath.Path)
			resOpts := resDef.GetOptions()
			resOpts.SetString(core.OptKeyRextResourceDefHTTPMethod, resPath.HttpMethod)
		case "metrics_fordwader":
			resDef = createResourceFrom4ExposedMetrics(resPath)
		default:
			log.WithField("resource_path_type", resPath.PathType).Errorln("valid types are rest_api or metrics_fordwader")
			return service, core.ErrKeyInvalidType
		}
		service.AddResource(resDef)
	}
	return service, err
}

func createResourceFrom4API(mtrN2Metric map[string]tomlconfig.Metric, resPath tomlconfig.ResourcePath) (resDef core.RextResourceDef) {
	resDef = &memconfig.ResourceDef{}
	resDef.SetType(resPath.PathType)
	resDef.SetResourceURI(resPath.Path)
	resOpts := resDef.GetOptions()
	resOpts.SetString(core.OptKeyRextResourceDefHTTPMethod, resPath.HttpMethod)
	for _, mtrName := range resPath.MetricNames {
		mtr /*, foundMetric*/ := mtrN2Metric[mtrName]
		// if !foundMetric {
		// 	continue
		// }
		nodeSolver := &memconfig.NodeSolver{MType: resPath.NodeSolverType}
		nodeSolver.SetNodePath(mtr.Path)
		metric := &memconfig.MetricDef{}
		metric.SetMetricName(mtr.Name)
		metric.SetMetricType(mtr.Options.Type)
		metric.SetMetricDescription(mtr.Options.Description)
		mtrOpts := metric.GetOptions()
		for _, tomlLabel := range mtr.Options.Labels {
			label := &memconfig.LabelDef{}
			label.SetName(tomlLabel.Name)
			lns := &memconfig.NodeSolver{}
			lns.SetNodePath(tomlLabel.Path)
			label.SetNodeSolver(lns)
			metric.AddLabel(label)
		}
		if mtr.Options.Type == core.KeyMetricTypeHistogram {
			mtrOpts.SetObject(core.OptKeyRextMetricDefHMetricBuckets, mtr.HistogramOptions.Buckets)
		}
		if mtr.Options.Labels != nil && len(mtr.Options.Labels) > 0 {
			mtrOpts.SetString(core.OptKeyRextMetricDefVecItemPath, mtr.Options.ItemPath)
		}
		metric.SetNodeSolver(nodeSolver)
		resDef.AddMetricDef(metric)
	}
	return resDef
}

func createResourceFrom4ExposedMetrics(resPath tomlconfig.ResourcePath) (resDef core.RextResourceDef) {
	resDef = &memconfig.ResourceDef{}
	resDef.SetType(resPath.PathType)
	resDef.SetResourceURI(resPath.Path)
	return resDef
}

func Fill(conf tomlconfig.RootConfig) (root core.RextRoot, err error) {
	root = &memconfig.RootConfig{}
	metricsMapping := buildMetricsMapping(conf)
	for _, srv := range conf.Services {
		var service core.RextServiceDef
		if service, err = createService(srv, metricsMapping); err != nil {
			log.WithError(err).Errorln("can not fill service info")
			return root, err
		}
		root.AddService(service)
	}
	return root, err
}
