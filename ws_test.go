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
