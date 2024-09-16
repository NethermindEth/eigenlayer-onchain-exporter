package avsexporter

import (
	"context"

	"github.com/NethermindEth/eigenda-blob-scrapper/internal/config"
)

type AVSExporter interface {
	Run(context.Context, *config.Config) error
}
