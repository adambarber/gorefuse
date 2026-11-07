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
// slot, so a hostile slot cannot forge terminal output. A slot marked
// sensitive renders as the fixed token Redacted, never its content — the
// caller classifies what is sensitive; the renderer guarantees only the
// marked case (see Sensitive). A cause marked CauseSensitive renders nothing.
// Bidi overrides such as U+202E pass through: they are valid Unicode and a
// grapheme-level display concern, not something the library sanitizes.
func (b Because[R]) Error() string {
	var msg string

	if frame := b.What.frame(); frame != "" {
		msg = frame
	}

	for _, part := range []string{
		fmt.Sprint(b.Reason),
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

	return msg + b.What.sentenceTail(msg) +
		scrub(strings.TrimSpace(string(markedSensitiveWhy(b.What.Why, WhySensitive, b.What.Sensitive))))
}

// scrub guarantees one slot's rendered text is well-formed UTF-8 with no raw
// carriage return and no ESC byte: invalid UTF-8 and ESC become U+FFFD, each
// \r becomes a space. Embedded newlines pass through untouched.
func scrub(s string) string {
	s = strings.ToValidUTF8(s, "\uFFFD")
	s = strings.ReplaceAll(s, "\x1b", "\uFFFD")

	return strings.ReplaceAll(s, "\r", " ")
}

// causeTail renders the cause's message, or "" when there is no cause or the
// cause is marked sensitive: a marked cause contributes nothing, as if absent.
// The cause is a tail, not a slot, so it has no single-token redacted form.
func causeTail[R any](b Because[R]) string { //nolint:forbidigo // R must accept any reason kind, as in from.go
	if b.Cause == nil || b.What.Sensitive&CauseSensitive != 0 {
		return ""
	}

	return b.Cause.Error()
}

// frame renders the file:line prefix, or "". The File is trimmed first: empty-after-trim means no frame — the uniform
// skip-if-empty rule applies to File like every other slot. The guard is total: on EVERY construction path, a File that
// could forge the file:line pairing — newline-bearing or trailing ":digits" — renders the fixed placeholder invalidFrame
// with no Line, so no path can emit a false file:line pairing. At refuses such input up front; this render-time backstop
// covers hand-built Refusal structs (e.g. through From/Describe). The same placeholder covers a colon-terminated File
// (a dangling ":" in the frame) and unknown Sensitive bits: the render path refuses what At refuses. A non-positive
// Line never pairs: it is skipped.
func (r Refusal) frame() string {
	file := scrub(strings.TrimSpace(string(r.File)))
	if file == "" {
		return ""
	}

	if strings.ContainsAny(file, "\n") || trailingDigits.MatchString(file) || trailingColon.MatchString(file) ||
		r.Sensitive&^sensitiveMask != 0 {
		return invalidFrame
	}

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
