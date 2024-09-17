package main

import (
	"log/slog"
	"os"

	"github.com/NethermindEth/eigenlayer-onchain-exporter/cmd"
)

func main() {
	rootCmd := cmd.RootCmd()
	err := rootCmd.Execute()
	if err != nil {
		slog.Error("error executing root command", "error", err)
		os.Exit(1)
	}
}
