package gorefuse

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

type fuzzDescriber struct{ refusal Refusal }

func (f fuzzDescriber) Describe() (Refusal, bool) { return f.refusal, true }

// verdictCarrier is both an error and a labelled renderer.
type verdictCarrier interface {
	error
	Verdict() string
}

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

	for _, carrier := range fuzzCarriers(file, reason, slot, line, marks) {
		first := carrier.Error()
		second := carrier.Error()
		assertContentProps(t, first, second)
		assertOneVerdict(t, renderText(carrier.Verdict()))
	}

	// At masks unknown bits before construction, so the fixed-reason
	// carriers built here are the ones that exercise the fail-closed path.
	assertSensitiveProps(t, slot, marks)
}

func fuzzCarriers(file, reason, slot string, line int64, marks uint8) []verdictCarrier {
	what := Refusal{
		File: File(file), Line: Line(line), Subject: Subject(slot),
		Why: Why(slot), Sensitive: Sensitive(marks),
	}

	carriers := []verdictCarrier{
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

// assertContentProps fails unless Error holds the content properties:
// valid UTF-8, idempotent, no unscrubbed ESC byte, and no trailing connector.
func assertContentProps(t *testing.T, first, second string) {
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
}

const sensitiveMarker = "\x02"

// assertSensitiveProps fails when an unknown Sensitive bit lets a populated
// slot render in the clear on a hand-built or From carrier.
func assertSensitiveProps(t *testing.T, slot string, marks uint8) {
	t.Helper()

	const fixed = "fixed-reason"

	what := Refusal{
		File: File(fixed), Line: Line(1), Subject: Subject(sensitiveMarker + slot + sensitiveMarker),
		Why: Why(fixed), Sensitive: Sensitive(marks),
	}

	carriers := []error{
		Because[string]{Reason: fixed, What: what, Cause: confError(fixed)},
		From(testReason(fixed), fuzzDescriber{refusal: what}),
	}

	for _, carrier := range carriers {
		if failure := sensitiveProbeFailure(slot, marks, carrier); failure != "" {
			t.Fatal(failure)
		}
	}
}

// sensitiveProbeFailure checks one carrier against the fail-closed oracle. A
// fixed reason keeps a fuzz slot that reads [redacted] from satisfying the
// token check by itself; the marker must disappear from the render too.
func sensitiveProbeFailure(slot string, marks uint8, carrier error) string {
	if Sensitive(marks)&^sensitiveMask == 0 || strings.TrimSpace(scrub(slot)) == "" {
		return ""
	}

	render := carrier.Error()
	if !strings.Contains(render, Redacted) {
		return fmt.Sprintf("Error() = %q, unknown Sensitive bits did not redact slot %q", render, slot)
	}

	if strings.Contains(render, sensitiveMarker) {
		return fmt.Sprintf("Error() = %q, raw slot survived redaction: %q", render, slot)
	}

	return ""
}

func TestSensitiveOracleRequiresRawSlotAbsent(t *testing.T) {
	const slot = "raw-secret"

	what := Refusal{
		File: File("f.go"), Line: Line(1),
		Subject: Subject(sensitiveMarker + slot + sensitiveMarker),
		Why:     Why("fixed"), Sensitive: 0,
	}

	// The old oracle accepted this render because the reason alone carried
	// the redaction token, even though the marked slot stayed in the clear.
	leaky := Because[string]{Reason: Redacted, What: what, Cause: confError("fixed")}
	if failure := sensitiveProbeFailure(slot, uint8(1<<7), leaky); failure == "" {
		t.Fatal("oracle accepted a render with an unredacted slot")
	}

	marked := what
	marked.Sensitive = sensitiveMask | Sensitive(1<<7)

	redacted := Because[string]{Reason: "fixed", What: marked, Cause: confError("fixed")}
	if failure := sensitiveProbeFailure(slot, uint8(1<<7), redacted); failure != "" {
		t.Fatalf("oracle rejected a redacted render: %s", failure)
	}
}
