package gorefuse

import "testing"

// nilDescriber implements Describable on the pointer with a value-receiver
// call into the zero value, so a nil pointer reaching Describe panics —
// the hostile shape From must harden against.
type nilDescriber struct {
	name testReason
}

func (d *nilDescriber) Describe() (Refusal, bool) {
	// Reaching here through a nil receiver panics on name — the hostile shape.
	return Refusal{File: "", Line: 0, Subject: "", Why: ""}, d.name == ""
}

func TestFromTypedNilDescriber(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("From panicked on a typed-nil Describable: %v", r)
		}
	}()

	var d *nilDescriber

	b := From(testRefused, d)
	if got := b.Error(); got != testRendered {
		t.Fatalf("Error() = %q, want %q", got, testRendered)
	}
}

func TestFromNilUntypedInterface(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("From panicked on an untyped-nil Describable: %v", r)
		}
	}()

	var d Describable

	b := From(testRefused, d)
	if got := b.Error(); got != testRendered {
		t.Fatalf("Error() = %q, want %q", got, testRendered)
	}
}

func TestAtRefusesLineWithoutFile(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("At accepted Line > 0 with an empty File; want a refusal")
		}
	}()

	At("", 3, testSubject, "")
}
