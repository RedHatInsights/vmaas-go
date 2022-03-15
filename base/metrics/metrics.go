package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	EngineVersion = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Help:      "Patchman project deployment information",
		Namespace: "patchman_engine",
		Subsystem: "core",
		Name:      "info",
	}, []string{"version"})

	// ENGINEVERSION - DO NOT EDIT this variable MANUALLY - it is modified by generate_docs.sh
	ENGINEVERSION = "v0.0.2"
)

func init() {
	prometheus.MustRegister(EngineVersion)
	EngineVersion.WithLabelValues(ENGINEVERSION).Set(1)
}
