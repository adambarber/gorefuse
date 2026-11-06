package gorefuse

// At builds a Refusal: the single construction path. Line 0 means the whole
// file; a subject names the thing inside it; Why states the fix. A Line
// without a File is malformed — file:line is one frame — so At refuses it.
func At(f File, line Line, subject Subject, why Why) Refusal {
	if f == "" && line != 0 {
		panic("gorefuse: At refuses a Line without a File; file:line is one frame")
	}

	return Refusal{File: f, Subject: subject, Why: why, Line: line}
}

// Describable is anything that can say what was refused and where. The bool
// reports whether a description was available at all.
type Describable interface {
	Describe() (Refusal, bool)
}
