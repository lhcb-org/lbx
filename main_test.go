package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestInit(t *testing.T) {
	err := os.Setenv("LHCBPROJECTPATH", "testdata/projects")
	if err != nil {
		t.Fatalf("error setting LHCBPROJECTPATH: %v\n", err)
	}

	_ = os.RemoveAll("testdata/test-init")

	cmd := exec.Command("lbx", "init", "-lvl=-2", "-user-area=testdata/test-init", "gaudi")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Fatalf("error running lbx-init: %v\n", err)
	}

	_ = os.RemoveAll("testdata/test-init")
}
