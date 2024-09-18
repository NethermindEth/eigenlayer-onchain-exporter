package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/avs/eigenda"
	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/config"
	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/prometheus"
	"github.com/spf13/cobra"
)

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
			go prometheus.StartPrometheusServer(":9090")
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				wg          sync.WaitGroup
				ctx, cancel = context.WithCancel(cmd.Context())
				avsEnvs     = make(map[string]bool)
				errChan     = make(chan error, len(avsEnvs))
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
						slog.Error(err.Error())
						cancel()
						wg.Wait()
						return err
					}
					wg.Add(1)
					go func() {
						defer wg.Done()
						errChan <- exporter.Run(ctx, c)
					}()
				default:
					gracefulExit(cancel, &wg, fmt.Errorf("invalid AVS environment: %s", env))
				}
			}
			// Wait for all exporters to finish
			for {
				select {
				case err := <-errChan:
					gracefulExit(cancel, &wg, err)
				case <-ctx.Done():
					gracefulExit(cancel, &wg, nil)
				}
			}
		},
	}
}

func gracefulExit(cancel context.CancelFunc, wg *sync.WaitGroup, err error) {
	slog.Debug("Graceful exit")
	cancel()
	wg.Wait()
	if err != nil {
		slog.Error(err.Error())
	}
	os.Exit(1)
}
