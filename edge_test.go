package gorefuse

import "testing"

func TestEmptyFrameWithWhyDoesNotPanic(t *testing.T) {
	assertNoPanic(t, "Error() panicked on an empty message with a Why", func() {
		assertRender(t, because(At("", 0, "", Why(testFix)), "", nil), "Error: "+testFix)
	})
}

func TestFromNilDescribable(t *testing.T) {
	assertRender(t, From(testRefused, nil), "Error: "+testRendered)
}

// TestEmbeddedNewlinePassesThrough pins the newline contract once: raw and
// scrubbed rendering take the same code path, so one assertion covers both
// sides of the scrub boundary.
func TestEmbeddedNewlinePassesThrough(t *testing.T) {
	assertRender(t, because(Refusal{File: "f.toml", Line: 3, Subject: "", Why: "", Sensitive: 0}, testReason("line one\nline two"), nil),
		"Error: "+wantMultilineFrame)
}
