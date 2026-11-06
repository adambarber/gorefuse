package gorefuse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Error renders the refusal in a stable order: frame prefix, reason, subject,
// cause tail, then the fix as a separate sentence. Trailing whitespace is
// trimmed from the assembled message before the punctuation check, and a
// part whose trimmed contribution is empty (a whitespace-only cause, say)
// is skipped entirely — no connector, no join. An embedded newline passes
// through, so a multi-line reason keeps its shape under the grep-able frame.
func (b Because[R]) Error() string {
	msg := b.What.frame() + fmt.Sprint(b.Reason)
	if b.What.Subject != "" {
		msg += ": " + string(b.What.Subject)
	}

	if cause := strings.TrimRightFunc(causeText(b.Cause), unicode.IsSpace); cause != "" {
		msg += ": " + cause
	}

	msg = strings.TrimRightFunc(msg, unicode.IsSpace)

	return msg + b.What.sentenceTail(msg) + string(b.What.Why)
}

// causeText is the cause's message, or "" when there is no cause.
func causeText(cause error) string {
	if cause == nil {
		return ""
	}

	return cause.Error()
}

// frame renders the file:line prefix, with the trailing separator, or "".
// The file:line pairing holds for At-constructed (and From-of-honest-Describe)
// values; a directly constructed Refusal with a Line and no File drops the
// line silently — that is the caller's contract breach, not the renderer's.
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
