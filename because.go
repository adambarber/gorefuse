package gorefuse

// File is the name of the file a refusal speaks about.
type File string

// Line is a one-based line in a File; zero means the whole file.
type Line int

// Subject names the thing inside a file that was refused.
type Subject string

// Why is the fix, as a full sentence ending in punctuation.
type Why string

// Refusal describes what was refused: where it sits, what it is, and the fix.
// Build one only through At.
type Refusal struct {
	File    File
	Subject Subject
	Why     Why
	Line    Line
}

// Because carries one refusal of reason type R. Two domains with different R
// types cannot be mixed: the type parameter is the compile-time wall.
type Because[R any] struct {
	Reason R
	Cause  error
	What   Refusal
}

// Unwrap exposes the Cause to errors.Is and errors.As.
func (b Because[R]) Unwrap() error { return b.Cause }
