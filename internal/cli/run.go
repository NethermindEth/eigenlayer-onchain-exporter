package cli

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/avs/eigenda"
	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/avsexporter"
	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/config"
	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/prometheus"
	"github.com/spf13/cobra"
)

type exporterError struct {
	exporter avsexporter.AVSExporter
	err      error
}

func runCommand() *cobra.Command {
	var c *config.Config
	return &cobra.Command{
		Use:   "run",
		Short: "Run the application",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}
			c, err = config.GetConfig(configPath)
			if err != nil {
				return err
			}
			logLevel := slog.Level(slog.LevelInfo)
			if err := logLevel.UnmarshalText([]byte(c.LogLevel)); err != nil {
				return err
			}
			slog.SetLogLoggerLevel(logLevel)
			go func() {
				err := prometheus.StartPrometheusServer(":9090")
				if err != nil {
					slog.Error("Error starting Prometheus server", "error", err)
				}
			}()
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				wg              sync.WaitGroup
				ctx             = cmd.Context()
				avsEnvs         = make(map[string]bool)
				exporterErrorCh = make(chan exporterError, len(avsEnvs))
			)

			// Add all AVS environments from operators
			for _, operator := range c.Operators {
				for _, env := range operator.AVSEnvs {
					avsEnvs[env] = true
				}
			}

			// Run exporters for each AVS environment
			for env := range avsEnvs {
				switch env {
				case config.AVSEnvEigenDAMainnet, config.AVSEnvEigenDAHolesky:
					// Initialize and run the AVS environment exporter
					exporter, err := eigenda.NewEigenDAOnChainExporter(env, c)
					if err != nil {
						return err
					}
					runExporter(ctx, exporter, &wg, exporterErrorCh, c)
				default:
					return fmt.Errorf("invalid AVS environment: %s", env)
				}
			}

			for {
				select {
				case exporterError := <-exporterErrorCh:
					slog.Debug("exporter error", "exporter", exporterError.exporter.Name(), "error", exporterError.err)
					runExporter(ctx, exporterError.exporter, &wg, exporterErrorCh, c)
				case <-ctx.Done():
					slog.Debug("context done", "error", ctx.Err())
					return gracefulExit(&wg, nil)
				}
			}
		},
	}
}

// runExporter starts an exporter and adds it to the wait group. It also sends
// any errors to the exporterErrorCh channel.
func runExporter(
	ctx context.Context,
	exporter avsexporter.AVSExporter,
	wg *sync.WaitGroup,
	exporterErrorCh chan<- exporterError,
	c *config.Config,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := exporter.Run(ctx, c)
		if err != nil {
			if ctx.Err() == nil {
				slog.Error("exporter error", "exporter", exporter.Name(), "error", err)
				exporterErrorCh <- exporterError{exporter, err}
			}
		}
	}()
}

func gracefulExit(wg *sync.WaitGroup, err error) error {
	slog.Debug("Shutting down exporters...")
	wg.Wait()
	slog.Debug("Exporters shutdown complete")
	return err
}
