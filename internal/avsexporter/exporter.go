package avsexporter

import (
	"context"

	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/config"
)

type AVSExporter interface {
	Run(context.Context, *config.Config) error
}
