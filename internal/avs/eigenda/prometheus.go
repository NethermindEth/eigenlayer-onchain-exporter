package eigenda

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metricExporterLatestBlock = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "eoe",
		Name:      "eigenda_exporter_latest_block",
		Help:      "Latest block number that the exporter has processed",
	}, []string{"network"})
	metricOnchainBatchesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "eoe",
		Name:      "eigenda_onchain_batches_total",
		Help:      "Total number of eigenda onchain batches",
	}, []string{"network"})
	metricOnchainBatches = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "eoe",
		Name:      "eigenda_onchain_batches",
		Help:      "Number of eigenda onchain batches",
	}, []string{"operator", "network", "status"})
	metricOnchainQuorumStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "eoe",
		Name:      "eigenda_onchain_quorum_status",
		Help:      "Quorum status of eigenda onchain",
	}, []string{"operator", "network", "quorum"})
	metricExporterStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "eoe",
		Name:      "eigenda_exporter_up",
		Help:      "Status of the exporter",
	}, []string{"avsEnv"})
)
