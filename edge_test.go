package gorefuse

import "testing"

const (
	testFix      = "Fix it."
	testRendered = "refused"
)

func TestEmptyFrameWithWhyDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Error() panicked on an empty message with a Why: %v", r)
		}
	}()

	b := Because[testReason]{What: At("", 0, "", Why(testFix)), Reason: "", Cause: nil}
	if got := b.Error(); got != testFix {
		t.Fatalf("Error() = %q, want %q", got, testFix)
	}
}

func TestFromNilDescribable(t *testing.T) {
	b := From(testRefused, nil)
	if got := b.Error(); got != testRendered {
		t.Fatalf("Error() = %q, want %q", got, testRendered)
	}
}

func TestEmbeddedNewlinePassesThrough(t *testing.T) {
	b := Because[testReason]{
		What:   Refusal{File: "f.toml", Line: 3, Subject: "", Why: ""},
		Reason: testReason("line one\nline two"),
		Cause:  nil,
	}
	got := b.Error()

	want := "f.toml:3: line one\nline two"
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}
