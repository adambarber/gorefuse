package gorefuse

import (
	"strings"
	"testing"
)

// At accepts a File ending in a bare colon; the render path substitutes the
// invalid-frame placeholder for the dangling separator, never a panic.
func TestAtAcceptsDanglingColonFile(t *testing.T) {
	assertNoPanic(t, "At panicked on a colon-terminated File", func() {
		At(File("a.go:"), 7, Subject("s"), Why("w"))
	})
}

// A colon-terminated File is refused on every path: a hand-built Refusal with
// one renders the invalid-frame placeholder, never a dangling "a.go::line".
func TestHandBuiltDanglingColonFileRendersPlaceholder(t *testing.T) {
	b := Because[string]{
		Reason: "r",
		Cause:  nil,
		What:   Refusal{File: File("a.go:"), Subject: Subject("s"), Why: Why(""), Line: 7, Sensitive: 0},
	}

	if got := b.Error(); got != "Error: "+invalidFrame+": r: s" {
		t.Fatalf("render = %q, want %q", got, "Error: "+invalidFrame+": r: s")
	}
}

// Stored File equals the rendered frame only after the scrub stage: a File
// with dirty bytes is stored raw (the caller supplied it) and rendered clean.
func TestStoredFileEqualsRenderedModuloScrub(t *testing.T) {
	r := At(File("a\x80.go"), 1, Subject("s"), Why("w"))
	if r.File != File("a\x80.go") {
		t.Fatalf("stored File = %q, want the raw trimmed value", r.File)
	}

	b := Because[string]{Reason: "x", Cause: nil, What: r}
	if got := b.Error(); !strings.Contains(got, "\uFFFD.go:1") {
		t.Fatalf("render = %q, want the scrubbed frame", got)
	}
}

// Hand-built Refusals with unknown Sensitive bits fail closed on every path:
// the frame becomes invalidFrame, marked slots are redacted, and the cause is
// suppressed even when no File was supplied.
func TestHandBuiltUnknownBitsFailClosed(t *testing.T) {
	b := Because[string]{
		Reason: "r",
		Cause:  confError("cause-secret"),
		What:   Refusal{File: File("f"), Subject: Subject("subject-secret"), Why: Why("why-secret"), Line: 1, Sensitive: 0xF8},
	}

	want := "Error: " + invalidFrame + ": r: " + Redacted + ". " + Redacted
	if got := b.Error(); got != want {
		t.Fatalf("render = %q, want %q", got, want)
	}
}

func TestHandBuiltUnknownBitsWithoutFileFailClosed(t *testing.T) {
	b := Because[string]{
		Reason: "r",
		Cause:  confError("cause-secret"),
		What:   Refusal{File: "", Subject: "subject-secret", Why: "why-secret", Line: 0, Sensitive: 0x80},
	}

	want := "Error: " + invalidFrame + ": r: " + Redacted + ". " + Redacted
	if got := b.Error(); got != want {
		t.Fatalf("render = %q, want %q", got, want)
	}
}
