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

const wantFrameSubject = "Error: f.toml: refused: server.tls"

// TestVerdictPrefixPinsErrorLabel pins the owner's verdict-label ruling:
// Error() output begins with "Error: " so a render is marked as an error.
func TestVerdictPrefixPinsErrorLabel(t *testing.T) {
	assertRender(t, because(Refusal{File: testFile, Line: 3, Subject: testSubject, Why: "", Sensitive: 0}, testRefused, nil),
		renderText("Error: "+string(testFile)+":3: "+string(testRefused)+": "+string(testSubject)))
}

func TestReasonAlone(t *testing.T) {
	assertRender(t, because(Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, testMissing, nil), "Error: missing")
}

func TestFramePrefixWholeFile(t *testing.T) {
	assertRender(t, because(Refusal{File: "config.toml", Line: 0, Subject: "", Why: "", Sensitive: 0}, testRefused, nil),
		"Error: config.toml: refused")
}

func TestFramePrefixWithLine(t *testing.T) {
	assertRender(t, because(Refusal{File: "config.toml", Line: 12, Subject: "", Why: "", Sensitive: 0}, testRefused, nil),
		"Error: config.toml:12: refused")
}

func TestSubjectAppends(t *testing.T) {
	assertRender(t, because(Refusal{File: testFile, Line: 0, Subject: testSubject, Why: "", Sensitive: 0}, testRefused, nil),
		wantFrameSubject)
}

func TestWhyAppendsAsSentence(t *testing.T) {
	assertRender(t, because(Refusal{File: testFile, Line: 0, Subject: testSubject, Why: "Rotate the certificate.", Sensitive: 0}, testRefused, nil),
		"Error: f.toml: refused: server.tls. Rotate the certificate.")
}

func TestCauseTail(t *testing.T) {
	assertRender(t, because(Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, testRefused, errTestSentinel),
		"Error: refused: root cause")
}

func TestFullRenderStableOrder(t *testing.T) {
	assertRender(t, because(Refusal{File: testFile, Line: 3, Subject: testSubject, Why: "Renew it.", Sensitive: 0}, testRefused, errTestSentinel),
		"Error: f.toml:3: refused: server.tls: root cause. Renew it.")
}

func TestUnwrapToCause(t *testing.T) {
	b := because(Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, testRefused, errTestSentinel)
	if !errors.Is(b, errTestSentinel) {
		t.Fatal("errors.Is(Because, sentinel) = false, want true")
	}
}

func TestUnwrapNilCause(t *testing.T) {
	b := because(Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, testRefused, nil)
	if b.Unwrap() != nil {
		t.Fatalf("Unwrap() = %v, want nil", b.Unwrap())
	}
}

func TestReasonRendersStringer(t *testing.T) {
	assertRender(t, because(Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, stringerReason{}, nil), "Error: rendered")
}

type stringerReason struct{}

func (stringerReason) String() string { return "rendered" }

func TestDistinctReasonTypesAreDistinctTypes(t *testing.T) {
	a := because(Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, testRefused, nil)
	b := because(Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, otherRefused, nil)

	ta, tb := fmt.Sprintf("%T", a), fmt.Sprintf("%T", b)
	if ta == tb {
		t.Fatalf("two domains' refusals collapsed to one type %s", ta)
	}

	if a.Error() == b.Error() {
		t.Fatal("distinct domains rendered identical messages; they must not share rendering state")
	}
}
