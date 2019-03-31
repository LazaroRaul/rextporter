package rxt

import (
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/core"
	"github.com/simelo/rextporter/src/util"
)

// ASTDefEnv buildsthe syntax tree
type ASTDefEnv struct {
	Options config.OptionsMap
}

func NewASTDefEnv() *ASTDefEnv {
	return &ASTDefEnv{
		Options: config.NewOptionsMap(),
	}
}
// NewServiceScraper ...
func (env *ASTDefEnv) NewServiceScraper() (core.RextServiceScraper, error) {
	return &ASTDefScraperDataset{
		SupportedServiceNames: nil,
		SupportedStackNames:   nil,
		Definitions:           make(map[string]interface{}),
		Sources:               nil,
		Options:               config.NewOptionsMap(),
	}, nil
}

// NewAuthStrategy ...
func (env *ASTDefEnv) NewAuthStrategy(authtype string, options core.RextKeyValueStore) (core.RextAuth, error) {
	auth := ASTDefAuth{
		AuthType: authtype,
		Options:  config.NewOptionsMap(),
	}
	if err := util.MergeStoresInplace(auth.Options, options); err != nil {
		return nil, err
	}
	return &auth, nil
}

// NewMetricsExtractor ...
func (env *ASTDefEnv) NewMetricsExtractor(scraperType string, options core.RextKeyValueStore, metrics []core.RextMetricDef) (core.RextMetricsExtractor, error) {
	astMetrics := make([]*ASTDefMetric, len(metrics))
	extractor := ASTDefExtract{
		Type:    scraperType,
		Metrics: astMetrics,
		Options: config.NewOptionsMap(),
	}
	for idx, m := range metrics {
		if astMetric, isAST := m.(*ASTDefMetric); isAST {
			astMetrics[idx] = astMetric
		} else {
			astMetrics[idx] = &ASTDefMetric{
				Name:        m.GetMetricName(),
				Type:        m.GetMetricType(),
				Description: m.GetMetricDescription(),
				Labels:      m.GetMetricLabels(),
				Options:     config.NewOptionsMap(),
			}
			if err := util.MergeStoresInplace(astMetrics[idx].Options, m.GetOptions()); err != nil {
				return nil, err
			}
		}
	}
	return &extractor, nil
}

// NewMetricsDatasource ...
func (env *ASTDefEnv) NewMetricsDatasource(srcType string) core.RextDataSource {
	return &ASTDefSource{
		Method:   "",
		Type:     srcType,
		Location: "",
		Scrapers: nil,
		Options:  config.NewOptionsMap(),
	}
}

// RegisterScraperForServices ...
func (env *ASTDefEnv) RegisterScraperForServices(s core.RextServiceScraper, services ...string) error {
	if ds, isAST := s.(*ASTDefScraperDataset); isAST {
		ds.SupportedServiceNames = services
		return nil
	}
	return core.ErrInvalidType
}

// RegisterScraperForStacks ...
func (env *ASTDefEnv) RegisterScraperForStacks(s core.RextServiceScraper, stacks ...string) error {
	if ds, isAST := s.(*ASTDefScraperDataset); isAST {
		ds.SupportedStackNames = stacks
		return nil
	}
	return core.ErrInvalidType
}

// GetOptions ...
func (env *ASTDefEnv) GetOptions() core.RextKeyValueStore {
	return env.Options
}

// ASTDefScraperDataset parse tree node
type ASTDefScraperDataset struct {
	SupportedServiceNames []string
	SupportedStackNames   []string
	Definitions           map[string]interface{}
	Sources               []*ASTDefSource
	Options               config.OptionsMap
}

// AddAuthStrategy ...
func (ds *ASTDefScraperDataset) AddAuthStrategy(auth core.RextAuth, name string) {
	ds.Definitions[name] = auth
}

// AddSource ...
func (ds *ASTDefScraperDataset) AddSource(source core.RextDataSource) {
	if astSrc, isAST := source.(*ASTDefSource); isAST {
		ds.Sources = append(ds.Sources, astSrc)
	}
}

// AddSources ...
func (ds *ASTDefScraperDataset) AddSources(sources ...core.RextDataSource) {
	for _, source := range sources {
		ds.AddSource(source)
	}
}

// GetOptions ...
func (ds *ASTDefScraperDataset) GetOptions() core.RextKeyValueStore {
	return ds.Options
}

// ASTDefAuth parse tree node
type ASTDefAuth struct {
	AuthType string
	Options  config.OptionsMap
}

// GetOptions ...
func (auth *ASTDefAuth) GetOptions() core.RextKeyValueStore {
	return auth.Options
}

// GetOptions ...
func (auth *ASTDefAuth) GetAuthType() string {
	return auth.AuthType
}

// ASTDefSource parse tree node
type ASTDefSource struct {
	Method   string
	Type     string
	Location string
	Scrapers []*ASTDefExtract
	Options  config.OptionsMap
}

// SetBaseURL ...
func (src *ASTDefSource) SetBaseURL(url string) {
	// Do nothing. No info about this in AST
}

// GetMethod ...
func (src *ASTDefSource) GetMethod() string {
	return src.Method
}

// SetMethod ...
func (src *ASTDefSource) SetMethod(s string) {
	src.Type = s
}

// GetResourceLocation ...

func (src *ASTDefSource) GetResourceLocation() string {
	return src.Location
}

// SetResourceLocation ...
func (src *ASTDefSource) SetResourceLocation(s string) error {
	src.Location = s
	return nil
}

// ActivateScraper ...
func (src *ASTDefSource) ActivateScraper(scraper core.RextMetricsExtractor) (err error) {
	if extAST, isAST := scraper.(*ASTDefExtract); isAST {
		src.Scrapers = append(src.Scrapers, extAST)
	} else {
		err = core.ErrInvalidType
	}
	return
}

// GetOptions ...
func (src *ASTDefSource) GetOptions() core.RextKeyValueStore {
	return src.Options
}

// ASTDefExtract parse tree node
type ASTDefExtract struct {
	Type    string
	Metrics []*ASTDefMetric
	Options config.OptionsMap
}

// Apply ...
func (scraper *ASTDefExtract) Apply(rule core.RextMetricDef) (core.RextMetricDef, error) {
	var m *ASTDefMetric
	if astRule, isAST := rule.(*ASTDefMetric); isAST {
		m = astRule
	} else {
		m = NewMetricFromDef(rule)
	}
	scraper.Metrics = append(scraper.Metrics, m)
	return m, nil
}

// ApplyMany ...
func (scraper *ASTDefExtract) ApplyMany(rules []core.RextMetricDef) (newRules []core.RextMetricDef, err error) {
	newRules = make([]core.RextMetricDef, len(rules))
	for idx, md := range rules {
		if newRules[idx], err = scraper.Apply(md); err != nil {
			return
		}
	}
	return
}

// ExtractMetrics ...
func (scraper *ASTDefExtract) ExtractMetrics(target interface{}) ([]core.RextMetric, error) {
	return nil, nil
}

// GetOptions ...
func (scraper *ASTDefExtract) GetOptions() core.RextKeyValueStore {
	return scraper.Options
}

// ASTDefMetric parse tree node
type ASTDefMetric struct {
	Type        string
	Name        string
	Description string
	Labels      []string
	Options     config.OptionsMap
}

func NewMetricFromDef(m core.RextMetricDef) *ASTDefMetric {
	return &ASTDefMetric{
		Name:        m.GetMetricName(),
		Type:        m.GetMetricType(),
		Description: m.GetMetricDescription(),
		Labels:      m.GetMetricLabels(),
	}
}

// GetMetricName ...
func (m *ASTDefMetric) GetMetricName() string {
	return m.Name
}

// GetMetricType ...
func (m *ASTDefMetric) GetMetricType() string {
	return m.Type
}

// GetMetricDescription ...
func (m *ASTDefMetric) GetMetricDescription() string {
	return m.Description
}

// GetMetricLabels ...
func (m *ASTDefMetric) GetMetricLabels() []string {
	return m.Labels
}

// SetMetricName ...
func (m *ASTDefMetric) SetMetricName(s string) {
	m.Name = s
}

// SetMetricType ...
func (m *ASTDefMetric) SetMetricType(s string) {
	m.Type = s
}

// SetMetricDescription ...
func (m *ASTDefMetric) SetMetricDescription(s string) {
	m.Description = s
}

// SetMetricLabels ...
func (m *ASTDefMetric) SetMetricLabels(labels []string) {
	m.Labels = labels
}

// GetOptions ...
func (m *ASTDefMetric) GetOptions() core.RextKeyValueStore {
	return m.Options
}
