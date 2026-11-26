package gorefuse

import (
	"errors"
	"strings"
	"testing"
)

// The cross-domain guarantee is enforced by the compiler: the assignment in
// compilecheck_test.go must fail to build. This test pins the runtime half:
// errors.As never matches across reason domains.
func TestCrossDomainErrorsAsDistinct(t *testing.T) {
	var chain error = Because[testReason]{What: Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, Reason: testRefused, Cause: nil}

	var other Because[otherReason]
	if errors.As(chain, &other) {
		t.Fatalf("As matched a Because[testReason] as Because[otherReason]: %v", other)
	}

	var own Because[testReason]
	if !errors.As(chain, &own) {
		t.Fatalf("As failed to match its own domain: %v", own)
	}
}

func TestEmptyBecauseRendersEmpty(t *testing.T) {
	b := Because[testReason]{What: Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, Reason: "", Cause: nil}
	if got := b.Error(); got != "" {
		t.Fatalf("Error() = %q, want empty", got)
	}
}

func TestNilErrorReasonRendersEmpty(t *testing.T) {
	var reason error
	if got := From(reason, nil).Error(); got != "" {
		t.Fatalf("Error() = %q, want empty for a nil error reason", got)
	}
}

func TestCauseTailAfterSentencePunctuation(t *testing.T) {
	b := Because[testReason]{
		What:   Refusal{File: "", Line: 0, Subject: "", Why: "Renew it.", Sensitive: 0},
		Reason: testRefused,
		Cause:  errTestSentinel,
	}

	got := b.Error()
	if !strings.HasSuffix(got, "refused: root cause. Renew it.") {
		t.Fatalf("Error() = %q, want fix as the final sentence", got)
	}
}

func TestCauseTailNoDoublePunctuation(t *testing.T) {
	b := Because[testReason]{
		What:   Refusal{File: "", Line: 0, Subject: "server.tls!", Why: "Rename it.", Sensitive: 0},
		Reason: testRefused,
		Cause:  nil,
	}

	got := b.Verdict()
	if got != "Error: refused: server.tls! Rename it." {
		t.Fatalf("Verdict() = %q, want single punctuation before the fix", got)
	}
}

func TestDescribableInterfaceSatisfied(t *testing.T) {
	var d Describable = testDescriber{what: Refusal{File: testFile, Line: 0, Subject: testSubject, Why: "", Sensitive: 0}, misses: false}

	b := From(testRefused, d)

	if b.Verdict() != wantFrameSubject {
		t.Fatalf("Verdict() = %q, want %q", b.Verdict(), wantFrameSubject)
	}
}
