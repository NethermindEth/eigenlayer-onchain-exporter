package e2e

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/config"
	"gopkg.in/yaml.v2"
)

func copyOutputString(t *testing.T, reader io.ReadCloser, out *string) {
	t.Helper()
	outBuffer := new(bytes.Buffer)
	_, err := io.Copy(outBuffer, reader)
	if err != nil {
		panic("Failed to copy output: " + err.Error())
	}
	*out = outBuffer.String()
}

func setupConfigFile(t *testing.T, appPath string, config *config.Config) {
	t.Helper()
	t.Logf("appPath: %s", appPath)
	configFile, err := os.Create(filepath.Join(appPath, "eoe-config.yml"))
	if err != nil {
		t.Fatalf("Failed to create config file: %s", err)
	}
	defer configFile.Close()

	encoder := yaml.NewEncoder(configFile)
	if err := encoder.Encode(config); err != nil {
		t.Fatalf("Failed to write config to file: %s", err)
	}
}
