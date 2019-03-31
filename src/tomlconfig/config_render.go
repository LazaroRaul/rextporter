package tomlconfig

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"text/template"

// 	"github.com/shibukawa/configdir"
// 	"github.com/simelo/rextporter/src/util"
// 	"github.com/simelo/rextporter/src/util/file"
// 	"github.com/spf13/viper"
// )

// type templateData struct {
// 	ServicesConfigPath     string
// 	MetricsForServicesPath string
// }

// type metricsForServiceTemplateData struct {
// 	ServiceNameToMetricsConfPath map[string]string
// }

// type metricsForServiceConfigTemplateData struct {
// 	TmplData metricsForServiceTemplateData
// }

// type mainConfigData struct {
// 	mainConfigPath                  string
// 	tmplData                        templateData
// 	metricsForServiceConfigTmplData metricsForServiceConfigTemplateData
// }

// func (confData mainConfigData) ServicesConfigPath() string {
// 	return confData.tmplData.ServicesConfigPath
// }

// func (confData mainConfigData) metricsForServicesPath() string {
// 	return confData.tmplData.MetricsForServicesPath
// }

// func (confData mainConfigData) MetricsConfigPath(serviceName string) string {
// 	return confData.metricsForServiceConfigTmplData.TmplData.ServiceNameToMetricsConfPath[serviceName]
// }

// func (confData mainConfigData) MainConfigPath() string {
// 	return confData.mainConfigPath
// }

// const mainConfigFileContentTemplate = `
// servicesConfigTransport = "file" # "file" | "consulCatalog"
// servicesConfigPath = "{{.ServicesConfigPath}}"
// servicesConfigPath = "{{.ServicesConfigPath}}"
// metricsForServicesPath = "{{.MetricsForServicesPath}}"
// `

// const serviceConfigFileContentTemplate = `
// # Services configuration.
// [[services]]
//   name = "skycoin"
//   modes = ["rest_api"]
//   scheme = "http"
//   port = 8080
//   authType = "CSRF"
//   tokenHeaderKey = "X-CSRF-Token"
//   genTokenEndpoint = "/api/v1/csrf"
//   tokenKeyFromEndpoint = "csrf_token"

//   [services.location]
//     location = "localhost"
// `
// const skycoinMetricsConfigFileContentTemplate = `
// # All metrics to be measured.
// [[metrics]]
//   name = "seq"
//   url = "/api/v1/health"
//   httpMethod = "GET"
//   path = "/blockchain/head/seq"

//   [metrics.options]
//     type = "Counter"
//     description = "I am running since"

// [[metrics]]
// 	name = "openConnections"
// 	url = "/api/v1/network/connections"
// 	httpMethod = "GET"
// 	path = "/connections/unconfirmed_verify_transaction/burn_factor"

// 	[metrics.options]
// 		type = "Histogram"
// 		description = "Connections amount"

// 	[metrics.histogramOptions]
// 		buckets = [1, 2, 3]

// [[metrics]]
// 	name = "burn_factor_by_service"
// 	url = "/api/v1/network/connections"
// 	httpMethod = "GET"
// 	path = "/connections"

// 	[metrics.options]
// 		type = "Gauge"
// 		itemPath = "/unconfirmed_verify_transaction/burn_factor"
// 		description = "I am running since"

// 	[[metrics.options.labels]]
// 		name = "ip_port"
// 		path = "/address"

// [[metrics]]
// 	name = "connectedAtBySocketAndByBurnFactor"
// 	url = "/api/v1/network/connections"
// 	httpMethod = "GET"
// 	path = "/connections"

// 	[metrics.options]
// 		type = "Gauge"
// 		itemPath = "/unconfirmed_verify_transaction/burn_factor"
// 		description = "Connections amount"

// 		[[metrics.options.labels]]
// 		  name = "ip:port"
// 			path = "/address"

// 		[[metrics.options.labels]]
// 		  name = "ip:port"
// 			path = "/address"

// # TODO(denisacostaq@gmail.com):
// # if you refer(under "metrics_for_host") to a not previously defined host or metric it will be raise an error and the process will not start
// # if in all your definition you not use some host or metric the process will raise a warning and the process will start normally.
// # from https://github.com/simelo/rextporter/pull/17
// # @denisacostaq services should be match against whole templates , rather than individual metrics.
// # The match is not for hosts directly . The match is for service types . Works like this
// # metrics <- m:1 -> templates <- 1:n -> services <- m:n -> (physical | virtual) hosts
// `

// const metricsForServiceMappingConfFileContentTemplate = `
// serviceNameToMetricsConfPath = [{{range $key, $value := .}}
// 	{ {{$key}} = "{{$value}}" },{{end}}
// ]
// `

// const (
// 	systemVendorName = "simelo"
// 	// SystemProgramName is the program's name
// 	SystemProgramName                = "rextporter"
// 	mainConfigFileName               = "main.toml"
// 	servicesConfigFileName           = "services.toml"
// 	metricsForServicesConfigFileName = "metricsForServices.toml"
// 	skycoinMetricsConfigFileName     = "skycoinMetrics.toml"
// 	walletMetricsConfigFileName      = "walletMetrics.toml"
// )

// func (confData mainConfigData) existServicesConfigFile() bool {
// 	return file.ExistFile(confData.ServicesConfigPath())
// }

// // createServicesConfigFile creates the service file or return an error if any,
// // if the file already exist does no thin.
// func (confData mainConfigData) createServicesConfigFile() (err error) {
// 	generalScopeErr := "error creating service config file"
// 	if confData.existServicesConfigFile() {
// 		return nil
// 	}
// 	tmpl := template.New("serviceConfig")
// 	var templateEngine *template.Template
// 	if templateEngine, err = tmpl.Parse(serviceConfigFileContentTemplate); err != nil {
// 		errCause := "error parsing service config: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	if err = file.CreateFullPathForFile(confData.ServicesConfigPath()); err != nil {
// 		errCause := "error creating directory for service file: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	var servicesConfigFile *os.File
// 	if servicesConfigFile, err = os.Create(confData.ServicesConfigPath()); err != nil {
// 		errCause := "error creating service config file: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	if err = templateEngine.Execute(servicesConfigFile, nil); err != nil {
// 		errCause := "error writing main config file: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	return err
// }

// // createMetricsConfigFile creates the metrics file or return an error if any,
// // if the file already exist does no thin.
// func createMetricsConfigFile(metricConfPath string) (err error) {
// 	generalScopeErr := "error creating metrics config file"
// 	if file.ExistFile(metricConfPath) {
// 		return nil
// 	}
// 	tmpl := template.New("metricsConfig")
// 	var templateEngine *template.Template
// 	if templateEngine, err = tmpl.Parse(skycoinMetricsConfigFileContentTemplate); err != nil {
// 		errCause := "error parsing metrics config: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	if err = file.CreateFullPathForFile(metricConfPath); err != nil {
// 		errCause := "error creating directory for metrics file: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	var metricsConfigFile *os.File
// 	if metricsConfigFile, err = os.Create(metricConfPath); err != nil {
// 		errCause := "error creating metrics config file: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	if err = templateEngine.Execute(metricsConfigFile, nil); err != nil {
// 		errCause := "error writing metrics config file: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	return err
// }

// func (confData mainConfigData) existMetricsForServicesConfigFile() bool {
// 	return file.ExistFile(confData.tmplData.MetricsForServicesPath)
// }

// // createMetricsForServicesConfFile creates the metrics for services file or return an error if any,
// // if the file already exist does no thin.
// func (confData mainConfigData) createMetricsForServicesConfFile() (err error) {
// 	generalScopeErr := "error creating metrics for services config file"
// 	if confData.existMetricsForServicesConfigFile() {
// 		// FIXME(denisacostaq@gmail.com): check each metric file and create one of not exist
// 		return nil
// 	}
// 	tmpl := template.New("metricsForServiceConfig")
// 	var templateEngine *template.Template
// 	if templateEngine, err = tmpl.Parse(metricsForServiceMappingConfFileContentTemplate); err != nil {
// 		errCause := "error parsing metrics for services config: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	if err = file.CreateFullPathForFile(confData.metricsForServicesPath()); err != nil {
// 		errCause := "error creating directory for metrics for services file: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	var metricsForServiceConfigFile *os.File
// 	if metricsForServiceConfigFile, err = os.Create(confData.metricsForServicesPath()); err != nil {
// 		errCause := "error creating metrics for services config file: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	if err = templateEngine.Execute(metricsForServiceConfigFile, confData.metricsForServiceConfigTmplData.TmplData.ServiceNameToMetricsConfPath); err != nil {
// 		errCause := "error writing metrics for services config file: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	for key, val := range confData.metricsForServiceConfigTmplData.TmplData.ServiceNameToMetricsConfPath {
// 		if err = createMetricsConfigFile(val); err != nil {
// 			errCause := fmt.Sprintf("error creating metrics config file for service %s: %s", key, err.Error())
// 			return util.ErrorFromThisScope(errCause, generalScopeErr)
// 		}
// 	}
// 	return err
// }

// func (confData mainConfigData) existMainConfigFile() bool {
// 	return file.ExistFile(confData.MainConfigPath())
// }

// // createMainConfigFile creates the main file or return an error if any,
// // if the file already exist does no thin.
// func (confData mainConfigData) createMainConfigFile() (err error) {
// 	generalScopeErr := "error creating main config file"
// 	if confData.existMainConfigFile() {
// 		return nil
// 	}
// 	tmpl := template.New("mainConfig")
// 	var templateEngine *template.Template
// 	if templateEngine, err = tmpl.Parse(mainConfigFileContentTemplate); err != nil {
// 		errCause := "error parsing main config: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	if err = file.CreateFullPathForFile(confData.MainConfigPath()); err != nil {
// 		errCause := "error creating directory for main file: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	var mainConfigFile *os.File
// 	if mainConfigFile, err = os.Create(confData.MainConfigPath()); err != nil {
// 		errCause := "error creating main config file: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	if err = templateEngine.Execute(mainConfigFile, confData.tmplData); err != nil {
// 		errCause := "error writing main config file: " + err.Error()
// 		return util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	return err
// }

// func servicesDefaultConfigPath(conf *configdir.Config) (path string) {
// 	return file.DefaultConfigPath(servicesConfigFileName, conf)
// }

// func mainDefaultConfigPath(conf *configdir.Config) (path string) {
// 	return file.DefaultConfigPath(mainConfigFileName, conf)
// }

// func metricsForServicesDefaultConfigPath(conf *configdir.Config) (path string) {
// 	return file.DefaultConfigPath(metricsForServicesConfigFileName, conf)
// }

// func skycoinMetricsConfigPath(conf *configdir.Config) (path string) {
// 	return file.DefaultConfigPath(skycoinMetricsConfigFileName, conf)
// }

// func walletMetricsConfigPath(conf *configdir.Config) (path string) {
// 	return file.DefaultConfigPath(walletMetricsConfigFileName, conf)
// }

// func defaultTmplData(conf *configdir.Config) (tmplData templateData) {
// 	tmplData = templateData{
// 		ServicesConfigPath:     servicesDefaultConfigPath(conf),
// 		MetricsForServicesPath: metricsForServicesDefaultConfigPath(conf),
// 	}
// 	return tmplData
// }

// func defaultMetricsForServiceTmplData(conf *configdir.Config) (tmplData metricsForServiceConfigTemplateData) {
// 	tmplData = metricsForServiceConfigTemplateData{
// 		TmplData: metricsForServiceTemplateData{
// 			ServiceNameToMetricsConfPath: map[string]string{
// 				"skycoin": skycoinMetricsConfigPath(conf),
// 				"wallet":  walletMetricsConfigPath(conf),
// 			},
// 		},
// 	}
// 	return tmplData
// }

// func tmplDataFromMainFile(mainConfigFilePath string) (tmpl templateData, err error) {
// 	generalScopeErr := "error filling template data"
// 	viper.SetConfigFile(mainConfigFilePath)
// 	viper.SetConfigType("toml")
// 	if err := viper.ReadInConfig(); err != nil {
// 		errCause := fmt.Sprintln("error reading config file: ", mainConfigFilePath, err.Error())
// 		return tmpl, util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	var mainConf templateData
// 	if err := viper.Unmarshal(&mainConf); err != nil {
// 		errCause := fmt.Sprintln("can not decode the config data: ", err.Error())
// 		return tmpl, util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	tmpl.ServicesConfigPath, tmpl.MetricsForServicesPath = mainConf.ServicesConfigPath, mainConf.MetricsForServicesPath
// 	return tmpl, err
// }

// func (tmpl templateData) metricsForServicesTmplDataFromFile() (metricsForServicesTmpl metricsForServiceConfigTemplateData, err error) {
// 	generalScopeErr := "error filling template data"
// 	viper.SetConfigFile(tmpl.MetricsForServicesPath)
// 	viper.SetConfigType("toml")
// 	if err := viper.ReadInConfig(); err != nil {
// 		errCause := fmt.Sprintln("error reading config file: ", tmpl.MetricsForServicesPath, err.Error())
// 		return metricsForServicesTmpl, util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	if err := viper.Unmarshal(&(metricsForServicesTmpl.TmplData)); err != nil {
// 		errCause := fmt.Sprintln("can not decode the config data: ", err.Error())
// 		return metricsForServicesTmpl, util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	return metricsForServicesTmpl, err
// }

// func metricsForServicesTmplData(conf *configdir.Config) metricsForServiceConfigTemplateData {
// 	return defaultMetricsForServiceTmplData(conf)
// }

// func newMainConfigData(path string) (mainConf mainConfigData, err error) {
// 	generalScopeErr := "can not create main config instance"
// 	if file.IsADirectoryPath(path) {
// 		path = filepath.Join(path, mainConfigFileName)
// 	}
// 	var tmplData templateData
// 	var metricsForServiceTmplData metricsForServiceConfigTemplateData
// 	if len(path) == 0 || !file.ExistFile(path) {
// 		// TODO(denisacostaq@gmail.com): move homeConf to fn defaultTmplData
// 		var homeConf *configdir.Config
// 		if homeConf, err = file.HomeConfigFolder(systemVendorName, SystemProgramName); err != nil {
// 			errCause := "error looking for config folder under home: " + err.Error()
// 			return mainConf, util.ErrorFromThisScope(errCause, generalScopeErr)
// 		}
// 		path = mainDefaultConfigPath(homeConf)
// 		tmplData = defaultTmplData(homeConf)
// 		metricsForServiceTmplData = metricsForServicesTmplData(homeConf)
// 	} else {
// 		if tmplData, err = tmplDataFromMainFile(path); err != nil {
// 			errCause := "error reading template data from file: " + err.Error()
// 			return mainConf, util.ErrorFromThisScope(errCause, generalScopeErr)
// 		}
// 		// BUG(denisacostaq@gmail.com): if file not exist, metricsForServicesTmplDataFromFile panics
// 		if metricsForServiceTmplData, err = tmplData.metricsForServicesTmplDataFromFile(); err != nil {
// 			errCause := "error reading template data from file: " + err.Error()
// 			return mainConf, util.ErrorFromThisScope(errCause, generalScopeErr)
// 		}
// 	}
// 	if len(tmplData.ServicesConfigPath) == 0 || len(tmplData.MetricsForServicesPath) == 0 {
// 		var homeConf *configdir.Config
// 		if homeConf, err = file.HomeConfigFolder(systemVendorName, SystemProgramName); err != nil {
// 			errCause := "error looking for config folder under home: " + err.Error()
// 			return mainConf, util.ErrorFromThisScope(errCause, generalScopeErr)
// 		}
// 		tmpTmplData := defaultTmplData(homeConf)
// 		if len(tmplData.ServicesConfigPath) == 0 {
// 			tmplData.ServicesConfigPath = tmpTmplData.ServicesConfigPath
// 		}
// 		if len(tmplData.MetricsForServicesPath) == 0 {
// 			tmplData.MetricsForServicesPath = tmpTmplData.MetricsForServicesPath
// 		}
// 		metricsForServiceTmplData = metricsForServicesTmplData(homeConf)
// 	}
// 	mainConf = mainConfigData{
// 		mainConfigPath:                  path,
// 		tmplData:                        tmplData,
// 		metricsForServiceConfigTmplData: metricsForServiceTmplData,
// 	}
// 	if err = mainConf.createMainConfigFile(); err != nil {
// 		errCause := "error creating main config file: " + err.Error()
// 		return mainConf, util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	if err = mainConf.createServicesConfigFile(); err != nil {
// 		errCause := "error creating service config file: " + err.Error()
// 		return mainConf, util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	if err = mainConf.createMetricsForServicesConfFile(); err != nil {
// 		errCause := "error creating metrics for services config file: " + err.Error()
// 		return mainConf, util.ErrorFromThisScope(errCause, generalScopeErr)
// 	}
// 	return mainConf, err
// }
