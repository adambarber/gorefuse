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
	var chain error = Because[testReason]{What: Refusal{File: "", Line: 0, Subject: "", Why: ""}, Reason: testRefused, Cause: nil}

	if _, ok := errors.AsType[Because[otherReason]](chain); ok {
		t.Fatal("errors.AsType matched a Because[testReason] as Because[otherReason]; the domains leaked")
	}

	if _, ok := errors.AsType[Because[testReason]](chain); !ok {
		t.Fatal("errors.AsType failed to match its own domain")
	}
}

func TestEmptyBecauseRendersEmpty(t *testing.T) {
	b := Because[testReason]{What: Refusal{File: "", Line: 0, Subject: "", Why: ""}, Reason: "", Cause: nil}
	if got := b.Error(); got != "" {
		t.Fatalf("Error() = %q, want empty", got)
	}
}

func TestCauseTailAfterSentencePunctuation(t *testing.T) {
	b := Because[testReason]{
		What:   Refusal{File: "", Line: 0, Subject: "", Why: "Renew it."},
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
		What:   Refusal{File: "", Line: 0, Subject: "server.tls!", Why: "Rename it."},
		Reason: testRefused,
		Cause:  nil,
	}

	got := b.Error()
	if got != "refused: server.tls! Rename it." {
		t.Fatalf("Error() = %q, want single punctuation before the fix", got)
	}
}

func TestDescribableInterfaceSatisfied(t *testing.T) {
	var d Describable = testDescriber{what: Refusal{File: testFile, Line: 0, Subject: testSubject, Why: ""}, misses: false}

	b := From(testRefused, d)

	if b.Error() != wantFrameSubject {
		t.Fatalf("Error() = %q, want %q", b.Error(), wantFrameSubject)
	}
}
