# gorefuse

The typed refusal vocabulary library: one way to say why a thing was
refused. Sister library to goconf (config validation); follows the same
conventions: birdcage policy gate, strictest merged lint config, TDD.

## Status

All batches landed and review defects fixed: `Because[R]` (error, unwraps to
Cause, renders frame prefix → reason → subject → cause tail → fix sentence;
never panics on an empty message; for string-kinded reasons the zero value
renders "" by documented contract; embedded newlines pass through so the frame
stays grep-able; a whitespace-only reason or cause contributes nothing),
`Refusal` with named kinds `File`/`Line`/`Subject`/`Why`, `At()` as the single
constructor, `Describable`, `From[R]` (nil Describable means no description).
Tail join trims trailing whitespace and decodes the last rune, honoring ASCII
and ideographic terminals (。！？). Compile-time domain wall tested two ways:
`errors.As` never matches across reason domains, and a tag-gated assignment
(`compilecheck_test.go`) is proven to fail the build by
`TestCrossDomainAssignmentDoesNotCompile`.

Security hardening landed (2026-11-06, five clusters): the rendered `Error()`
output is guaranteed well-formed UTF-8 with no raw carriage return and no ESC
byte — invalid UTF-8 and \x1b become U+FFFD, each \r becomes a space, while
embedded \n still passes through (grep-ability contract) and bidi overrides
such as U+202E deliberately pass through (documented, probed: visual reversal
only). `At()` refuses a File containing a newline or a trailing `:digits`
pattern, so the file:line pairing is unambiguous for At-constructed values —
a log parser's `^(\S+):\d+:` match locks onto the true file; interior colons
are accepted and documented. A whitespace-only Why joins the skip-if-empty
discipline and renders trimmed. A Cause that panics in Error is the caller's
defect — documented, not armored.

## Backlog (from the parity review, 2026-11-06)

- **Structured/JSON rendering of refusals**: a Refusable must be serializable into a machine-readable format. The prose render stays the human path; serialization must not fork the render contract (stable order, whitespace rules), so the JSON form is derived from the same described parts, not re-rendered.
