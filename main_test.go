package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {

	for _, k := range []string{"CMAKE_PREFIX_PATH", "CMTPROJECTPATH", "LHCBPROJECTPATH"} {
		os.Setenv(k, "")
	}

	const testinit = "testdata/test-init"

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = os.Setenv("LHCBPROJECTPATH", filepath.Join(pwd, "testdata/projects"))
	if err != nil {
		t.Fatalf("error setting LHCBPROJECTPATH: %v\n", err)
	}

	_ = os.RemoveAll(testinit)

	err = os.MkdirAll(testinit, 0755)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = os.Chdir(testinit)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	cmd := exec.Command("lbx", "init", "-lvl=-2", "gaudi")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Fatalf("error running lbx-init: %v\n", err)
	}

	_ = os.Chdir(pwd)
	_ = os.RemoveAll(testinit)
}
