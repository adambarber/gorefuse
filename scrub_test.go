package gorefuse

import (
	"strings"
	"testing"
)

func TestCarriageReturnNeutralized(t *testing.T) {
	cases := map[string]Because[testReason]{
		"CR in subject": {What: Refusal{File: "", Line: 0, Subject: "a\rb", Why: ""}, Reason: testRefused, Cause: nil},
		"CR in why":     {What: Refusal{File: "", Line: 0, Subject: "", Why: "Fix it.\rSoon."}, Reason: testRefused, Cause: nil},
		"CR in file":    {What: Refusal{File: "a\rb.go", Line: 3, Subject: "", Why: ""}, Reason: testRefused, Cause: nil},
		"CR in reason":  {What: Refusal{File: "", Line: 0, Subject: "", Why: ""}, Reason: "re\rfused", Cause: nil},
	}

	wants := map[string]string{
		"CR in subject": "refused: a b",
		"CR in why":     "refused. Fix it. Soon.",
		"CR in file":    "a b.go:3: refused",
		"CR in reason":  "re fused",
	}

	for name, b := range cases {
		t.Run(name, func(t *testing.T) {
			assertRender(t, b, renderText(wants[name]))

			if got := b.Error(); strings.ContainsRune(got, '\r') {
				t.Fatalf("Error() = %q still contains a raw carriage return", got)
			}
		})
	}
}

func TestInvalidUTF8BecomesReplacementChar(t *testing.T) {
	// strings.ToValidUTF8 replaces each run of invalid bytes with one U+FFFD.
	assertRender(t, because(Refusal{File: "", Line: 0, Subject: "", Why: ""}, testReason("bad \xff\xfe tail"), nil),
		"bad \uFFFD tail")
}

func TestESCEscapeNeutralized(t *testing.T) {
	b := because(Refusal{File: "", Line: 0, Subject: "a\x1b[31mb", Why: ""}, "re\x1bfused", nil)

	assertRender(t, b, "re\uFFFDfused: a\uFFFD[31mb")

	if got := b.Error(); strings.ContainsRune(got, '\x1b') {
		t.Fatalf("Error() = %q still contains a raw ESC byte", got)
	}
}

func TestBidiOverridePassesThrough(t *testing.T) {
	assertRender(t, because(Refusal{File: "", Line: 0, Subject: "", Why: ""}, testReason("dir\u202Ecru/"), nil),
		"dir\u202Ecru/")
}

func TestWhitespaceOnlyWhyContributesNothing(t *testing.T) {
	cases := map[string]Because[testReason]{
		"why is space and tab":      {What: Refusal{File: "", Line: 0, Subject: "", Why: " \t"}, Reason: testRefused, Cause: nil},
		"why is newline":            {What: Refusal{File: "", Line: 0, Subject: "", Why: "\n"}, Reason: testRefused, Cause: nil},
		"real why still renders":    {What: Refusal{File: "", Line: 0, Subject: "", Why: " \t" + Why(testFix) + " "}, Reason: testRefused, Cause: nil},
		"why-less frame unaffected": {What: Refusal{File: testFile, Line: 0, Subject: "", Why: ""}, Reason: testRefused, Cause: nil},
	}

	wants := map[string]string{
		"why is space and tab":      testRendered,
		"why is newline":            testRendered,
		"real why still renders":    testRendered + ". " + testFix,
		"why-less frame unaffected": string(testFile) + ": " + testRendered,
	}

	for name, b := range cases {
		t.Run(name, func(t *testing.T) {
			assertRender(t, b, renderText(wants[name]))
		})
	}
}
