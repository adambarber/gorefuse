package gorefuse

import (
	"strings"
	"testing"
)

// At refuses a File ending in a bare colon: the rendered frame would carry a
// dangling separator ("a.go::7") that no documented frame shape blesses.
func TestAtRefusesDanglingColonFile(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("At accepted a colon-terminated File; want panic")
		}

		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic value is not a string: %v", r)
		}

		if !strings.Contains(msg, "a.go:") {
			t.Fatalf("panic does not echo the offending File: %v", r)
		}
	}()

	At(File("a.go:"), 7, Subject("s"), Why("w"))
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

// Hand-built Refusals with unknown Sensitive bits render the invalid-frame
// placeholder: the render path is the backstop for every construction path,
// so it refuses what At refuses.
func TestHandBuiltUnknownBitsRenderPlaceholder(t *testing.T) {
	b := Because[string]{
		Reason: "r",
		Cause:  nil,
		What:   Refusal{File: File("f"), Subject: Subject("s"), Why: Why(""), Line: 1, Sensitive: 0xF8},
	}

	if got := b.Error(); !strings.Contains(got, invalidFrame) {
		t.Fatalf("render = %q, want the invalid-frame placeholder", got)
	}
}
