package gorefuse

import (
	"errors"
	"fmt"
	"strings"
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

const wantFrameSubject = "Error: f.toml: refused: server.tls"

// TestVerdictPrefixPinsErrorLabel pins the owner's verdict-label ruling:
// Error() output begins with "Error: " so a render is marked as an error.
func TestVerdictPrefixPinsErrorLabel(t *testing.T) {
	assertRender(t, Wrap(testRefused, Refusal{File: testFile, Line: 3, Subject: testSubject, Why: "", Sensitive: 0}, nil),
		renderText("Error: "+string(testFile)+":3: "+string(testRefused)+": "+string(testSubject)))
}

func TestReasonAlone(t *testing.T) {
	assertRender(t, Wrap(testMissing, Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, nil), "Error: missing")
}

func TestFramePrefixWholeFile(t *testing.T) {
	assertRender(t, Wrap(testRefused, Refusal{File: "config.toml", Line: 0, Subject: "", Why: "", Sensitive: 0}, nil),
		"Error: config.toml: refused")
}

func TestFramePrefixWithLine(t *testing.T) {
	assertRender(t, Wrap(testRefused, Refusal{File: "config.toml", Line: 12, Subject: "", Why: "", Sensitive: 0}, nil),
		"Error: config.toml:12: refused")
}

func TestSubjectAppends(t *testing.T) {
	assertRender(t, Wrap(testRefused, Refusal{File: testFile, Line: 0, Subject: testSubject, Why: "", Sensitive: 0}, nil),
		wantFrameSubject)
}

func TestWhyAppendsAsSentence(t *testing.T) {
	assertRender(t, Wrap(testRefused, Refusal{File: testFile, Line: 0, Subject: testSubject, Why: "Rotate the certificate.", Sensitive: 0}, nil),
		"Error: f.toml: refused: server.tls. Rotate the certificate.")
}

func TestCauseTail(t *testing.T) {
	assertRender(t, Wrap(testRefused, Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, errTestSentinel),
		"Error: refused: root cause")
}

// TestNestedBecauseCauseStaysUnlabelled pins the owner's ruling: a Because
// used as a cause contributes content, not a second verdict label. Only the
// outer Verdict labels the line, so wrapping and nesting read as prose.
func TestNestedBecauseCauseStaysUnlabelled(t *testing.T) {
	inner := Because[testReason]{Reason: testRefused, Cause: nil, What: Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}}
	outer := Because[testReason]{Reason: testMissing, Cause: inner, What: Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}}

	assertRender(t, outer, "Error: missing: refused")

	for _, line := range strings.Split(outer.Error(), "\n") {
		if strings.Contains(line, "Error: ") {
			t.Fatalf("Error() = %q, a nested carrier must not embed a label", line)
		}
	}
}

func TestFullRenderStableOrder(t *testing.T) {
	assertRender(t, Wrap(testRefused, Refusal{File: testFile, Line: 3, Subject: testSubject, Why: "Renew it.", Sensitive: 0}, errTestSentinel),
		"Error: f.toml:3: refused: server.tls: root cause. Renew it.")
}

func TestUnwrapToCause(t *testing.T) {
	b := Wrap(testRefused, Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, errTestSentinel)
	if !errors.Is(b, errTestSentinel) {
		t.Fatal("errors.Is(Because, sentinel) = false, want true")
	}
}

func TestUnwrapNilCause(t *testing.T) {
	b := Wrap(testRefused, Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, nil)
	if b.Unwrap() != nil {
		t.Fatalf("Unwrap() = %v, want nil", b.Unwrap())
	}
}

func TestReasonRendersStringer(t *testing.T) {
	assertRender(t, Wrap(stringerReason{}, Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, nil), "Error: rendered")
}

type stringerReason struct{}

func (stringerReason) String() string { return "rendered" }

func TestDistinctReasonTypesAreDistinctTypes(t *testing.T) {
	a := Wrap(testRefused, Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, nil)
	b := Wrap(otherRefused, Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, nil)

	ta, tb := fmt.Sprintf("%T", a), fmt.Sprintf("%T", b)
	if ta == tb {
		t.Fatalf("two domains' refusals collapsed to one type %s", ta)
	}

	if a.Error() == b.Error() {
		t.Fatal("distinct domains rendered identical messages; they must not share rendering state")
	}
}
