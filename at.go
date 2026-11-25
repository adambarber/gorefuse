package gorefuse

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// trailingDigits matches a File whose rendered frame would end in ":digits",
// forging the file:line pairing that log parsers lock onto.
var trailingDigits = regexp.MustCompile(`:[0-9]+$`)

// trailingColon matches a File ending in a bare colon: its frame carries a
// dangling separator ("a.go::7") no documented frame shape blesses.
var trailingColon = regexp.MustCompile(`:$`)

// invalidFrame is the placeholder the render path substitutes for a File that
// would forge a file:line pairing; it never appears as a real file name.
const invalidFrame = "<invalid-frame>"

// At builds a Refusal; Line 0 means the whole file. Every File name is accepted
// (ambiguous ones render invalidFrame); a negative Line, a Line without a File,
// and unknown Sensitive bits panic as caller mistakes.
func At(f File, line Line, subject Subject, why Why, marks ...Sensitive) Refusal {
	if strings.TrimRightFunc(string(f), unicode.IsSpace) == "" && line != 0 {
		panic("gorefuse: At refuses Line " + strconv.Itoa(int(line)) +
			" without a File (" + strconv.Quote(string(f)) + "); file:line is one frame")
	}

	file := strings.TrimSpace(string(f))

	if line < 0 {
		panic("gorefuse: At refuses negative Line " + strconv.Itoa(int(line)) +
			" for File " + strconv.Quote(string(f)) + "; lines are one-based")
	}

	sensitive := Sensitive(0)

	for _, m := range marks {
		if unknownSensitive(m) {
			panic("gorefuse: At refuses Sensitive marks " + strconv.Itoa(int(m)) +
				"; unknown bits outside SubjectSensitive|WhySensitive|CauseSensitive")
		}

		sensitive |= m
	}

	return Refusal{File: File(file), Subject: subject, Why: why, Line: line, Sensitive: sensitive}
}

// Describable is anything that can say what was refused and where; the bool
// reports whether a description was available. Describe must not panic —
// From does not recover.
type Describable interface {
	Describe() (Refusal, bool)
}
