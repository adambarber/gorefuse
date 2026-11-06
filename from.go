package gorefuse

// From builds a Because[R] from whatever the describable can say about the
// refusal; when it describes nothing, the reason stands alone.
func From[R any](reason R, d Describable) Because[R] {
	b := Because[R]{Reason: reason, Cause: nil, What: Refusal{File: "", Line: 0, Subject: "", Why: ""}}

	if what, ok := d.Describe(); ok {
		b.What = what
	}

	return b
}
