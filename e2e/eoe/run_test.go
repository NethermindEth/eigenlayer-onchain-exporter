package e2e

import (
	"os"
	"os/exec"
	"testing"
	"time"

	base "github.com/NethermindEth/eigenlayer-onchain-exporter/e2e"
	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestE2E_ConfigNotFound(t *testing.T) {
	// Test context
	var (
		cmdErr error
		stdErr string
	)
	// Build test case
	e2eTestCase := newE2eEOETestCase(t,
		//Arrange
		func(t *testing.T, appPath string) error {
			return nil
		},
		//Act
		func(t *testing.T, appPath string) *exec.Cmd {
			cmd, _, stdErrReader := base.RunCommandCMD(t, appPath, "eoe", "run")
			go copyOutputString(t, stdErrReader, &stdErr)
			cmdErr = cmd.Wait()
			return cmd
		},
		//Assert
		func(t *testing.T) {
			// CMD should fail
			assert.Error(t, cmdErr)
			// Check if the error message contains the expected error
			assert.Contains(t, stdErr, "Error: open eoe-config.yml: no such file or directory")
		},
	)
	// Run test case
	e2eTestCase.run()
}

func TestE2E_ConfigFound(t *testing.T) {
	// Test context
	var (
		cmdErr error
	)
	// Build test case
	e2eTestCase := newE2eEOETestCase(t,
		//Arrange
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
		//Act
		func(t *testing.T, appPath string) *exec.Cmd {
			cmd, _, _ := base.RunCommandCMD(t, appPath, "eoe", "run")

			time.Sleep(2 * time.Second)
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				t.Fatalf("Failed to send interrupt signal: %s", err)
			}

			cmdErr = cmd.Wait()
			return cmd
		},
		//Assert
		func(t *testing.T) {
			// CMD should fail
			assert.NoError(t, cmdErr)
		},
	)
	// Run test case
	e2eTestCase.run()
}
