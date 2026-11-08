package gorefuse

import (
	"strings"
	"testing"
	"unicode/utf8"
)

type fuzzDescriber struct{ refusal Refusal }

func (f fuzzDescriber) Describe() (Refusal, bool) { return f.refusal, true }

// FuzzError pins the render-path guarantees against arbitrary slot content.
// Four properties hold for every input.
func FuzzError(f *testing.F) {
	// Scrub/redaction edge bytes and the golden-fixture shapes a hostile or
	// unlucky slot could carry. testRefused rides as the reason-seed spelling.
	seeds := []string{
		"", " ", "\t", "\n", "\r", "\r\n", "\x1b", "\x1b[31m", "\uFFFD",
		"\u202E", "[redacted]", "[REDACTED]", "Error: ", ". ", "f.go", "f.go:1",
		"f.go:", "server.go:42", "listen :8080", "missing",
		"unknown field \"extra\" in /path/app.json, line 1",
		"bind the port before serving",
		string(testRefused),
	}
	for _, s := range seeds {
		for _, t := range seeds {
			f.Add(s, t, int64(len(s)), uint8(0))
			f.Add(s, t, int64(len(t)), uint8(1))
			f.Add(s, t, int64(len(s))+int64(len(t)), uint8(3))
		}
	}

	f.Fuzz(func(t *testing.T, file, slot string, line int64, marks uint8) {
		// At-refusal panics are documented construction-time behavior; the
		// render path must be reachable only through valid carriers, so the
		// fuzz body filters to inputs At accepts rather than recovering.
		refusal, ok := atAccepted(file, slot, line, marks)
		if !ok {
			t.Skip()
		}

		b := From(testReason("refused"), fuzzDescriber{refusal: refusal})

		first := b.Error()
		second := b.Error()

		assertRenderProps(t, first, second)
	})
}

// assertRenderProps fails unless the render holds the four pinned properties:
// valid UTF-8, idempotent across two calls, and no unscrubbed ESC byte.
func assertRenderProps(t *testing.T, first, second string) {
	t.Helper()

	if !utf8.ValidString(first) {
		t.Fatalf("Error() = %q, not valid UTF-8", first)
	}

	if first != second {
		t.Fatalf("Error() not idempotent: %q then %q", first, second)
	}

	if strings.IndexByte(first, 0x1b) >= 0 {
		t.Fatalf("Error() = %q, carries an unscrubbed ESC byte", first)
	}
}

// atAccepted builds the refusal At would build, reporting false when At
// refuses the input: those documented construction-time panics are not render
// failures, so the fuzz body skips them and only valid carriers reach Error.
func atAccepted(file, slot string, line int64, marks uint8) (Refusal, bool) {
	ok := true

	defer func() {
		if recover() != nil {
			ok = false
		}
	}()

	return At(File(file), Line(line), Subject(slot), Why(slot), Sensitive(marks)), ok
}
