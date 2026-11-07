package gorefuse

import (
	"errors"
	"strings"
	"testing"
)

func TestMarkedSubjectRendersToken(t *testing.T) {
	b := because(Refusal{File: testFile, Line: 0, Subject: testSubject, Why: "", Sensitive: 0}, testRefused, nil)
	b.What.Sensitive = SubjectSensitive

	assertRender(t, b, "f.toml: refused: "+Redacted)
}

func TestMarkedWhyRendersToken(t *testing.T) {
	b := because(Refusal{File: "", Line: 0, Subject: "", Why: "Rotate the secret now.", Sensitive: 0}, testRefused, nil)
	b.What.Sensitive = WhySensitive

	assertRender(t, b, testRendered+". "+Redacted)
}

func TestBothSlotsMarkedRenderToken(t *testing.T) {
	b := because(Refusal{File: "", Line: 0, Subject: testSubject, Why: "Rotate the secret now.", Sensitive: 0}, testRefused, nil)
	b.What.Sensitive = SubjectSensitive | WhySensitive

	assertRender(t, b, testRendered+": "+Redacted+". "+Redacted)
}

func TestMarkedEmptyWhyStillContributesNothing(t *testing.T) {
	// The empty-after-trim rule wins even for a marked slot: an empty slot
	// contributes nothing, connector and all — redaction never fabricates one.
	b := because(Refusal{File: "", Line: 0, Subject: "", Why: "  ", Sensitive: 0}, testRefused, nil)
	b.What.Sensitive = WhySensitive

	assertRender(t, b, testRendered)
}

func TestMarkedSlotTokenSurvivesScrub(t *testing.T) {
	b := because(Refusal{File: "", Line: 0, Subject: "a\rb\x1bc\xff", Why: "x\r\x1b\xff", Sensitive: 0}, testRefused, nil)
	b.What.Sensitive = SubjectSensitive | WhySensitive

	assertRender(t, b, testRendered+": "+Redacted+". "+Redacted)
}

func TestMarkingLeavesIdentityIntact(t *testing.T) {
	b := because(Refusal{File: "", Line: 0, Subject: testSubject, Why: "Rotate it.", Sensitive: 0}, testRefused, errTestSentinel)
	b.What.Sensitive = SubjectSensitive | WhySensitive

	if !errors.Is(b, errTestSentinel) {
		t.Fatal("errors.Is(Because, sentinel) = false with marked slots, want true")
	}
}

func TestAtMarksSensitiveSlots(t *testing.T) {
	r := At(testFile, 0, testSubject, "Rotate it.", SubjectSensitive, WhySensitive)

	if r.Sensitive != SubjectSensitive|WhySensitive {
		t.Fatalf("Sensitive = %d, want both bits set", r.Sensitive)
	}
}

func TestUnmarkedAtUnchanged(t *testing.T) {
	r := At(testFile, 3, testSubject, "Rotate it.")

	if r.Sensitive != 0 {
		t.Fatalf("Sensitive = %d, want 0", r.Sensitive)
	}
}

func TestRedactedTokenIsClean(t *testing.T) {
	if strings.TrimSpace(Redacted) == "" {
		t.Fatal("token is whitespace")
	}
}
