package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartPrometheusServer starts a Prometheus server on the given port.
func StartPrometheusServer(port string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(port, nil)
}
