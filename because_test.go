package gorefuse

import (
	"errors"
	"fmt"
	"testing"
)

type testReason string

type otherReason string

const (
	testMissing  testReason  = "missing"
	testRefused  testReason  = "refused"
	otherRefused otherReason = "refused otherwise"
)

var errTestSentinel = errors.New("root cause")

const testFile = File("f.toml")

const testSubject = Subject("server.tls")

const wantFrameSubject = "f.toml: refused: server.tls"

func TestReasonAlone(t *testing.T) {
	b := Because[testReason]{What: Refusal{File: "", Line: 0, Subject: "", Why: ""}, Reason: testMissing, Cause: nil}
	got := b.Error()

	want := "missing"
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestFramePrefixWholeFile(t *testing.T) {
	b := Because[testReason]{What: Refusal{File: "config.toml", Line: 0, Subject: "", Why: ""}, Reason: testRefused, Cause: nil}
	got := b.Error()

	want := "config.toml: refused"
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestFramePrefixWithLine(t *testing.T) {
	b := Because[testReason]{What: Refusal{File: "config.toml", Line: 12, Subject: "", Why: ""}, Reason: testRefused, Cause: nil}
	got := b.Error()

	want := "config.toml:12: refused"
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestSubjectAppends(t *testing.T) {
	b := Because[testReason]{
		What:   Refusal{File: testFile, Line: 0, Subject: testSubject, Why: ""},
		Reason: testRefused,
		Cause:  nil,
	}
	got := b.Error()

	want := wantFrameSubject
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestWhyAppendsAsSentence(t *testing.T) {
	b := Because[testReason]{
		What:   Refusal{File: testFile, Line: 0, Subject: testSubject, Why: "Rotate the certificate."},
		Reason: testRefused,
		Cause:  nil,
	}
	got := b.Error()

	want := "f.toml: refused: server.tls. Rotate the certificate."
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestCauseTail(t *testing.T) {
	b := Because[testReason]{What: Refusal{File: "", Line: 0, Subject: "", Why: ""}, Reason: testRefused, Cause: errTestSentinel}
	got := b.Error()

	want := "refused: root cause"
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestFullRenderStableOrder(t *testing.T) {
	b := Because[testReason]{
		What:   Refusal{File: testFile, Line: 3, Subject: testSubject, Why: "Renew it."},
		Reason: testRefused,
		Cause:  errTestSentinel,
	}
	got := b.Error()

	want := "f.toml:3: refused: server.tls: root cause. Renew it."
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestUnwrapToCause(t *testing.T) {
	b := Because[testReason]{What: Refusal{File: "", Line: 0, Subject: "", Why: ""}, Reason: testRefused, Cause: errTestSentinel}
	if !errors.Is(b, errTestSentinel) {
		t.Fatal("errors.Is(Because, sentinel) = false, want true")
	}
}

func TestUnwrapNilCause(t *testing.T) {
	b := Because[testReason]{What: Refusal{File: "", Line: 0, Subject: "", Why: ""}, Reason: testRefused, Cause: nil}
	if b.Unwrap() != nil {
		t.Fatalf("Unwrap() = %v, want nil", b.Unwrap())
	}
}

func TestReasonRendersStringer(t *testing.T) {
	b := Because[stringerReason]{What: Refusal{File: "", Line: 0, Subject: "", Why: ""}, Reason: stringerReason{}, Cause: nil}
	got := b.Error()

	want := "rendered"
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

type stringerReason struct{}

func (stringerReason) String() string { return "rendered" }

func TestDistinctReasonTypesAreDistinctTypes(t *testing.T) {
	a := Because[testReason]{What: Refusal{File: "", Line: 0, Subject: "", Why: ""}, Reason: testRefused, Cause: nil}
	b := Because[otherReason]{What: Refusal{File: "", Line: 0, Subject: "", Why: ""}, Reason: otherRefused, Cause: nil}

	ta, tb := fmt.Sprintf("%T", a), fmt.Sprintf("%T", b)
	if ta == tb {
		t.Fatalf("two domains' refusals collapsed to one type %s", ta)
	}

	if a.Error() == b.Error() {
		t.Fatal("distinct domains rendered identical messages; they must not share rendering state")
	}
}
