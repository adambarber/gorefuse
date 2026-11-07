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
// Build one only through At; the render path guards the frame even for
// hand-built structs. A slot marked in Sensitive renders as Redacted.
type Refusal struct {
	File      File
	Subject   Subject
	Why       Why
	Line      Line
	Sensitive Sensitive
}

// Because carries one refusal of reason type R; the type parameter is the
// compile-time wall (zero value renders ""; see TestEmptyBecauseRendersEmpty).
// A Cause that panics in Error is the caller's defect — no recovery.
type Because[R any] struct {
	Reason R
	Cause  error
	What   Refusal
}

// Unwrap exposes the Cause to errors.Is and errors.As.
func (b Because[R]) Unwrap() error { return b.Cause }
