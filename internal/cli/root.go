package cli

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "eoe",
		Short: "EigenLayer On-chain Exporter (eoe) exposes Prometheus metrics about AVS protocols and EigenLayer's Node Operator.",
	}
	rootCmd.PersistentFlags().StringP("config", "c", "eoe-config.yml", "path to config file")
	rootCmd.AddCommand(runCommand())

	return rootCmd
}
