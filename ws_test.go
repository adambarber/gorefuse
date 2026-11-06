package gorefuse

import "testing"

func TestWhitespaceContract(t *testing.T) {
	const refusedSpace = "refused "

	const refusedTab = "refused\t"

	const refusedNewline = "refused\n"

	cases := []struct {
		name   string
		file   File
		reason testReason
		why    Why
		want   string
	}{
		{"why-less trailing space trimmed", "", refusedSpace, "", testRendered},
		{"why-less trailing tab trimmed", "", refusedTab, "", "refused"},
		{"why-less trailing newline trimmed", "", refusedNewline, "", "refused"},
		{"trailing space then why", "", refusedSpace, Why(testFix), "refused. " + testFix},
		{"trailing tab then why", "", refusedTab, Why(testFix), "refused. " + testFix},
		{"trailing newline then why", "", refusedNewline, Why(testFix), "refused. " + testFix},
		{"embedded newline kept with why", "", "refused\nmore", Why(testFix), "refused\nmore. " + testFix},
		{"whitespace-only reason degenerates", "", " ", Why(testFix), testFix},
		{"why without terminal punctuation", "", testRendered, "no period", "refused. no period"},
		{"frame with trailing slash", "/etc/f.toml/", testRendered, Why(testFix), "/etc/f.toml/: refused. " + testFix},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := Because[testReason]{
				What:   Refusal{File: tc.file, Line: 0, Subject: "", Why: tc.why},
				Reason: tc.reason,
				Cause:  nil,
			}
			if got := b.Error(); got != tc.want {
				t.Fatalf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSkipIfEmptyUniform(t *testing.T) {
	cases := []struct {
		name    string
		reason  testReason
		subject Subject
		file    File
		why     Why
		want    string
		line    Line
	}{
		{"empty reason with subject renders subject alone", "", testSubject, testFile, "", string(testFile) + ": " + string(testSubject), 0},
		{"empty reason with subject and why renders subject and why", "", testSubject, testFile, Why(testFix), string(testFile) + ": " + string(testSubject) + ". " + testFix, 0},
		{"empty reason with line", "", "", testFile, "", string(testFile) + ":3", 3},
		{"whitespace-only reason with subject and why degenerates", " ", testSubject, "", Why(testFix), string(testSubject) + ". " + testFix, 0},
		{"whitespace-only reason with subject alone", " ", testSubject, "", "", string(testSubject), 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := Because[testReason]{
				What:   Refusal{File: tc.file, Line: tc.line, Subject: tc.subject, Why: tc.why},
				Reason: tc.reason,
				Cause:  nil,
			}
			if got := b.Error(); got != tc.want {
				t.Fatalf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestFrameWhitespaceFileSuppressed(t *testing.T) {
	b := Because[testReason]{
		What:   Refusal{File: " ", Line: 3, Subject: testSubject, Why: ""},
		Reason: testRefused,
		Cause:  nil,
	}

	want := string(testRefused) + ": " + string(testSubject)
	if got := b.Error(); got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}
