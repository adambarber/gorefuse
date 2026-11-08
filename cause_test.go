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
		{"all-whitespace cause drops out", trailingError("   "), "refused. " + testFix},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertRender(t, because(Refusal{File: "", Line: 0, Subject: "", Why: Why(testFix), Sensitive: 0}, testRefused, tc.cause), renderText("Error: " + tc.want))
		})
	}
}

// trailingError is a cause whose message carries the trailing whitespace under
// test; Error() is its literal text.
type trailingError string

func (e trailingError) Error() string { return string(e) }
