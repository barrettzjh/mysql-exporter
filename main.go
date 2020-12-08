package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"mysql-exporter/module"
	"net/http"
)

type Exporter struct {
	newStatusMetric map[string]*prometheus.Desc
}

func newECSMetric(metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName("mysql", "status", metricName),
		docString, labels, nil,
	)
}

func newExporter() *Exporter {
	desc := make(map[string]*prometheus.Desc)
	status := module.GetMysqlStatus()
	desc["disk"] = newECSMetric("disk", "disk", []string{"mysql", "database"})

	for _, m := range status{
		for k := range m{
			desc[k] = newECSMetric(k, k, []string{"mysql"})
		}
	}
	return &Exporter{
		newStatusMetric: desc,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.newStatusMetric {
		ch <- m
	}
}

func (e Exporter) Collect(ch chan<- prometheus.Metric) {
	var newValue,newV float64
	var err error
	status := module.GetMysqlStatus()
	storage := module.GetMysqlStorage()

	for index, m := range status{
		for key, value := range m{
			newV, err = module.StringToFloat(value)
			if err != nil{
				continue
			}
			ch <- prometheus.MustNewConstMetric(e.newStatusMetric[key], prometheus.GaugeValue, newV, index)
		}
	}
	for i, n := range storage{

		for k, v := range n{
			newValue, err = module.StringToFloat(v)
			if err != nil{
				continue
			}
			ch <- prometheus.MustNewConstMetric(e.newStatusMetric["disk"], prometheus.GaugeValue, newValue, i, k)
		}
	}
}

var (
	listenAddress   = flag.String("telemetry.address", ":" + module.Configure.Port, "Address on which to expose metrics.")
	metricsEndpoint = flag.String("telemetry.endpoint", module.Configure.Endpoint, "Path under which to expose metrics.")
	//insecure        = flag.Bool("insecure", true, "Ignore server certificate if using https")
)


func main() {

	flag.Parse()
	exporter := newExporter()
	prometheus.MustRegister(exporter)
	prometheus.Unregister(prometheus.NewGoCollector())

	http.Handle(*metricsEndpoint, promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}
