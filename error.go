package gorefuse

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Error renders the refusal in a stable order: frame prefix, reason, subject,
// cause tail, then the fix as a separate sentence. Every contribution slot —
// reason, subject, cause, File, and Why — obeys one rule: the slot's trimmed
// text is computed once, and empty-after-trim means the slot contributes
// nothing, no connector. An embedded newline passes through, so a multi-line
// reason keeps its shape under the grep-able frame. The rendered output is
// guaranteed well-formed UTF-8 with no raw carriage returns: invalid UTF-8
// bytes become U+FFFD (strings.ToValidUTF8) so downstream UTF-8 log encoders
// never break, and each \r becomes a space — a raw \r overwrites terminal
// lines and has no meaning a refusal wants to preserve. Escape sequences are
// neutralized at the byte level: the ESC byte (\x1b) becomes U+FFFD in every
// slot, so a hostile slot cannot forge terminal output. Bidi overrides such as
// U+202E pass through: they are valid Unicode and a grapheme-level display
// concern, not something the library sanitizes.
func (b Because[R]) Error() string {
	var msg string

	if frame := b.What.frame(); frame != "" {
		msg = frame
	}

	for _, part := range []string{fmt.Sprint(b.Reason), string(b.What.Subject), causeText(b.Cause)} {
		if part = scrub(strings.TrimSpace(part)); part != "" {
			if msg != "" {
				msg += ": "
			}

			msg += part
		}
	}

	return msg + b.What.sentenceTail(msg) + scrub(strings.TrimSpace(string(b.What.Why)))
}

// scrub guarantees one slot's rendered text is well-formed UTF-8 with no raw
// carriage return and no ESC byte: invalid UTF-8 and ESC become U+FFFD, each
// \r becomes a space. Embedded newlines pass through untouched.
func scrub(s string) string {
	s = strings.ToValidUTF8(s, "\uFFFD")
	s = strings.ReplaceAll(s, "\x1b", "\uFFFD")

	return strings.ReplaceAll(s, "\r", " ")
}

// causeText is the cause's message, or "" when there is no cause.
func causeText(cause error) string {
	if cause == nil {
		return ""
	}

	return cause.Error()
}

// frame renders the file:line prefix, or "". The File is trimmed first: empty-after-trim means no frame — the uniform
// skip-if-empty rule applies to File like every other slot. For At-constructed
// values, the file:line pairing is unambiguous: At refuses any File containing
// a newline or ending in ":digits", so a log parser's ^(\S+):\d+: match locks
// onto the true file. A directly constructed Refusal with a Line and no File
// drops the line silently — that is the caller's contract breach, not the
// renderer's.
func (r Refusal) frame() string {
	file := scrub(strings.TrimSpace(string(r.File)))
	if file == "" {
		return ""
	}

	if r.Line == 0 {
		return file
	}

	return fmt.Sprintf("%s:%d", file, r.Line)
}

// sentenceTail returns the punctuation joining the fix to what precedes it.
// An empty or whitespace-only message gets no separator — the fix stands
// alone. The check decodes the last rune and honors ASCII and ideographic
// terminals (。！？), so a multi-byte terminal still joins with a space.
func (r Refusal) sentenceTail(msg string) string {
	if strings.TrimSpace(string(r.Why)) == "" {
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
