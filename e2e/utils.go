package e2e

import (
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func RunCommand(t *testing.T, path string, binaryName string, args ...string) error {
	_, _, err := runCommandOutput(t, path, binaryName, args...)
	return err
}

func RunCommandCMD(t *testing.T, path string, binaryName string, args ...string) (cmd *exec.Cmd, stdOut io.ReadCloser, stdErr io.ReadCloser) {
	t.Helper()
	t.Logf("Running command: %s %s", binaryName, strings.Join(args, " "))
	cmd = exec.Command(filepath.Join(path, binaryName), args...)

	cmd.Dir = path

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to pipe stdout: %s", err)
	}
	stdErr, err = cmd.StderrPipe()
	if err != nil {
		t.Fatalf("Failed to pipe stderr: %s", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start command: %s %s", binaryName, strings.Join(args, " "))
	}

	return cmd, stdOut, stdErr
}

func runCommandOutput(t *testing.T, path string, binaryName string, args ...string) ([]byte, *exec.Cmd, error) {
	t.Helper()
	t.Logf("Binary path: %s", path)
	t.Logf("Running command: %s %s", binaryName, strings.Join(args, " "))
	cmd := exec.Command(path, args...)
	out, err := cmd.CombinedOutput()
	t.Logf("===== OUTPUT =====\n%s\n==================", out)
	return out, cmd, err
}

func LogAndPipeError(t *testing.T, prefix string, err error) error {
	t.Helper()
	if err != nil {
		t.Log(prefix, err)
	}
	return err
}
