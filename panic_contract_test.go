package gorefuse

import (
	"fmt"
	"strings"
	"testing"
)

// panicText is a named kind for panic payloads and expected fragments.
type panicText string

// panickingDescriber violates Describe's contract: Describe must not panic.
type panickingDescriber struct{}

func (panickingDescriber) Describe() (Refusal, bool) { panic(panicText("describe panicked")) }

// panickingCauseError violates Cause's contract: Error must not panic.
type panickingCauseError struct{}

func (panickingCauseError) Error() string { panic(panicText("cause panicked")) }

// assertPanicPropagates pins that a caller-misbehaving hook panics through the
// library rather than being recovered.
func assertPanicPropagates(t *testing.T, accepted testContext, call func(), want panicText) {
	t.Helper()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("%s; want a panic", accepted)
		}

		if got := fmt.Sprint(r); !strings.Contains(got, string(want)) {
			t.Fatalf("panic %q does not echo %q", got, want)
		}
	}()

	call()
}

func TestPanickingDescribePropagates(t *testing.T) {
	assertPanicPropagates(t, "From recovered a panicking Describe", func() {
		got := From(testRefused, panickingDescriber{})
		_ = got
	}, "describe panicked")
}

func TestPanickingCausePropagates(t *testing.T) {
	b := Because[testReason]{Reason: testRefused, What: testEmptyWhat("", ""), Cause: panickingCauseError{}}

	assertPanicPropagates(t, "Error recovered a panicking Cause", func() {
		_ = b.Error()
	}, "cause panicked")
}
