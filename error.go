package gorefuse

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Error renders one verdict per line: label, frame, reason, subject, cause,
// fix. Empty-after-trim slots contribute nothing. Scrub makes the output valid
// UTF-8, replaces ESC and invalid bytes with U+FFFD, and flattens every line
// break (including U+2028/U+2029 and NUL) to a space so no slot can forge a
// second line-start verdict. Marked slots render as [redacted] (CauseSensitive
// contributes nothing); bidi overrides pass through as valid Unicode.
func (b Because[R]) Error() string {
	var msg string

	if frame := b.What.frame(); frame != "" {
		msg = frame
	}

	for _, part := range []string{
		reasonPart(b.Reason),
		string(markedSensitive(b.What.Subject, SubjectSensitive, b.What.Sensitive)),
		causeTail(b),
	} {
		if part = scrub(strings.TrimSpace(part)); part != "" {
			if msg != "" {
				msg += ": "
			}

			msg += part
		}
	}

	why := scrub(strings.TrimSpace(string(markedSensitiveWhy(b.What.Why, WhySensitive, b.What.Sensitive))))

	if msg == "" && why == "" {
		return ""
	}

	return "Error: " + msg + b.What.sentenceTail(msg) + why
}

// scrub guarantees one slot's text is well-formed UTF-8 on a single line:
// invalid bytes and ESC become U+FFFD; line breaks become spaces.
func scrub(s string) string {
	s = strings.ToValidUTF8(s, "\uFFFD")
	s = strings.ReplaceAll(s, "\x1b", "\uFFFD")

	return strings.Map(func(r rune) rune {
		switch r {
		case '\r', '\n', '\v', '\f', 0, '\u2028', '\u2029':
			return ' '
		default:
			return r
		}
	}, s)
}

// reasonPart renders a reason value. A nil error reason contributes nothing;
// any other value is rendered by fmt.
func reasonPart[R any](reason R) string {
	v := any(reason)
	if v == nil {
		return ""
	}

	if err, ok := v.(error); ok && err == nil {
		return ""
	}

	return fmt.Sprint(v)
}

// causeTail renders the cause's message, or "" when there is no cause or the
// cause is marked sensitive: a marked cause contributes nothing, as if absent.
// The cause is a tail, not a slot, so it has no single-token redacted form.
func causeTail[R any](b Because[R]) string {
	if b.Cause == nil || b.What.Sensitive&CauseSensitive != 0 {
		return ""
	}

	return b.Cause.Error()
}

// frame renders the file:line prefix, or "". An empty File has no frame. The
// guard reads the raw (trimmed) File before scrubbing: a file name carrying a
// line break, trailing ":digits", or a dangling colon would forge a pairing, so
// it renders invalidFrame with no Line. This is the backstop for hand-built
// Refusals; At no longer panics on file names. Unknown Sensitive bits use the
// same placeholder, and a non-positive Line never pairs.
func (r Refusal) frame() string {
	file := strings.TrimSpace(string(r.File))
	if file == "" {
		return ""
	}

	if strings.ContainsAny(file, "\n\v\f\x00\u2028\u2029") || trailingDigits.MatchString(file) ||
		trailingColon.MatchString(file) || r.Sensitive&^sensitiveMask != 0 {
		return invalidFrame
	}

	file = scrub(file)
	if r.Line <= 0 {
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
