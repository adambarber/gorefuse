package gorefuse

import (
	"encoding/json"
	"fmt"
)

type jsonCause string

// MarshalJSON keeps marked slots redacted in structured encoders. The prose
// renderer is not the only output surface of a refusal.
func (r Refusal) MarshalJSON() ([]byte, error) {
	n := r.normalize()

	subject := n.Subject
	if n.Sensitive&SubjectSensitive != 0 {
		subject = Subject(Redacted)
	}

	why := n.Why
	if n.Sensitive&WhySensitive != 0 {
		why = Why(Redacted)
	}

	encoded, err := json.Marshal(struct {
		File      File      `json:"File"`
		Subject   Subject   `json:"Subject"`
		Why       Why       `json:"Why"`
		Line      Line      `json:"Line"`
		Sensitive Sensitive `json:"Sensitive"`
	}{
		File:      n.File,
		Subject:   subject,
		Why:       why,
		Line:      n.Line,
		Sensitive: n.Sensitive,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal refusal: %w", err)
	}

	return encoded, nil
}

// MarshalJSON applies the same slot decisions as Error: a marked cause is
// omitted, and Reason remains the caller's domain value.
func (b Because[R]) MarshalJSON() ([]byte, error) {
	n := b.What.normalize()

	out := struct {
		Reason R         `json:"Reason"`
		Cause  jsonCause `json:"Cause,omitempty"`
		What   Refusal   `json:"What"`
	}{
		Reason: b.Reason,
		Cause:  jsonCause(""),
		What:   b.What,
	}

	if b.Cause != nil && n.Sensitive&CauseSensitive == 0 {
		out.Cause = jsonCause(b.Cause.Error())
	}

	encoded, err := json.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("marshal refusal carrier: %w", err)
	}

	return encoded, nil
}
