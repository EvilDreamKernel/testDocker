package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	//	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	monitorDirectory = flag.String("fs.monitor-directory", "files/",
		"Directory to monitor")
	listenAddress = flag.String("web.listen-address", ":9145",
		"Address to listen on for telemetry")
	metricsPath = flag.String("web.metrics-path", "/metrics",
		"Path under which to expose metrics")
)

func main() {
	flag.Parse()
	router := mux.NewRouter()
	promRegistry := prometheus.NewRegistry()
	fileExporterCollector := newFileExporterCollector(*monitorDirectory)
	promRegistry.MustRegister(fileExporterCollector)
	router.Handle(*metricsPath, promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}))

	log.Println("Listening on " + *listenAddress)
	log.Println("Monitoring "+ *monitorDirectory + " directory")
	log.Println("Metrics Path: " + *metricsPath)
	log.Fatal(http.ListenAndServe(*listenAddress, router))
}
