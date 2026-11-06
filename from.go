package gorefuse

// From builds a Because[R] from whatever the describable can say about the
// refusal. A nil Describable, or one whose Describe reports no description,
// leaves the reason standing alone with no frame.
func From[R any](reason R, d Describable) Because[R] {
	b := Because[R]{Reason: reason, Cause: nil, What: Refusal{File: "", Line: 0, Subject: "", Why: ""}}

	if d == nil {
		return b
	}

	if what, ok := d.Describe(); ok {
		b.What = what
	}

	return b
}
