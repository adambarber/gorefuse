package gorefuse

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Error renders one verdict per line: label, frame, reason, subject, cause,
// fix. Empty-after-scrub slots contribute nothing. Scrub makes the output valid
// UTF-8, replaces ESC and invalid bytes with U+FFFD, and flattens every line
// break (including U+0085, U+2028/U+2029 and NUL) to a space so no slot can forge a
// second line-start verdict. Marked slots render as [redacted] (CauseSensitive
// contributes nothing); bidi overrides pass through as valid Unicode.
func (b Because[R]) Error() string {
	r := b.What.normalize()

	var msg string
	if frame := r.frame(); frame != "" {
		msg = frame
	}

	slots := []struct {
		text slotText
		mark Sensitive
	}{
		{text: slotText(reasonPart(b.Reason)), mark: 0},
		{text: slotText(r.Subject), mark: SubjectSensitive},
		{text: slotText(causeTail(b, r.Sensitive)), mark: 0},
	}
	for _, slot := range slots {
		part := renderSlot(slot.text, slot.mark, r.Sensitive)
		if part == "" {
			continue
		}

		if msg != "" {
			msg += ": "
		}

		msg += string(part)
	}

	why := renderSlot(slotText(r.Why), WhySensitive, r.Sensitive)
	if msg == "" && why == "" {
		return ""
	}

	return "Error: " + msg + sentenceTail(msg, string(why)) + string(why)
}

// scrub guarantees one slot's text is well-formed UTF-8 on a single line:
// invalid bytes and ESC become U+FFFD; line breaks become spaces.
func scrub(s string) string {
	s = strings.ToValidUTF8(s, "\uFFFD")
	s = strings.ReplaceAll(s, "\x1b", "\uFFFD")

	return strings.Map(func(r rune) rune {
		switch r {
		case '\r', '\n', '\v', '\f', 0, '\u0085', '\u2028', '\u2029':
			return ' '
		default:
			return r
		}
	}, s)
}

// reasonPart renders a reason value. A nil error reason contributes nothing;
// any other value is rendered by fmt.
func reasonPart[R any](reason R) string {
	if isNilValue(reason) {
		return ""
	}

	return fmt.Sprint(reason)
}

// causeTail renders the cause's message, or "" when there is no cause or the
// cause is marked sensitive: a marked cause contributes nothing, as if absent.
// The cause is a tail, not a slot, so it has no single-token redacted form.
func causeTail[R any](b Because[R], marks Sensitive) string {
	if b.Cause == nil || marks&CauseSensitive != 0 {
		return ""
	}

	return b.Cause.Error()
}

// frame renders the file:line prefix, or "". An empty File has no frame. The
// guard reads the raw (trimmed) File before scrubbing: a file name carrying a
// line break, trailing ":digits", or a dangling colon would forge a pairing, so
// it renders invalidFrame with no Line. A non-positive Line never pairs.
func (r Refusal) frame() string {
	file := strings.TrimSpace(string(r.File))
	if file == "" {
		return ""
	}

	if strings.ContainsAny(file, "\n\r\v\f\x00\u2028\u2029") || trailingDigits.MatchString(file) ||
		trailingColon.MatchString(file) {
		return invalidFrame
	}

	file = scrub(file)
	if r.Line <= 0 {
		return file
	}

	return fmt.Sprintf("%s:%d", file, r.Line)
}

// sentenceTail returns the punctuation joining the fix to what precedes it.
// An empty fix gets no separator; an empty message stands alone. ASCII and
// ideographic terminals (。！？) join with a space, other ends with ". ".
func sentenceTail(msg, why string) string {
	if why == "" || msg == "" {
		return ""
	}

	term, _ := utf8.DecodeLastRuneInString(msg)
	if strings.ContainsRune(".!?。！？", term) {
		return " "
	}

	return ". "
}
