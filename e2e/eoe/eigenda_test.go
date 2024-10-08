package e2e

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	base "github.com/NethermindEth/eigenlayer-onchain-exporter/e2e"
	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestEigenDAExporterUpMetric(t *testing.T) {
	var (
		cmdErr         error
		prometheusResp *http.Response
		prometheusErr  error
	)

	e2eTestCase := newE2eEOETestCase(t,
		// Arrange
		func(t *testing.T, appPath string) error {
			setupConfigFile(t, appPath, &config.Config{
				Operators: []config.OperatorConfig{
					{
						Name:         "nethermind",
						Address:      "0x57b6FdEF3A23B81547df68F44e5524b987755c99",
						BLSPublicKey: [2]string{"8888183187486914528692107799849671390221086122063975348075796070706039667533", "1162660161480410110225128994312394399428655142287492115882227161635275660953"},
						AVSEnvs:      []string{"eigenda-holesky"},
						EigenDAConfig: config.EigenDAConfig{
							Quorums: map[int]bool{
								0: true,
							},
						},
					},
				},
				RPCs: map[string]string{
					"holesky": "https://ethereum-holesky-rpc.publicnode.com",
				},
				LogLevel: "debug",
			})
			return nil
		},
		// Act
		func(t *testing.T, appPath string) *exec.Cmd {
			cmd, _, _ := base.RunCommandCMD(t, appPath, "eoe", "run")

			// Wait for the Prometheus server to start and metrics to be collected
			time.Sleep(10 * time.Second)

			// Get Prometheus metrics
			prometheusResp, prometheusErr = http.Get("http://localhost:9090/metrics")

			// Stop EOE
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				t.Fatalf("Failed to send interrupt signal: %s", err)
			}

			// Wait for EOE to exit
			cmdErr = cmd.Wait()
			return cmd
		},
		// Assert
		func(t *testing.T) {
			// Check command ran successfully
			assert.NoError(t, cmdErr)
			
			// Check Prometheus metrics were collected successfully
			assert.NoError(t, prometheusErr)
			assert.Equal(t, http.StatusOK, prometheusResp.StatusCode)

			// Check for the presence of eigenda_exporter_up metrics
			responseBody, err := io.ReadAll(prometheusResp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %s", err)
			}
			assert.Contains(t, string(responseBody), "eigenda_exporter_up")
		},
	)

	e2eTestCase.run()
}

func TestEigenDAExporterLatestBlockMetric(t *testing.T) {
	var (
		cmdErr         error
		prometheusResp *http.Response
		prometheusErr  error
	)

	e2eTestCase := newE2eEOETestCase(t,
		// Arrange
		func(t *testing.T, appPath string) error {
			setupConfigFile(t, appPath, &config.Config{
				Operators: []config.OperatorConfig{
					{
						Name:         "nethermind",
						Address:      "0x57b6FdEF3A23B81547df68F44e5524b987755c99",
						BLSPublicKey: [2]string{"8888183187486914528692107799849671390221086122063975348075796070706039667533", "1162660161480410110225128994312394399428655142287492115882227161635275660953"},
						AVSEnvs:      []string{"eigenda-holesky"},
						EigenDAConfig: config.EigenDAConfig{
							Quorums: map[int]bool{
								0: true,
							},
						},
					},
				},
				RPCs: map[string]string{
					"holesky": "https://ethereum-holesky-rpc.publicnode.com",
				},
				LogLevel: "debug",
			})
			return nil
		},
		// Act
		func(t *testing.T, appPath string) *exec.Cmd {
			cmd, _, _ := base.RunCommandCMD(t, appPath, "eoe", "run")

			// Wait for the Prometheus server to start and metrics to be collected
			time.Sleep(time.Minute)

			// Get Prometheus metrics
			prometheusResp, prometheusErr = http.Get("http://localhost:9090/metrics")

			// Stop EOE
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				t.Fatalf("Failed to send interrupt signal: %s", err)
			}

			// Wait for EOE to exit
			cmdErr = cmd.Wait()
			return cmd
		},
		// Assert
		func(t *testing.T) {
			// Check command ran successfully
			assert.NoError(t, cmdErr)
			
			// Check Prometheus metrics were collected successfully
			assert.NoError(t, prometheusErr)
			assert.Equal(t, http.StatusOK, prometheusResp.StatusCode)

			// Check for the presence of eigenda_exporter_latest_block metric
			responseBody, err := io.ReadAll(prometheusResp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %s", err)
			}
			assert.Contains(t, string(responseBody), "eoe_eigenda_exporter_latest_block")
		},
	)

	e2eTestCase.run()
}
