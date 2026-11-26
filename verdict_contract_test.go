package gorefuse

import (
	"fmt"
	"strings"
	"testing"
)

func TestErrorDoesNotEmbedVerdictLabel(t *testing.T) {
	b := because(At("server.go", 42, "listen :8080", "bind the port"), testRefused, nil)
	want := "server.go:42: refused: listen :8080. bind the port"

	if got := b.Error(); got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}

	if strings.Contains(b.Error(), "Error: ") {
		t.Fatalf("Error() embedded a verdict label: %q", b.Error())
	}
}

func TestVerdictOwnsSingleLabel(t *testing.T) {
	b := because(At("server.go", 42, "listen :8080", "bind the port"), testRefused, nil)
	want := "Error: server.go:42: refused: listen :8080. bind the port"

	if got := b.Verdict(); got != want {
		t.Fatalf("Verdict() = %q, want %q", got, want)
	}
}

func TestWrappedErrorStaysComposable(t *testing.T) {
	b := because(At("server.go", 42, "listen :8080", "bind the port"), testRefused, nil)
	wrapped := fmt.Errorf("startup: %w", b)
	want := "startup: server.go:42: refused: listen :8080. bind the port"

	if got := wrapped.Error(); got != want {
		t.Fatalf("wrapped Error() = %q, want %q", got, want)
	}
}
