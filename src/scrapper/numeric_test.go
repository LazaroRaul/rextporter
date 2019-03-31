package scrapper

// import (
// 	"net"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/simelo/rextporter/src/cache"
// 	"github.com/simelo/rextporter/src/client"
// 	"github.com/simelo/rextporter/src/config"
// 	log "github.com/sirupsen/logrus"
// 	"github.com/stretchr/testify/require"
// 	"github.com/stretchr/testify/suite"
// )

// func httpHandler(w http.ResponseWriter, r *http.Request) {
// 	const jsonResponse = `
// {
//     "blockchain": {
//         "head": {
//             "seq": 58894,
//             "block_hash": "3961bea8c4ab45d658ae42effd4caf36b81709dc52a5708fdd4c8eb1b199a1f6",
//             "previous_block_hash": "8eca94e7597b87c8587286b66a6b409f6b4bf288a381a56d7fde3594e319c38a",
//             "timestamp": 1537581604,
//             "fee": 485194,
//             "version": 0,
//             "tx_body_hash": "c03c0dd28841d5aa87ce4e692ec8adde923799146ec5504e17ac0c95036362dd",
//             "ux_hash": "f7d30ecb49f132283862ad58f691e8747894c9fc241cb3a864fc15bd3e2c83d3"
//         },
//         "unspents": 38171,
//         "unconfirmed": 1,
//         "time_since_last_block": "4m46s"
//     },
//     "version": {
//         "version": "0.24.1",
//         "commit": "8798b5ee43c7ce43b9b75d57a1a6cd2c1295cd1e",
//         "branch": "develop"
//     },
//     "open_connections": 8,
//     "uptime": "6m30.629057248s",
//     "csrf_enabled": true,
//     "csp_enabled": true,
//     "wallet_api_enabled": true,
//     "gui_enabled": true,
//     "unversioned_api_enabled": false,
//     "json_rpc_enabled": false
// }
// `
// 	w.WriteHeader(http.StatusOK)
// 	if _, err := w.Write([]byte(jsonResponse)); err != nil {
// 		log.WithError(err).Panicln("unable to write response")
// 	}
// }

// type numericStatsSuit struct {
// 	suite.Suite
// 	testServer *httptest.Server
// 	cache      cache.Cache
// }

// func (suite *numericStatsSuit) SetupSuite() {
// 	suite.cache = cache.NewCache()
// 	suite.testServer = stubSkycoin()
// 	suite.testServer.Start()
// }

// func (suite *numericStatsSuit) SetupTest() {
// 	suite.cache.Reset()
// }

// func (suite *numericStatsSuit) TearDownSuite() {
// 	suite.testServer.Close()
// }

// func TestNumericStatsSuit(t *testing.T) {
// 	suite.Run(t, new(numericStatsSuit))
// }

// func stubSkycoin() *httptest.Server {
// 	l, err := net.Listen("tcp", "127.0.0.1:8080")
// 	if err != nil {
// 		log.WithError(err).Fatal("unable to create listenner")
// 	}
// 	testServer := httptest.NewUnstartedServer(http.HandlerFunc(httpHandler))
// 	testServer.Listener.Close()
// 	testServer.Listener = l
// 	return testServer
// }

// func (suite *numericStatsSuit) TestMetricBlockChainHeadSeq() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 	# All metrics to be measured.
// 	[[services.metrics]]
// 		name = "seq"
// 		url = "/api/v1/health"
// 		httpMethod = "GET"
// 		path = "/blockchain/head/seq"

// 		[services.metrics.options]
// 			type = "Counter"
// 			description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal(float64(58894), val)
// }

// func (suite *numericStatsSuit) TestMetricBlockChainHeadBlockHash() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "block_hash"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "/blockchain/head/block_hash"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal("3961bea8c4ab45d658ae42effd4caf36b81709dc52a5708fdd4c8eb1b199a1f6", val)
// }

// func (suite *numericStatsSuit) TestMetricBlockChainHeadPreviousBlockHash() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "previous_block_hash"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "/blockchain/head/previous_block_hash"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal("8eca94e7597b87c8587286b66a6b409f6b4bf288a381a56d7fde3594e319c38a", val)
// }

// func (suite *numericStatsSuit) TestMetricBlockChainHeadTimestamp() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "timestamp"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "/blockchain/head/timestamp"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal(float64(1537581604), val)
// }

// func (suite *numericStatsSuit) TestMetricBlockChainHeadFee() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "fee"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "/blockchain/head/fee"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal(float64(485194), val)
// }

// func (suite *numericStatsSuit) TestMetricBlockChainHeadVersion() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "version"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "/blockchain/head/version"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal(float64(0), val)
// }

// func (suite *numericStatsSuit) TestMetricBlockChainHeadTxBodyHash() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "tx_body_hash"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "/blockchain/head/tx_body_hash"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal("c03c0dd28841d5aa87ce4e692ec8adde923799146ec5504e17ac0c95036362dd", val)
// }

// func (suite *numericStatsSuit) TestMetricBlockChainHeadUxHash() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "ux_hash"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "/blockchain/head/ux_hash"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal("f7d30ecb49f132283862ad58f691e8747894c9fc241cb3a864fc15bd3e2c83d3", val)
// }

// func (suite *numericStatsSuit) TestMetricBlockchainUnspens() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "unspents"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "/blockchain/unspents"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal(float64(38171), val)
// }

// func (suite *numericStatsSuit) TestMetricBlockchainUnconfirmed() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "unconfirmed"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "/blockchain/unconfirmed"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal(float64(1), val)
// }

// func (suite *numericStatsSuit) TestMetricBlockchainTimeSinceLastBlock() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "time_since_last_block"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "/blockchain/time_since_last_block"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal("4m46s", val)
// }

// func (suite *numericStatsSuit) TestMetricVersionVersion() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "version"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "/version/version"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal("0.24.1", val)
// }

// func (suite *numericStatsSuit) TestMetricVersionCommit() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "commit"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "/version/commit"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal("8798b5ee43c7ce43b9b75d57a1a6cd2c1295cd1e", val)
// }

// func (suite *numericStatsSuit) TestMetricVersionBranch() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "branch"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "version/branch"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal("develop", val)
// }

// func (suite *numericStatsSuit) TestMetricOpenConnections() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "open_connections"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "open_connections"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal(float64(8), val)
// }

// func (suite *numericStatsSuit) TestMetricUptime() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "uptime"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "uptime"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal("6m30.629057248s", val)
// }

// func (suite *numericStatsSuit) TestMetricCsrfEnabled() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "csrf_enabled"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "csrf_enabled"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal(true, val)
// }

// func (suite *numericStatsSuit) TestMetricCspEnabled() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "csp_enabled"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "csp_enabled"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal(true, val)
// }

// func (suite *numericStatsSuit) TestMetricWalletApiEnabled() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "wallet_api_enabled"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "wallet_api_enabled"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal(true, val)
// }

// func (suite *numericStatsSuit) TestMetricGuiEnabled() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "gui_enabled"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "gui_enabled"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal(true, val)
// }

// func (suite *numericStatsSuit) TestMetricUnversionedApiEnabled() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "unversioned_api_enabled"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "unversioned_api_enabled"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal(false, val)
// }

// func (suite *numericStatsSuit) TestMetricJsonRpcEnabled() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	var tomlConfig = `
// 	# Service configuration.
// 	[[services]]
// 		name = "wallet"
// 		modes = ["rest_api"]
// 		scheme = "http"
// 		port = 8080
// 		basePath = ""
// 		authType = "CSRF"
// 		tokenHeaderKey = "X-CSRF-Token"
// 		genTokenEndpoint = "/api/v1/csrf"
// 		tokenKeyFromEndpoint = "csrf_token"

// 		[services.location]
// 			location = "localhost"

// 		# All metrics to be measured.
// 		[[services.metrics]]
// 			name = "json_rpc_enabled"
// 			url = "/api/v1/health"
// 			httpMethod = "GET"
// 			path = "json_rpc_enabled"

// 			[services.metrics.options]
// 				type = "Counter"
// 				description = "I am running since"
// `
// 	require := require.New(suite.T())
// 	conf, err := config.MustConfigFromRawString(tomlConfig)
// 	require.Nil(err)
// 	require.Len(conf.Services, 1)
// 	require.Len(conf.Services[0].Metrics, 1)
// 	serviceConf := conf.Services[0]
// 	metricConf := conf.Services[0].Metrics[0]
// 	var ccf client.CacheableFactory
// 	ccf, err = client.CreateAPIRestCreator(metricConf, serviceConf)
// 	require.Nil(err)
// 	cc := client.CatcherCreator{Cache: suite.cache, ClientFactory: ccf}
// 	var sc Scrapper
// 	sc, err = NewScrapper(cc, JSONParser{}, metricConf, serviceConf)
// 	require.Nil(err)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var val interface{}
// 	val, err = sc.GetMetric()
// 	require.Nil(err, "Can not get the metric")

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Equal(false, val)
// }
