package e2e

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

type E2ETestCase struct {
	T          *testing.T
	TestDir    string
	RepoPath   string
	BinaryName string
}

func (e *E2ETestCase) BinaryPath() string {
	binaryName := e.BinaryName
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	return filepath.Join(e.TestDir, binaryName)
}

func (e *E2ETestCase) InstallGoModules() {
	e.T.Helper()
	cmd := exec.Command("go", "mod", "download")
	cmd.Dir = e.RepoPath
	e.T.Logf("Installing Go modules in %s", e.RepoPath)
	if err := cmd.Run(); err != nil {
		e.T.Fatalf("error installing Go modules: %v", err)
	} else {
		e.T.Logf("Go modules installed")
	}
}

func (e *E2ETestCase) Build() {
	e.T.Helper()
	e.T.Logf("Building binary to %s", e.BinaryPath())
	err := exec.Command("go", "build", "-o", e.BinaryPath(), filepath.Join(e.RepoPath, "cmd", e.BinaryName, "main.go")).Run()
	if err != nil {
		e.T.Fatalf("error building binary: %v", err)
	} else {
		e.T.Logf("binary built")
	}
}

func CheckGoInstalled(t *testing.T) {
	t.Helper()
	err := exec.Command("go", "version").Run()
	if err != nil {
		t.Fatalf("error checking Go installation: %v", err)
	} else {
		t.Logf("Go installed")
	}
}
