package gorefuse

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// confError is a string-kinded reason domain whose constants are errors —
// the shape goconf documents ("match it with errors.Is against the reason
// constants"); gorefuse must deliver the matching mechanism.
type confError string

func (r confError) Error() string { return string(r) }

const errConfUnknownField confError = `unknown field "extra"`

// The zero-What carrier under test is built through the real API: From with
// a nil describable yields the reason alone — the same shape the Is tests
// need, without re-implementing From's construction.

func TestIsMatchesErrorReasonValue(t *testing.T) {
	b := From(errConfUnknownField, nil)

	if !errors.Is(error(b), error(errConfUnknownField)) {
		t.Fatal("errors.Is(Because, reason value) = false, want true")
	}

	if !b.Is(errConfUnknownField) {
		t.Fatal("Is(reason itself) = false, want true")
	}
}

func TestIsWrappedBecauseStillMatches(t *testing.T) {
	wrapped := fmt.Errorf("parsing config: %w", From(errConfUnknownField, nil))

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
	b := From(isTestError{err: errTestSentinel}, nil)

	if !b.Is(errTestSentinel) {
		t.Fatal("Is(target matched by the reason's errors.Is) = false, want true")
	}
}

// assertNonErrorRefusesTarget fails unless b refuses the error target and
// nil — the shared refusals of a non-error reason domain.
func assertNonErrorRefusesTarget(t *testing.T, b Because[testReason]) {
	t.Helper()

	if b.Is(errTestSentinel) {
		t.Fatal("non-error reason matched an error target")
	}

	if b.Is(nil) {
		t.Fatal("nil target matched")
	}
}

func TestIsNonErrorReasonNeverMatches(t *testing.T) {
	// A non-error R cannot equal an error target: the dynamic types differ,
	// so the any-equality is false and never panics.
	assertNonErrorRefusesTarget(t, From(testRefused, nil))
}

func TestIsNilTargetFalse(t *testing.T) {
	if From(errConfUnknownField, nil).Is(nil) {
		t.Fatal("Is(nil) = true, want false")
	}
}

func TestIsCrossDomainNeverBridges(t *testing.T) {
	b := From(testRefused, nil)
	other := From(otherRefused, nil)

	if b.Is(other) {
		t.Fatal("Because[testReason] matched Because[otherReason]; the type wall leaked through Is")
	}

	assertNonErrorRefusesTarget(t, b)
}

// TestIsAfterWrappingFrom pins that the real construction path (From) and an
// errors.New-wrapped chain still answer Is through the carrier.
func TestIsAfterWrappingFrom(t *testing.T) {
	var chain error = fmt.Errorf("outer: %w", fmt.Errorf("mid: %w", From(errConfUnknownField, nil)))

	if !errors.Is(chain, errConfUnknownField) {
		t.Fatal("errors.Is(deeply wrapped From-built Because, reason) = false, want true")
	}

	if !strings.Contains(chain.Error(), errConfUnknownField.Error()) {
		t.Fatal("wrapped chain lost the reason text")
	}
}
