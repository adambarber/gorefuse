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
		nulReasonLabel:   {want: "", b: because(testEmptyWhat("", ""), testReason("\x00\x00"), nil)},
		nulSubjectLabel:  {want: nulRenderedRefusal, b: because(testEmptyWhat("\x00", ""), testRefused, nil)},
		nulWhyLabel:      {want: nulRenderedRefusal, b: because(testEmptyWhat("", "\x00"), testRefused, nil)},
		"all slots":      {want: "", b: because(testEmptyWhat("\x00", "\x00"), testReason("\x00"), nil)},
		"marked subject": {want: nulRenderedRefusal, b: because(Refusal{File: "", Subject: "\x00", Why: "", Line: 0, Sensitive: SubjectSensitive}, testRefused, nil)},
		"marked why":     {want: nulRenderedRefusal, b: because(Refusal{File: "", Subject: "", Why: "\x00", Line: 0, Sensitive: WhySensitive}, testRefused, nil)},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assertRender(t, tc.b, tc.want)
		})
	}
}
