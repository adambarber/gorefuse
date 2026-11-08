package gorefuse

import (
	"errors"
	"fmt"
	"testing"
)

// confError is a string-kinded reason domain whose constants are errors —
// the shape goconf documents ("match it with errors.Is against the reason
// constants"); gorefuse must deliver the matching mechanism.
type confError string

func (r confError) Error() string { return string(r) }

const errConfUnknownField confError = `unknown field "extra"`

func TestIsMatchesErrorReasonValue(t *testing.T) {
	b := Because[confError]{Reason: errConfUnknownField, Cause: nil, What: Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}}

	if !errors.Is(error(b), error(errConfUnknownField)) {
		t.Fatal("errors.Is(Because, reason value) = false, want true")
	}

	if !b.Is(errConfUnknownField) {
		t.Fatal("Is(reason itself) = false, want true")
	}
}

func TestIsWrappedBecauseStillMatches(t *testing.T) {
	b := Because[confError]{Reason: errConfUnknownField, Cause: nil, What: Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}}
	wrapped := fmt.Errorf("parsing config: %w", b)

	if !errors.Is(wrapped, errConfUnknownField) {
		t.Fatal("errors.Is(wrapped Because, reason) = false, want true through the chain")
	}
}

// isTestError is an error reason whose own Is delegates to the wrapped error —
// errors.Is semantics on the reason value must be honored, not bypassed.
type isTestError struct{ err error }

func (r isTestError) Error() string { return "isReason: " + r.err.Error() }

func (r isTestError) Is(target error) bool { return errors.Is(r.err, target) }

func TestIsReasonOwnIsMethodHonored(t *testing.T) {
	b := Because[isTestError]{Reason: isTestError{err: errTestSentinel}, Cause: nil, What: Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}}

	if !b.Is(errTestSentinel) {
		t.Fatal("Is(target matched by the reason's errors.Is) = false, want true")
	}
}

func TestIsNonErrorReasonNeverMatches(t *testing.T) {
	// A non-error R cannot equal an error target: the dynamic types differ,
	// so the any-equality is false and never panics.
	b := Because[testReason]{Reason: testRefused, Cause: nil, What: Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}}

	if b.Is(errTestSentinel) {
		t.Fatal("non-error reason matched an error target")
	}

	if b.Is(nil) {
		t.Fatal("nil target matched")
	}
}

func TestIsNilTargetFalse(t *testing.T) {
	b := Because[confError]{Reason: errConfUnknownField, Cause: nil, What: Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}}

	if b.Is(nil) {
		t.Fatal("Is(nil) = true, want false")
	}
}

func TestIsCrossDomainNeverBridges(t *testing.T) {
	b := Because[testReason]{Reason: testRefused, Cause: nil, What: Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}}
	other := Because[otherReason]{Reason: otherRefused, Cause: nil, What: Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}}

	if b.Is(other) {
		t.Fatal("Because[testReason] matched Because[otherReason]; the type wall leaked through Is")
	}

	if b.Is(errTestSentinel) {
		t.Fatal("non-error reason matched an error target")
	}

	if b.Is(nil) {
		t.Fatal("nil target matched")
	}
}
