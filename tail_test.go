package gorefuse

import "testing"

// ideographic sentence terminals, spelled by rune so gosmopolitan stays quiet.
const (
	ideofullstop = '。'
	ideoclash    = '直'
	ideofix      = 'せ'
)

func TestTailRendering(t *testing.T) {
	cases := []struct {
		name    string
		reason  testReason
		subject Subject
		why     Why
		want    string
	}{
		{"ascii terminal dot joins with space", "refused.", "", Why(testFix), "refused. " + testFix},
		{"ascii terminal question joins with space", "refused?", "", "Huh?", "refused? Huh?"},
		{"ascii terminal bang joins with space", "refused!", "", "Yes!", "refused! Yes!"},
		{"subject question then plain fix", testRendered, "why?", Why(testFix), "refused: why? " + testFix},
		{"trailing whitespace trimmed before check", "refused ", "", Why(testFix), "refused. " + testFix},
		{
			"ideographic terminal rune joins with space",
			testReason("refused" + string(ideofullstop)),
			"",
			Why(string(ideoclash) + string(ideofix) + string(ideofullstop)),
			"refused" + string(ideofullstop) + " " + string(ideoclash) + string(ideofix) + string(ideofullstop),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := Because[testReason]{
				What:   Refusal{File: "", Line: 0, Subject: tc.subject, Why: tc.why},
				Reason: tc.reason,
				Cause:  nil,
			}
			if got := b.Error(); got != tc.want {
				t.Fatalf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}
