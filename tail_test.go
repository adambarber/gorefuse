package gorefuse

import "testing"

// Han-script rune fires gosmopolitan; the others are rune-spelled for symmetry.
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
			assertRender(t, because(Refusal{File: "", Line: 0, Subject: tc.subject, Why: tc.why, Sensitive: 0}, tc.reason, nil), renderText(tc.want))
		})
	}
}
