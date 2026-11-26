package gorefuse

import (
	"strings"
	"testing"
)

func TestCarriageReturnNeutralized(t *testing.T) {
	cases := map[string]Because[testReason]{
		"CR in subject": {What: Refusal{File: "", Line: 0, Subject: "a\rb", Why: "", Sensitive: 0}, Reason: testRefused, Cause: nil},
		"CR in why":     {What: Refusal{File: "", Line: 0, Subject: "", Why: "Fix it.\rSoon.", Sensitive: 0}, Reason: testRefused, Cause: nil},
		"CR in file":    {What: Refusal{File: "a\rb.go", Line: 3, Subject: "", Why: "", Sensitive: 0}, Reason: testRefused, Cause: nil},
		"CR in reason":  {What: Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, Reason: "re\rfused", Cause: nil},
	}

	wants := map[string]string{
		"CR in subject": "refused: a b",
		"CR in why":     "refused. Fix it. Soon.",
		"CR in file":    invalidFrame + ": refused",
		"CR in reason":  "re fused",
	}

	for name, b := range cases {
		t.Run(name, func(t *testing.T) {
			assertRender(t, b, renderText("Error: "+wants[name]))

			if got := b.Error(); strings.ContainsRune(got, '\r') {
				t.Fatalf("Error() = %q still contains a raw carriage return", got)
			}
		})
	}
}

func TestInvalidUTF8BecomesReplacementChar(t *testing.T) {
	// strings.ToValidUTF8 replaces each run of invalid bytes with one U+FFFD.
	assertRender(t, Wrap(testReason("bad \xff\xfe tail"), Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, nil),
		"Error: bad \uFFFD tail")
}

func TestESCEscapeNeutralized(t *testing.T) {
	b := Wrap("re\x1bfused", Refusal{File: "", Line: 0, Subject: "a\x1b[31mb", Why: "", Sensitive: 0}, nil)

	assertRender(t, b, "Error: re\uFFFDfused: a\uFFFD[31mb")

	if got := b.Error(); strings.ContainsRune(got, '\x1b') {
		t.Fatalf("Error() = %q still contains a raw ESC byte", got)
	}
}

func TestLineBreaksFlattened(t *testing.T) {
	cases := map[string]Because[testReason]{
		"line feed in reason":       {What: testEmptyWhat("", ""), Reason: testReason("a\nb"), Cause: nil},
		"line separator in subject": {What: testEmptyWhat("a\u2028b", ""), Reason: testRefused, Cause: nil},
		"nel in subject":            {What: testEmptyWhat("a\u0085b", ""), Reason: testRefused, Cause: nil},
		"nel cannot forge verdict":  {What: testEmptyWhat("a\u0085Error: forged", ""), Reason: testRefused, Cause: nil},
		"nul in why":                {What: testEmptyWhat("", "a\x00b"), Reason: testRefused, Cause: nil},
	}
	wants := map[string]string{
		"line feed in reason":       "Error: a b",
		"line separator in subject": "Error: refused: a b",
		"nel in subject":            "Error: refused: a b",
		"nel cannot forge verdict":  "Error: refused: a Error: forged",
		"nul in why":                "Error: refused. a b",
	}

	for name, b := range cases {
		t.Run(name, func(t *testing.T) {
			assertRender(t, b, renderText(wants[name]))

			assertOneVerdict(t, renderText(b.Verdict()))
		})
	}
}

func TestBidiOverridePassesThrough(t *testing.T) {
	assertRender(t, Wrap(testReason("dir\u202Ecru/"), Refusal{File: "", Line: 0, Subject: "", Why: "", Sensitive: 0}, nil),
		"Error: dir\u202Ecru/")
}

func TestWhitespaceOnlyWhyContributesNothing(t *testing.T) {
	cases := map[string]Because[testReason]{
		"why is space and tab":      {What: Refusal{File: "", Line: 0, Subject: "", Why: " \t", Sensitive: 0}, Reason: testRefused, Cause: nil},
		"why is newline":            {What: Refusal{File: "", Line: 0, Subject: "", Why: "\n", Sensitive: 0}, Reason: testRefused, Cause: nil},
		"real why still renders":    {What: Refusal{File: "", Line: 0, Subject: "", Why: " \t" + Why(testFix) + " ", Sensitive: 0}, Reason: testRefused, Cause: nil},
		"why-less frame unaffected": {What: Refusal{File: testFile, Line: 0, Subject: "", Why: "", Sensitive: 0}, Reason: testRefused, Cause: nil},
	}

	wants := map[string]string{
		"why is space and tab":      testRendered,
		"why is newline":            testRendered,
		"real why still renders":    testRendered + ". " + testFix,
		"why-less frame unaffected": string(testFile) + ": " + testRendered,
	}

	for name, b := range cases {
		t.Run(name, func(t *testing.T) {
			assertRender(t, b, renderText("Error: "+wants[name]))
		})
	}
}
