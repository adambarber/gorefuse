package gorefuse

import (
	"strings"
	"testing"
)

// Test fixtures shared across the render tests.
const (
	testFix      = "Fix it."
	testRendered = "refused"

	// wantMultilineFrame pins the embedded-newline passthrough the scrub
	// tests assert from both sides: raw rendering and post-scrub rendering.
	wantMultilineFrame = "f.toml:3: line one\nline two"
)

// testContext is a named kind for the situation a panic assertion describes.
type testContext string

// renderText is a named kind for expected and rendered refusal messages.
type renderText string

// because assembles a Because with the given slots — the shape most render
// tests construct. Reason carries the generic parameter, so pass a typed
// reason (e.g. testRefused) to pin the domain.
func because[R any](what Refusal, reason R, cause error) Because[R] {
	return Because[R]{What: what, Reason: reason, Cause: cause}
}

// assertRender fails unless b renders to want. For follow-up checks (raw
// byte, suffix, contains) call b.Error() again — the assertion above pins it
// to want.
func assertRender[R any](t *testing.T, b Because[R], want renderText) {
	t.Helper()

	if got := b.Error(); got != string(want) {
		t.Fatalf("Error() = %q, want %q", got, want)
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
