package gorefuse

import "testing"

type testDescriber struct {
	what   Refusal
	misses bool
}

func (d testDescriber) Describe() (Refusal, bool) {
	return d.what, !d.misses
}

func TestAtBuildsRefusal(t *testing.T) {
	got := At(testFile, 7, testSubject, Why("Fix it."))

	want := Refusal{File: testFile, Subject: testSubject, Why: "Fix it.", Line: 7}
	if got != want {
		t.Fatalf("At() = %+v, want %+v", got, want)
	}
}
func TestAtZeroLineIsWholeFile(t *testing.T) {
	got := At(testFile, 0, testSubject, "")
	if got.Line != 0 {
		t.Fatalf("At() Line = %d, want 0", got.Line)
	}
}

func TestBecauseFromDescribable(t *testing.T) {
	b := From(testRefused, testDescriber{what: Refusal{File: testFile, Line: 0, Subject: testSubject, Why: ""}, misses: false})
	got := b.Error()

	want := "f.toml: refused: server.tls"
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestFromWithoutDescription(t *testing.T) {
	b := From(testRefused, testDescriber{what: Refusal{File: "", Line: 0, Subject: "", Why: ""}, misses: true})
	got := b.Error()

	want := "refused"
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}
