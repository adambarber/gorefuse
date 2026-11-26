package gorefuse

import (
	"strings"
	"testing"
)

// Test fixtures shared across the render tests.
const (
	testFix      = "Fix it."
	testRendered = "refused"

	// wantFlattenedFrame pins that a line break cannot create a second
	// line-start verdict: the reason stays on one rendered line.
	wantFlattenedFrame = "f.toml:3: line one line two"
)

// testContext is a named kind for the situation a panic assertion describes.
type testContext string

// renderText is a named kind for expected and rendered refusal messages.
type renderText string

// assertRender pins both surfaces: Verdict owns the label, Error omits it.
// For follow-up checks call b.Error() or b.Verdict() again.
func assertRender[R any](t *testing.T, b Because[R], want renderText) {
	t.Helper()

	if got := b.Verdict(); got != string(want) {
		t.Fatalf("Verdict() = %q, want %q", got, want)
	}

	unlabelled := strings.TrimPrefix(string(want), "Error: ")
	if got := b.Error(); got != unlabelled {
		t.Fatalf("Error() = %q, want %q", got, unlabelled)
	}
}

// assertNoPanic fails if call panics; context names the operation under test.
func assertNoPanic(t *testing.T, context testContext, call func()) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("%s: %v", context, r)
		}
	}()

	call()
}

// assertRefusalPanics fails unless call panics with a string that echoes each
// value in wantEcho; accepted describes the input the constructor must refuse.
func assertRefusalPanics(t *testing.T, accepted testContext, call func(), wantEcho ...string) {
	t.Helper()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("%s; want a refusal", accepted)
		}

		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic is %T, want string", r)
		}

		for _, want := range wantEcho {
			if !strings.Contains(msg, want) {
				t.Fatalf("panic %q does not echo %q", msg, want)
			}
		}
	}()

	call()
}
