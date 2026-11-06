package gorefuse

import "fmt"

// Error renders the refusal in a stable order: frame prefix, reason, subject,
// cause tail, then the fix as a separate sentence.
func (b Because[R]) Error() string {
	msg := b.What.frame() + fmt.Sprint(b.Reason)
	if b.What.Subject != "" {
		msg += ": " + string(b.What.Subject)
	}

	if b.Cause != nil {
		msg += ": " + b.Cause.Error()
	}

	return msg + b.What.sentenceTail(msg) + string(b.What.Why)
}

// frame renders the file:line prefix, with the trailing separator, or "".
func (r Refusal) frame() string {
	if r.File == "" {
		return ""
	}

	if r.Line == 0 {
		return string(r.File) + ": "
	}

	return fmt.Sprintf("%s:%d: ", r.File, r.Line)
}

// sentenceTail returns the punctuation joining the fix to what precedes it.
func (r Refusal) sentenceTail(msg string) string {
	if r.Why == "" {
		return ""
	}

	last := msg[len(msg)-1]
	if last == '.' || last == '!' || last == '?' {
		return " "
	}

	return ". "
}
