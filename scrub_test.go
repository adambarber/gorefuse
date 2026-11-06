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
			got := b.Error()
			if got != wants[name] {
				t.Fatalf("Error() = %q, want %q", got, wants[name])
			}

			if strings.ContainsRune(got, '\r') {
				t.Fatalf("Error() = %q still contains a raw carriage return", got)
			}
		})
	}
}

func TestInvalidUTF8BecomesReplacementChar(t *testing.T) {
	b := Because[testReason]{
		What:   Refusal{File: "", Line: 0, Subject: "", Why: ""},
		Reason: testReason("bad \xff\xfe tail"),
		Cause:  nil,
	}

	// strings.ToValidUTF8 replaces each run of invalid bytes with one U+FFFD.
	want := "bad \uFFFD tail"
	if got := b.Error(); got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestEmbeddedNewlineStillPassesThroughScrubbed(t *testing.T) {
	b := Because[testReason]{
		What:   Refusal{File: "f.toml", Line: 3, Subject: "", Why: ""},
		Reason: testReason("line one\nline two"),
		Cause:  nil,
	}

	want := "f.toml:3: line one\nline two"
	if got := b.Error(); got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestESCEscapeNeutralized(t *testing.T) {
	b := Because[testReason]{
		What:   Refusal{File: "", Line: 0, Subject: "a\x1b[31mb", Why: ""},
		Reason: "re\x1bfused",
		Cause:  nil,
	}

	want := "re\uFFFDfused: a\uFFFD[31mb"

	got := b.Error()
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}

	if strings.ContainsRune(got, '\x1b') {
		t.Fatalf("Error() = %q still contains a raw ESC byte", got)
	}
}

func TestBidiOverridePassesThrough(t *testing.T) {
	b := Because[testReason]{
		What:   Refusal{File: "", Line: 0, Subject: "", Why: ""},
		Reason: testReason("dir\u202Ecru/"),
		Cause:  nil,
	}

	want := "dir\u202Ecru/"
	if got := b.Error(); got != want {
		t.Fatalf("Error() = %q, want %q — U+202E is valid Unicode and passes through", got, want)
	}
}
