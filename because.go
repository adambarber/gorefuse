package gorefuse

import "errors"

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

// Is matches only an error reason: errors.Is on the reason value. A non-error R
// carries no error identity — use errors.As on the carrier type to match its
// domain; nil never bridges domains.
func (b Because[R]) Is(target error) bool {
	if target == nil {
		return false
	}

	if err, ok := any(b.Reason).(error); ok {
		return errors.Is(err, target)
	}

	return false
}

// Refuse builds a carrier without a cause, so callers do not spell the struct
// literal at every refusal site.
func Refuse[R any](reason R, what Refusal) Because[R] {
	return Wrap(reason, what, nil)
}

// Wrap builds a carrier with one cause. errors.Is and errors.As traverse it
// through Unwrap.
func Wrap[R any](reason R, what Refusal, cause error) Because[R] {
	return Because[R]{Cause: cause, What: what, Reason: reason}
}

// Join builds a carrier whose cause is errors.Join of causes. The joined cause
// keeps standard multi-error identity; a marked Cause slot suppresses the
// whole joined tail in prose, JSON, and structured logs.
func Join[R any](reason R, what Refusal, causes ...error) Because[R] {
	return Because[R]{Reason: reason, Cause: errors.Join(causes...), What: what}
}
