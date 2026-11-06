package gorefuse

import "testing"

func TestCauseTailWhitespaceTrimmed(t *testing.T) {
	cases := []struct {
		name  string
		cause error
		want  string
	}{
		{"cause trailing space", trailingError("inner "), "refused: inner. " + testFix},
		{"cause trailing tab", trailingError("inner\t"), "refused: inner. " + testFix},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := Because[testReason]{
				What:   Refusal{File: "", Line: 0, Subject: "", Why: Why(testFix)},
				Reason: testRefused,
				Cause:  tc.cause,
			}
			if got := b.Error(); got != tc.want {
				t.Fatalf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}

// trailingError is a cause whose message carries the trailing whitespace under
// test; Error() is its literal text.
type trailingError string

func (e trailingError) Error() string { return string(e) }
