package gorefuse

import "strings"

// Redacted is the fixed token a marked-sensitive slot renders as — grep-
// stable, resembling no real content, and scrub-clean. Not authenticated:
// unmarked content containing the literal [redacted] renders identically.
const Redacted = "[redacted]"

// Sensitive marks which slots of a Refusal carry secret-classified content: a
// marked Subject or Why renders as Redacted; a marked Cause renders nothing.
// Marks never cover Reason (arbitrary caller type R) — keep secrets out of it.
type Sensitive uint8

const (
	// SubjectSensitive marks the Subject slot as sensitive.
	SubjectSensitive Sensitive = 1 << iota

	// WhySensitive marks the Why slot as sensitive.
	WhySensitive

	// CauseSensitive marks the Cause slot: when marked, the entire cause tail
	// renders nothing, as if there were no cause — a tail has no token form.
	CauseSensitive
)

// sensitiveMask is the set of valid Sensitive bits; At refuses anything else.
const sensitiveMask = SubjectSensitive | WhySensitive | CauseSensitive

// slotText is one rendered slot's text after domain conversion.
type slotText string

// unknownSensitive reports whether m carries bits no construction path owns.
func unknownSensitive(m Sensitive) bool { return m&^sensitiveMask != 0 }

// renderSlot scrubs and trims a slot, then redacts it when marked. Redacting
// after scrubbing keeps a NUL-only slot empty instead of fabricating a token.
func renderSlot(text slotText, mark, marks Sensitive) slotText {
	clean := slotText(strings.TrimSpace(scrub(string(text))))
	if clean == "" || marks&mark == 0 {
		return clean
	}

	return slotText(Redacted)
}

// normalize fails closed on unknown bits on every construction path: the frame
// becomes invalidFrame, all three slots are marked, and the cause is suppressed.
func (r Refusal) normalize() Refusal {
	if !unknownSensitive(r.Sensitive) {
		return r
	}

	return Refusal{File: invalidFrame, Subject: r.Subject, Why: r.Why, Line: 0, Sensitive: sensitiveMask}
}
