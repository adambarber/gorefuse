package gorefuse

import "testing"

const (
	nulReasonLabel     = "nul reason"
	nulSubjectLabel    = "nul subject"
	nulWhyLabel        = "nul why"
	nulRenderedRefusal = renderText("Error: " + string(testRefused))
)

// A slot whose only content is NUL becomes empty after scrubbing. An empty
// slot contributes no connector and fabricates no redaction token.
func TestNulOnlySlotsContributeNothing(t *testing.T) {
	cases := map[string]struct {
		want renderText
		b    Because[testReason]
	}{
		nulReasonLabel:   {want: "", b: Wrap(testReason("\x00\x00"), testEmptyWhat("", ""), nil)},
		nulSubjectLabel:  {want: nulRenderedRefusal, b: Wrap(testRefused, testEmptyWhat("\x00", ""), nil)},
		nulWhyLabel:      {want: nulRenderedRefusal, b: Wrap(testRefused, testEmptyWhat("", "\x00"), nil)},
		"all slots":      {want: "", b: Wrap(testReason("\x00"), testEmptyWhat("\x00", "\x00"), nil)},
		"marked subject": {want: nulRenderedRefusal, b: Wrap(testRefused, Refusal{File: "", Subject: "\x00", Why: "", Line: 0, Sensitive: SubjectSensitive}, nil)},
		"marked why":     {want: nulRenderedRefusal, b: Wrap(testRefused, Refusal{File: "", Subject: "", Why: "\x00", Line: 0, Sensitive: WhySensitive}, nil)},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assertRender(t, tc.b, tc.want)
		})
	}
}
