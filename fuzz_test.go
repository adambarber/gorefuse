package gorefuse

import (
	"strings"
	"testing"
	"unicode/utf8"
)

type fuzzDescriber struct{ refusal Refusal }

func (f fuzzDescriber) Describe() (Refusal, bool) { return f.refusal, true }

// FuzzError pins the render-path guarantees against arbitrary slot content.
// Every carrier is built from the fuzz inputs; no input is filtered out.
func FuzzError(f *testing.F) {
	addFuzzSeeds(f)

	f.Fuzz(func(t *testing.T, file, reason, slot string, line int64, marks uint8) {
		assertFuzzRender(t, file, reason, slot, line, marks)
	})
}

func addFuzzSeeds(f *testing.F) {
	f.Helper()

	seeds := []string{
		"", " ", "\t", "\n", "\r", "\r\n", "\x1b", "\x1b[31m", "\uFFFD",
		"\u2028", "\u2029", "\x00", "[redacted]", "[REDACTED]", "Error: ",
		"f.go", "f.go:1", "f.go:", "f.go:99", "f.go\nx.go", "server.go:42",
		"listen :8080", "missing", "Error: forged",
		"unknown field \"extra\" in /path/app.json, line 1",
		"bind the port before serving", string(testRefused),
	}
	for _, s := range seeds {
		for _, t := range seeds {
			f.Add(s, t, s, int64(len(s)), uint8(0))
			f.Add(s, t, s, int64(len(t)), uint8(1))
			f.Add(s, t, s, int64(len(s))+int64(len(t)), uint8(3))
		}
	}
}

func assertFuzzRender(t *testing.T, file, reason, slot string, line int64, marks uint8) {
	t.Helper()

	for i, carrier := range fuzzCarriers(file, reason, slot, line, marks) {
		first := carrier.Error()
		second := carrier.Error()
		assertRenderProps(t, first, second)

		// At masks unknown bits before construction, so only the hand-built
		// and From carriers exercise the fail-closed unknown-bit path.
		if i < 2 {
			assertSensitiveProps(t, first, slot, marks)
		}
	}
}

func fuzzCarriers(file, reason, slot string, line int64, marks uint8) []error {
	what := Refusal{
		File: File(file), Line: Line(line), Subject: Subject(slot),
		Why: Why(slot), Sensitive: Sensitive(marks),
	}

	carriers := []error{
		Because[string]{Reason: reason, What: what, Cause: confError(slot)},
		From(testReason(reason), fuzzDescriber{refusal: what}),
	}

	// At rejects only a whitespace-only File paired with a positive Line;
	// that is a documented construction rule, not a render property. Every
	// other file name, including ambiguous ones, goes through At here.
	if strings.TrimSpace(file) != "" && line >= 0 {
		carriers = append(carriers, Because[testReason]{
			Reason: testRefused, Cause: nil,
			What: At(File(file), Line(line), Subject(slot), Why(slot), Sensitive(marks)&sensitiveMask),
		})
	}

	return carriers
}

// assertRenderProps fails unless the render holds the pinned properties:
// valid UTF-8, idempotent, no unscrubbed ESC byte, and exactly one verdict on
// a single line.
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

	if strings.HasSuffix(first, " ") {
		t.Fatalf("Error() = %q, has a dangling connector or trailing space", first)
	}

	assertOneVerdict(t, renderText(first))
}

// assertSensitiveProps fails when an unknown Sensitive bit lets a populated
// slot render in the clear on a hand-built or From carrier.
func assertSensitiveProps(t *testing.T, render, slot string, marks uint8) {
	t.Helper()

	if Sensitive(marks)&^sensitiveMask == 0 || strings.TrimSpace(scrub(slot)) == "" {
		return
	}

	if !strings.Contains(render, Redacted) {
		t.Fatalf("Error() = %q, unknown Sensitive bits did not redact slot %q", render, slot)
	}
}
