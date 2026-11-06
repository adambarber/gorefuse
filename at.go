package gorefuse

import "strconv"

// At builds a Refusal: the single construction path. Line 0 means the whole
// file; a subject names the thing inside it; Why states the fix. A malformed
// Line (no File, or negative) is refused with the offending values echoed.
func At(f File, line Line, subject Subject, why Why) Refusal {
	if f == "" && line != 0 {
		panic("gorefuse: At refuses Line " + strconv.Itoa(int(line)) +
			" without a File; file:line is one frame")
	}

	if line < 0 {
		panic("gorefuse: At refuses negative Line " + strconv.Itoa(int(line)) +
			" for File " + strconv.Quote(string(f)) + "; lines are one-based")
	}

	return Refusal{File: f, Subject: subject, Why: why, Line: line}
}

// Describable is anything that can say what was refused and where; the bool
// reports whether a description was available. Describe must not panic —
// From does not recover.
type Describable interface {
	Describe() (Refusal, bool)
}
