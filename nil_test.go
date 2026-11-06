package gorefuse

import "testing"

// plainDescriber implements Describable on the pointer, so a nil receiver's
// Describe would panic — the hostile nil-able shape From must harden against.
// The nil-able-kind switch covers pointer/map/slice/func/chan for all kinds.
type plainDescriber struct{}

func (d *plainDescriber) Describe() (Refusal, bool) {
	_ = *d // nil receiver panics here; the hostile shape From must not reach

	return Refusal{File: "", Line: 0, Subject: "", Why: ""}, true
}

func TestFromNilChanDescriber(t *testing.T) {	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("From panicked on a nil-chan Describable: %v", r)
		}
	}()

	var d *plainDescriber

	b := From(testRefused, d)
	if got := b.Error(); got != testRendered {
		t.Fatalf("Error() = %q, want %q", got, testRendered)
	}
}

func TestAtRefusesNegativeLine(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("At accepted a negative Line; want a refusal")
		}
	}()

	At(testFile, -2, testSubject, "")
}
