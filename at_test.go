package gorefuse

import (
	"strconv"
	"strings"
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
	got := b.Error()

	want := wantFrameSubject
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestFromWithoutDescription(t *testing.T) {
	b := From(testRefused, testDescriber{what: Refusal{File: "", Line: 0, Subject: "", Why: ""}, misses: true})
	got := b.Error()

	want := testRendered
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestAtRefusesWhitespaceFile(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("At accepted a whitespace-only File; want a refusal")
		}

		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic is %T, want string", r)
		}

		if !strings.Contains(msg, strconv.Quote(" ")) {
			t.Fatalf("panic %q does not echo the offending File value %q", msg, " ")
		}
	}()

	At(" ", 3, testSubject, "")
}

func TestAtRefusesTrailingDigitsFile(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal(`At accepted File "real.go:99"; want a refusal of the :digits frame suffix`)
		}

		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic is %T, want string", r)
		}

		if !strings.Contains(msg, strconv.Quote("real.go:99")) {
			t.Fatalf("panic %q does not echo the offending File value %q", msg, "real.go:99")
		}
	}()

	At("real.go:99", 3, testSubject, "")
}

func TestAtRefusesNewlineFile(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal(`At accepted File "a\nb"; want a refusal of the embedded newline`)
		}

		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic is %T, want string", r)
		}

		if !strings.Contains(msg, strconv.Quote("a\nb")) {
			t.Fatalf("panic %q does not echo the offending File value %q", msg, "a\nb")
		}
	}()

	At("a\nb", 3, testSubject, "")
}

func TestFrameAcceptsInteriorColonFile(t *testing.T) {
	b := Because[testReason]{
		What:   At("a:b.go", 3, "", Why(testFix)),
		Reason: testRefused,
		Cause:  nil,
	}

	want := "a:b.go:3: refused. " + testFix
	if got := b.Error(); got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}
