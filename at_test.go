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

	assertRender(t, b, wantFrameSubject)
}

func TestFromWithoutDescription(t *testing.T) {
	b := From(testRefused, testDescriber{what: Refusal{File: "", Line: 0, Subject: "", Why: ""}, misses: true})

	assertRender(t, b, testRendered)
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

func TestFrameAcceptsInteriorColonFile(t *testing.T) {
	assertRender(t, because(At("a:b.go", 3, "", Why(testFix)), testRefused, nil),
		"a:b.go:3: refused. "+testFix)
}
