package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	base "github.com/NethermindEth/eigenlayer-onchain-exporter/e2e"
)

type (
	e2eArranger func(t *testing.T, testDir string) error
	e2eAct      func(t *testing.T, testDir string) *exec.Cmd
	e2eAssert   func(t *testing.T)
)

type e2eEOETestCase struct {
	base.E2ETestCase
	arranger e2eArranger
	act      e2eAct
	assert   e2eAssert
	pid      int
}

func newE2eEOETestCase(t *testing.T, arranger e2eArranger, act e2eAct, assert e2eAssert) *e2eEOETestCase {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tc := &e2eEOETestCase{
		E2ETestCase: base.E2ETestCase{
			T:          t,
			TestDir:    t.TempDir(),
			RepoPath:   filepath.Dir(filepath.Dir(wd)),
			BinaryName: "eoe",
		},
		arranger: arranger,
		act:      act,
		assert:   assert,
	}
	t.Logf("Creating new E2E test case (%p). TestDir: %s", tc, tc.TestDir)
	base.CheckGoInstalled(t)
	tc.E2ETestCase.InstallGoModules()
	tc.E2ETestCase.Build()
	return tc
}

func (e *e2eEOETestCase) run() {
	// Cleanup environment before and after test
	if e.arranger != nil {
		err := e.arranger(e.T, e.TestDir)
		if err != nil {
			e.T.Fatalf("error in Arrange step: %v", err)
		}
	}
	if e.act != nil {
		cmd := e.act(e.T, e.TestDir)
		e.pid = cmd.Process.Pid
	}
	if e.assert != nil {
		e.assert(e.T)
	}
}
