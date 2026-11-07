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

// At builds a Refusal: the single construction path. Line 0 means the whole
// file; a subject names the thing inside it; Why states the fix. At refuses
// File with a newline or trailing ":digits"; marked slots render Redacted.
func At(f File, line Line, subject Subject, why Why, marks ...Sensitive) Refusal {
	if strings.TrimRightFunc(string(f), unicode.IsSpace) == "" && line != 0 {
		panic("gorefuse: At refuses Line " + strconv.Itoa(int(line)) +
			" without a File (" + strconv.Quote(string(f)) + "); file:line is one frame")
	}

	file := strings.TrimSpace(string(f))
	if strings.ContainsAny(file, "\n") {
		panic("gorefuse: At refuses File " + strconv.Quote(file) +
			"; File contains a newline, which would forge its file:line frame")
	}

	if trailingDigits.MatchString(file) {
		panic("gorefuse: At refuses File " + strconv.Quote(file) +
			"; trailing :digits would forge its file:line frame")
	}

	if line < 0 {
		panic("gorefuse: At refuses negative Line " + strconv.Itoa(int(line)) +
			" for File " + strconv.Quote(string(f)) + "; lines are one-based")
	}

	sensitive := Sensitive(0)
	for _, m := range marks {
		sensitive |= m
	}

	return Refusal{File: f, Subject: subject, Why: why, Line: line, Sensitive: sensitive}
}

// Describable is anything that can say what was refused and where; the bool
// reports whether a description was available. Describe must not panic —
// From does not recover.
type Describable interface {
	Describe() (Refusal, bool)
}
