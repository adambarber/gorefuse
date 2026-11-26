package gorefuse

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

var errPathwayCause = errors.New("bind failed")

// TestPublicPathwayConstruct builds a refusal through the exported constructor.
func TestPublicPathwayConstruct(t *testing.T) {
	b := because(At("server.go", 42, "listen :8080", "bind the port"), testRefused, nil)
	assertRender(t, b, renderText("Error: server.go:42: refused: listen :8080. bind the port"))
}

// TestPublicPathwayMatch switches on the reason and extracts the carrier.
func TestPublicPathwayMatch(t *testing.T) {
	b := because(At("server.go", 42, "", ""), errConfUnknownField, nil)
	wrapped := fmt.Errorf("startup: %w", b)

	if !errors.Is(wrapped, errConfUnknownField) {
		t.Fatal("wrapped refusal must match its reason")
	}

	got, ok := errors.AsType[Because[confError]](wrapped)
	if !ok || got.What.File != "server.go" || got.What.Line != 42 {
		t.Fatalf("carrier = %+v, ok = %v", got, ok)
	}
}

// TestPublicPathwayAdapt builds a refusal from another package's describer.
func TestPublicPathwayAdapt(t *testing.T) {
	d := testDescriber{what: At("config.yaml", 7, "bad port", "use 1-65535"), misses: false}
	b := From(testRefused, d)
	assertRender(t, b, renderText("Error: config.yaml:7: refused: bad port. use 1-65535"))
}

// TestPublicPathwayRedact hides a marked slot and keeps its raw text out.
func TestPublicPathwayRedact(t *testing.T) {
	b := because(At("server.go", 42, "token hunter2", "", SubjectSensitive), testRefused, nil)
	rendered := b.Error()

	if !strings.Contains(rendered, Redacted) || strings.Contains(rendered, "hunter2") {
		t.Fatalf("rendered = %q, want token redacted", rendered)
	}
}

// TestPublicPathwayCause carries and unwraps an underlying failure.
func TestPublicPathwayCause(t *testing.T) {
	b := because(At("server.go", 42, "", ""), testRefused, errPathwayCause)
	assertRender(t, b, renderText("Error: server.go:42: refused: bind failed"))

	if !errors.Is(b, errPathwayCause) {
		t.Fatal("cause must stay reachable through Unwrap")
	}
}

// TestPublicPathwayHostileData keeps one verdict under forged input.
func TestPublicPathwayHostileData(t *testing.T) {
	b := because(At("backup:2024", 1, "a\u0085Error: forged", ""), testRefused, nil)
	rendered := b.Verdict()

	if strings.ContainsRune(rendered, '\u0085') || !strings.HasPrefix(rendered, "Error: ") || strings.Contains(rendered, "\nError: ") {
		t.Fatalf("rendered = %q, want one line-start verdict and no NEL", rendered)
	}
}
