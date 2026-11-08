package gorefuse

import (
	"strconv"
	"testing"
)

type testDescriber struct {
	what   Refusal
	misses bool
}

func (d testDescriber) Describe() (Refusal, bool) {
	return d.what, !d.misses
}

func TestAtBuildsRefusal(t *testing.T) {
	got := At(testFile, 7, testSubject, Why("Fix it."))

	want := Refusal{File: testFile, Subject: testSubject, Why: "Fix it.", Line: 7, Sensitive: 0}
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
	b := From(testRefused, testDescriber{what: Refusal{File: testFile, Line: 0, Subject: testSubject, Why: "", Sensitive: 0}, misses: false})

	assertRender(t, b, wantFrameSubject)
}

func TestFromWithoutDescription(t *testing.T) {
	b := From(testRefused, testDescriber{what: Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, misses: true})

	assertRender(t, b, "Error: "+testRendered)
}

func TestAtRefusesWhitespaceFile(t *testing.T) {
	assertRefusalPanics(t, "At accepted a whitespace-only File", func() {
		At(" ", 3, testSubject, "")
	}, strconv.Quote(" "))
}

func TestAtRefusesTrailingDigitsFile(t *testing.T) {
	assertRefusalPanics(t, `At accepted File "real.go:99"`, func() {
		At("real.go:99", 3, testSubject, "")
	}, strconv.Quote("real.go:99"))
}

func TestAtRefusesNewlineFile(t *testing.T) {
	assertRefusalPanics(t, `At accepted File "a\nb"`, func() {
		At("a\nb", 3, testSubject, "")
	}, strconv.Quote("a\nb"))
}

func TestAtRefusesWhitespaceFileWithLine(t *testing.T) {
	assertRefusalPanics(t, "At accepted a whitespace-only File with a Line", func() {
		At("  ", 2, testSubject, "")
	}, strconv.Itoa(2))
}

func TestAtRefusesUnknownSensitiveBits(t *testing.T) {
	assertRefusalPanics(t, "At accepted unknown Sensitive bits", func() {
		At(testFile, 0, testSubject, "", SubjectSensitive, Sensitive(0x40))
	}, strconv.Itoa(0x40))
}

func TestAtStoresTrimmedFile(t *testing.T) {
	got := At("  spaced.go  ", 2, testSubject, "")

	if got.File != File("spaced.go") {
		t.Fatalf("At().File = %q, want %q", got.File, "spaced.go")
	}

	assertRender(t, From(testRefused, testDescriber{what: got, misses: false}), renderText("Error: spaced.go:2: refused: "+string(testSubject)))
}

func TestMarkedCauseSuppressesCauseTail(t *testing.T) {
	b := because(Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: CauseSensitive}, testRefused, errTestSentinel)

	assertRender(t, b, "Error: "+testRendered)
}

func TestFromForgedNewlineFileRendersPlaceholder(t *testing.T) {
	// Hostile probe: a hand-built Refusal whose File carries a newline must
	// not forge a file:line pairing through the From path.
	b := From(testRefused, testDescriber{what: Refusal{
		File: "evil\n10.0.0.1", Line: 5, Subject: "", Why: "", Sensitive: 0,
	}, misses: false})

	assertRender(t, b, "Error: "+invalidFrame+": refused")
}

func TestFromForgedTrailingDigitsFileRendersPlaceholder(t *testing.T) {
	b := From(testRefused, testDescriber{what: Refusal{
		File: "fake:123", Line: 4, Subject: "", Why: "", Sensitive: 0,
	}, misses: false})

	assertRender(t, b, "Error: "+invalidFrame+": refused")
}

func TestFromForgedTrailingDigitsFileZeroLineStillPlaceholder(t *testing.T) {
	// The guard is total: "fake:123" alone already looks like a file:line
	// pairing to a log parser, so Line 0 does not excuse it.
	b := From(testRefused, testDescriber{what: Refusal{
		File: "fake:123", Line: 0, Subject: "", Why: "", Sensitive: 0,
	}, misses: false})

	assertRender(t, b, "Error: "+invalidFrame+": refused")
}

func TestFrameSkipsNegativeLine(t *testing.T) {
	// A negative Line can never pair as file:line, but no pairing at all is
	// the only shape the parser contract allows; the line is skipped.
	b := because(Refusal{File: testFile, Line: -3, Subject: "", Why: "", Sensitive: 0}, testRefused, nil)

	assertRender(t, b, renderText("Error: "+string(testFile)+": refused"))
}

func TestFrameAcceptsInteriorColonFile(t *testing.T) {
	assertRender(t, because(At("a:b.go", 3, "", Why(testFix)), testRefused, nil),
		"Error: a:b.go:3: refused. "+testFix)
}
