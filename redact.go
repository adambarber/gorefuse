package gorefuse

import "strings"

// Redacted is the fixed token a marked-sensitive slot renders as. It is
// grep-stable, resembles no real content, and is itself scrub-clean, so the
// other scrub stages transform it to itself.
const Redacted = "[redacted]"

// Sensitive marks which slots of a Refusal carry content the caller has
// classified as sensitive: a marked slot renders as Redacted, never its
// content. The caller classifies; this guarantees only the marked case.
type Sensitive uint8

const (
	// SubjectSensitive marks the Subject slot as sensitive.
	SubjectSensitive Sensitive = 1 << iota

	// WhySensitive marks the Why slot as sensitive.
	WhySensitive
)

// markedSensitive renders a slot the redaction stage produces: the token
// when marked, the text otherwise. An empty slot stays empty — redaction
// never fabricates a contribution.
func markedSensitive(text Subject, mark Sensitive, marks Sensitive) Subject {
	return markedSensitiveWhy(Why(text), mark, marks).toSubject()
}

// markedSensitiveWhy is the redaction stage for the Why slot.
func markedSensitiveWhy(text Why, mark Sensitive, marks Sensitive) Why {
	if strings.TrimSpace(string(text)) == "" || marks&mark == 0 {
		return text
	}

	return Why(Redacted)
}

// toSubject spells a redacted Why back in the Subject domain kind.
func (w Why) toSubject() Subject { return Subject(w) }
