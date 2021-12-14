package main

import (
	"flag"
	// "github.com/go-kit/kit/log/level"
	"net/http"
	"newrelic_exporter/config"
	"newrelic_exporter/exporter"
	"newrelic_exporter/newrelic"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	// "github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	// "github.com/prometheus/common/promlog"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "newrelic_exporter.yml", "Config file path. Defaults to 'newrelic_exporter.yml'")
	flag.Parse()

	log.SetLevel(log.DebugLevel)

	cfg, err := config.GetConfig(configFile)

	api := newrelic.NewAPI(cfg)

	exp := exporter.NewExporter(api, cfg)

	prometheus.MustRegister(exp)

	http.Handle(cfg.MetricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
<head><title>NewRelic exporter</title></head>
<body>
<h1>NewRelic exporter</h1>
<p><a href='` + cfg.MetricPath + `'>Metrics</a></p>
</body>
</html>
`))
	})

	log.Infof("Listening on %s.", cfg.ListenAddress)
	err = http.ListenAndServe(cfg.ListenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("HTTP server stopped.")
}
