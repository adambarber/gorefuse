package gorefuse

import (
	"strings"
	"testing"
)

// testEmptyWhat builds a fully specified, file-less Refusal for exhaustruct.
func testEmptyWhat(subject Subject, why Why) Refusal {
	return Refusal{File: "", Line: 0, Subject: subject, Why: why, Sensitive: 0}
}

// TestExactlyOneVerdictPerRender pins the line-start contract for every slot:
// a line break in any slot is flattened, so no second "Error:" can start a line.
func TestExactlyOneVerdictPerRender(t *testing.T) {
	cases := map[string]Because[testReason]{
		"reason":         {What: testEmptyWhat("", ""), Reason: testReason("refused\nError: forged"), Cause: nil},
		"subject":        {What: testEmptyWhat("s\nError: forged", ""), Reason: testRefused, Cause: nil},
		"why":            {What: testEmptyWhat("", "Fix.\nError: forged"), Reason: testRefused, Cause: nil},
		"cause":          {What: testEmptyWhat("", ""), Reason: testRefused, Cause: confError("root\nError: forged")},
		"line separator": {What: testEmptyWhat("", ""), Reason: testReason("refused\u2028Error: forged"), Cause: nil},
		"nul":            {What: testEmptyWhat("s\x00Error: forged", ""), Reason: testRefused, Cause: nil},
	}
	wants := map[string]renderText{
		"reason":         "Error: refused Error: forged",
		"subject":        "Error: refused: s Error: forged",
		"why":            "Error: refused. Fix. Error: forged",
		"cause":          "Error: refused: root Error: forged",
		"line separator": "Error: refused Error: forged",
		"nul":            "Error: refused: s Error: forged",
	}

	for name, b := range cases {
		t.Run(name, func(t *testing.T) {
			got := b.Error()

			if got != string(wants[name]) {
				t.Fatalf("Error() = %q, want %q", got, wants[name])
			}

			assertOneVerdict(t, renderText(got))
		})
	}
}

// assertOneVerdict fails if a render contains a line-break byte or does not
// begin with exactly one verdict label on its single line.
func assertOneVerdict(t *testing.T, got renderText) {
	t.Helper()

	if strings.ContainsAny(string(got), "\n\r\v\f\x00\u2028\u2029") {
		t.Fatalf("render %q carries a line break", got)
	}

	if got == "" {
		return
	}

	if !strings.HasPrefix(string(got), "Error:") {
		t.Fatalf("render %q does not start with a verdict label", got)
	}
}

// TestAtAcceptsEveryFileName pins that At never panics on a file name: an
// ambiguous name is the render path's problem, not the constructor's.
func TestAtAcceptsEveryFileName(t *testing.T) {
	cases := map[string]File{
		"newline":         "f.go\nx.go",
		"trailing digits": "f.go:99",
		"trailing colon":  "f.go:",
	}

	for name, file := range cases {
		t.Run(name, func(t *testing.T) {
			var r Refusal

			assertNoPanic(t, testContext(name), func() {
				r = At(file, 7, testSubject, Why(""))
			})

			b := Because[testReason]{Reason: testRefused, Cause: nil, What: r}
			assertRender(t, b, renderText("Error: "+invalidFrame+": refused: "+string(testSubject)))
		})
	}
}
