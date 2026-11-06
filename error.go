package gorefuse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Error renders the refusal in a stable order: frame prefix, reason, subject,
// cause tail, then the fix as a separate sentence. Trailing whitespace of any
// kind is trimmed before the punctuation check; an embedded newline passes
// through, so a multi-line reason keeps its shape under the grep-able frame.
func (b Because[R]) Error() string {
	msg := strings.TrimRightFunc(b.What.frame()+fmt.Sprint(b.Reason), unicode.IsSpace)
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
// An empty or whitespace-only message gets no separator — the fix stands
// alone. The check decodes the last rune and honors ASCII and ideographic
// terminals (。！？), so a multi-byte terminal still joins with a space.
func (r Refusal) sentenceTail(msg string) string {
	if r.Why == "" {
		return ""
	}

	if msg == "" {
		return ""
	}

	term, _ := utf8.DecodeLastRuneInString(msg)
	if strings.ContainsRune(".!?。！？", term) {
		return " "
	}

	return ". "
}
