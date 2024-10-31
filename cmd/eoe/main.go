package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/cli"
)

func main() {
	rootCmd := cli.RootCmd()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		cancel()
	}()
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		slog.Error("error executing root command", "error", err)
		os.Exit(1)
	}
}
