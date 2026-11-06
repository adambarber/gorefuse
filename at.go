package gorefuse

// At builds a Refusal: the single construction path. Line 0 means the whole
// file; a subject names the thing inside it; Why states the fix.
func At(f File, line Line, subject Subject, why Why) Refusal {
	return Refusal{File: f, Subject: subject, Why: why, Line: line}
}

// Describable is anything that can say what was refused and where. The bool
// reports whether a description was available at all.
type Describable interface {
	Describe() (Refusal, bool)
}
