package gorefuse

import (
	"errors"
	"strings"
	"testing"
)

var errFirstCause = errors.New("first cause")

var errSecondCause = errors.New("second cause")

func TestRefuseWrapAndJoin(t *testing.T) {
	what := At(File("f.go"), 1, Subject("bad port"), Why("use 8080."))

	plain := Refuse(testRefused, what)
	if plain.Cause != nil || plain.What.File != "f.go" {
		t.Fatalf("Refuse = %+v", plain)
	}

	wrapped := Wrap(testRefused, what, errTestSentinel)
	if !errors.Is(wrapped, errTestSentinel) {
		t.Fatalf("Wrap lost its cause: %v", wrapped)
	}

	joined := Join(testRefused, what, errFirstCause, errSecondCause)

	if !errors.Is(joined, errFirstCause) || !errors.Is(joined, errSecondCause) {
		t.Fatalf("Join lost multi-cause identity: %v", joined)
	}

	text := joined.Error()
	if !strings.Contains(text, "first cause") || !strings.Contains(text, "second cause") {
		t.Fatalf("Join render = %q, want both causes", text)
	}
}
