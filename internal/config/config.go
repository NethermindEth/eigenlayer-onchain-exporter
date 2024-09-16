package config

import (
	"log/slog"
)

const (
	// RPCNetwork is the network type for the RPC.
	RPCNetworkEthereum = "ethereum"
	RPCNetworkHolesky  = "holesky"

	// AVSEnv is the environment for the AVS.
	AVSEnvEigenDAHolesky = "eigenda-holesky"
	AVSEnvEigenDAMainnet = "eigenda-mainnet"
)

// Config is the configuration for the application.
type Config struct {
	// Operators is the list of operators to be tracked.
	Operators []OperatorConfig  `yaml:"operators"`
	// RPCs is the list of RPCs to be used for the AVS exporters.
	RPCs      map[string]string `yaml:"rpcs"`
	// LogLevel is the level of logging to be used.
	LogLevel  slog.Level        `yaml:"logLevel"`
}

// OperatorConfig holds the needed information for an operator to be tracked.
type OperatorConfig struct {
	// Name is the name of the operator.
	Name          string        `yaml:"name"`
	// Address is the address of the operator.
	Address       string        `yaml:"address"`
	// BLSPublicKey is the BLS public key of the operator.
	BLSPublicKey  [2]string     `yaml:"blsPublicKey"`
	// AVSEnvs is the list of AVS environments to be tracked.
	AVSEnvs       []string      `yaml:"avsEnvs"`
	// EigenDAConfig is the configuration for the EigenDA AVS.
	EigenDAConfig EigenDAConfig `yaml:"eigenDAConfig"`
}

type EigenDAConfig struct {
	// Quorums is the initial status of the operator's quorums. If the exporter
	// receives events in the future about the operator's quorum status, it will
	// update the status in the Prometheus metric. This map is only for bootstrapping
	// and does not have a missing Prometheus metric.
	Quorums map[int]bool `yaml:"quorums"`
}
