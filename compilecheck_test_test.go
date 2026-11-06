package gorefuse

import (
	"bytes"
	"context"
	"os/exec"
	"testing"
)

// The compile-time wall: Because[R1] and Because[R2] are distinct types, so
// the assignment in compilecheck_test.go must fail to build. The test runs
// the build with the file's tag on and expects failure.
func TestCrossDomainAssignmentDoesNotCompile(t *testing.T) {
	cmd := exec.CommandContext(context.Background(), "go", "test", "-run=^$", "-tags=compilecheck", "./...")

	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("cross-domain Because assignment compiled; the compile-time wall is gone")
	}

	if !bytes.Contains(out, []byte("cannot use a")) {
		t.Fatalf("build failed for the wrong reason: %s", out)
	}
}
