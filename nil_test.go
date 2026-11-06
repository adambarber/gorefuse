package gorefuse

import (
	"strings"
	"testing"
)

func TestAtRefusesNegativeLine(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("At accepted a negative Line; want a refusal")
		}

		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic is %T, want string", r)
		}

		for _, want := range []string{"-2", string(testFile)} {
			if !strings.Contains(msg, want) {
				t.Fatalf("panic %q does not echo %q", msg, want)
			}
		}
	}()

	At(testFile, -2, testSubject, "")
}

// plainDescriber implements Describable on the pointer, so a nil receiver's
// Describe would panic — the hostile nil-able shape From must harden against.
// Each nil-able-kind arm (pointer/map/slice/func/chan) of isNilDescribable
// must resolve to the reason alone without calling Describe.
type plainDescriber struct{}

func (d *plainDescriber) Describe() (Refusal, bool) {
	_ = *d // nil receiver panics here; the hostile shape From must not reach

	return Refusal{File: "", Line: 0, Subject: "", Why: ""}, true
}

// Each kind below implements Describable on a nil-able type, so a nil value in
// the interface is one of the hostile nil-able shapes isNilDescribable must
// detect before calling Describe.
type mapDescriber map[string]*plainDescriber

func (d mapDescriber) Describe() (Refusal, bool) {
	return d["x"].Describe()
}

type sliceDescriber []*plainDescriber

func (d sliceDescriber) Describe() (Refusal, bool) {
	return d[0].Describe()
}

type funcDescriber func() (Refusal, bool)

func (d funcDescriber) Describe() (Refusal, bool) {
	return d()
}

type chanDescriber chan *plainDescriber

func (d chanDescriber) Describe() (Refusal, bool) {
	return (<-d).Describe()
}

func TestFromNilDescribableKinds(t *testing.T) {
	kinds := map[string]Describable{
		"pointer": (*plainDescriber)(nil),
		"map":     mapDescriber(nil),
		"slice":   sliceDescriber(nil),
		"func":    funcDescriber(nil),
		"chan":    chanDescriber(nil),
	}

	for name, d := range kinds {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("From panicked on a nil-%s Describable: %v", name, r)
				}
			}()

			b := From(testRefused, d)
			if got := b.Error(); got != testRendered {
				t.Fatalf("Error() = %q, want %q", got, testRendered)
			}
		})
	}
}
