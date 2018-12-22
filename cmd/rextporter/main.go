package main

import (
	"flag"
	"os"

	"github.com/simelo/rextporter/src/core"
	"github.com/simelo/rextporter/src/exporter"
	"github.com/simelo/rextporter/src/toml2config"
	"github.com/simelo/rextporter/src/tomlconfig"
	log "github.com/sirupsen/logrus"
)

func main() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	mainConfigFile := flag.String("config", "/home/adacosta/.config/simelo/rextporter/main.toml", "Metrics main config file path.")
	defaultListenPort := 8080
	listenPort := flag.Uint("port", uint(defaultListenPort), "Listen port.")
	defaultHandlerEndpoint := "/metrics"
	handlerEndpoint := flag.String("handler", defaultHandlerEndpoint, "Handler endpoint.")
	flag.Parse()
	conf, err := tomlconfig.ReadConfigFromFileSystem(*mainConfigFile)
	if err != nil {
		log.WithError(err).Errorln("error reading config from file system")
		os.Exit(1)
	}
	var rootConf core.RextRoot
	rootConf, err = toml2config.Fill(conf)
	if err != nil {
		log.WithError(err).Errorln("error filling config info")
		os.Exit(1)
	}
	exporter.MustExportMetrics(*handlerEndpoint, uint16(*listenPort), rootConf)
	waitForEver := make(chan bool)
	<-waitForEver
}
