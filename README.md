# gorefuse

A typed refusal vocabulary for Go: one way to say **why** a thing was refused — with a `file:line` frame, a reason type as identity, and a guaranteed render shape.

```bash
go get github.com/adambarber/gorefuse
```

```go
package main

import (
	"errors"
	"fmt"

	"github.com/adambarber/gorefuse"
)

type Reason string

const Denied Reason = "denied"

func (r Reason) Error() string { return string(r) }

func main() {
	var err error = gorefuse.Because[Reason]{
		Reason: Denied,
		What: gorefuse.At("server.go", 42, gorefuse.Subject("listen :8080"),
			gorefuse.Why("bind the port before serving")),
	}
	fmt.Println(err)

	// Switch on the refusal's identity, not its words:
	wrapped := fmt.Errorf("starting server: %w", err)
	fmt.Println(errors.Is(wrapped, Denied))
}

// Output:
// Error: server.go:42: denied: listen :8080. bind the port before serving
// true
```

At a glance — use gorefuse if you want:

- a refusal carried as an `error` with a typed reason as its identity: `Because[R]`
- a reason you can switch on: `errors.Is(err, Denied)` matches through any wrapping (when `R` is an error type, matching follows `errors.Is` semantics on the reason value)
- one constructor that builds the refusal and guarantees the frame: `At`
- named slots instead of prose fragments: `Refusal` with `File`, `Line`, `Subject`, `Why`
- refusals described by values, nil-safe: `Describable` and `From[R]`
- a render you can parse and grep from v0.1.1 onward, byte for byte: the render contract below. The verdict label was added in v0.1.1 — the Status section owns that change.

## Why

`fmt.Errorf("...: %w", err)` gives you wrapping but no guaranteed render shape and no `file:line` frame that stays grep-able; `errors.Is` finds a cause but says nothing about what the error *looks like* in a log. gorefuse is the missing half: a refusal is a small struct — where, what, why — with a reason type that makes two domains two distinct types at compile time, rendered by one contract.

What it does **not** do: no structured/JSON rendering yet (roadmap), no retry or classification policy, no terminal-color output.

## The render contract

Render order is fixed: **verdict label → frame prefix → reason → subject → cause tail → fix sentence**. Every render begins with the `Error: ` verdict label — nothing about a log line should leave the reader guessing whether it marks a failure — and the label obeys the same skip-if-empty rule as every slot: the zero `Because[R]` for a string-kinded reason renders `""`, label and all, and is still a valid `error` with identity, because the reason type is the identity, not the words. Every contribution slot — reason, subject, cause, `File`, `Why` — obeys one rule: trimmed once, and empty-after-trim means the slot contributes nothing, no connector.

- **Well-formed UTF-8 on one line, always.** Invalid UTF-8 and the ESC byte become U+FFFD, so a hostile slot cannot break a UTF-8 log encoder or forge terminal output; `\r`, `\n`, `\v`, `\f`, NUL, U+2028 and U+2029 become spaces. Pinned by `scrub_test.go`, `verdict_test.go` and `edge_test.go`.
- **Marked-sensitive slots render as `[redacted]`.** gorefuse cannot know what a secret is — the caller classifies; a slot marked with `Sensitive` renders as the fixed token `[redacted]` regardless of content, a marked `Cause` renders nothing at all (the cause is a tail, so it has no single-token form), and the token itself is never scrubbed or transformed. The marks cover Subject, Why, and Cause — **never `Reason`**, which is of arbitrary caller type `R` and outside what gorefuse can redact: never put secret material in `Reason`. This is defense-in-depth, not the primary wall. Pinned by `redact_test.go`.
- **The `[redacted]` token is not authenticated.** Unmarked content containing the literal `[redacted]` renders identically to a real redaction; redaction is a rendering decision of this library, not a cryptographic guarantee. Pinned by `TestRedactedTokenIsByteExact`.
- **Grep-able frame.** Every render has exactly one `Error:` at line start: line breaks in any slot are flattened to spaces, so a hostile reason cannot forge a second verdict on a new line. Bidi overrides such as U+202E pass through deliberately — they are valid Unicode and a grapheme-level display concern, not something a refusal library sanitizes.
- **Unambiguous frame on every path.** `At` accepts every `File` name. The render path substitutes a fixed `<invalid-frame>` placeholder (with no `Line`) for a name that would forge a `file:line` pairing — a line break, trailing `:digits`, or a dangling trailing `:` — so even a hand-built `Refusal` cannot produce a false pairing. A non-positive `Line` never pairs. Pinned by `at_test.go` and `verdict_test.go`.
- **`At` panics only on caller mistakes.** The three conditions: a negative `Line`, a whitespace-only `File` with a non-zero `Line`, and unknown `Sensitive` bits. Hostile file names are data, not programmer errors. A `File` ending in `:` followed by non-ASCII digits is prose by contract: only ASCII `:digits` is the pairing grammar. `At` stores the `File` trimmed; the render scrub is the only transformation. Pinned by `at_test.go`.
- **errors.Is / errors.As unwrap to `Cause`** via `Unwrap`. Pinned by `because_test.go`.
- **errors.Is against the reason.** `Because[R]` carries `Is(target error) bool` on the carrier, so `errors.Is(err, someReason)` matches through any wrapping. Only an error-kinded `R` has a reason identity: matching follows `errors.Is` on the reason value. A non-error `R` matches nothing through `Is` — use `errors.As` on the carrier type; a nil target never matches and a `Because[otherR]` never bridges domains. Pinned by `is_test.go`.
- **A Because used as a cause carries its own verdict label.** `Because{Cause: innerBecause}` renders `Error: missing: Error: refused` — the nested label is never stripped and never sits at line start, so line-start grep (`^Error:`) still sees exactly one verdict per render. Pinned by `TestNestedBecauseCauseCarriesItsOwnLabel`.
- **Composition contract.** A consuming library puts its location prose (`unknown field "extra" in /path/app.json, line 1`) in the Reason and the remediation sentence in the Why; the fixed render order delivers the documented model shape byte for byte, and the grep-able `path:line:` frame is the same refusal's machine-parseable form. Pinned by `composition_test.go`.
- **The compile-time wall.** The type parameter `R` is the reason domain: `errors.As` never matches `Because[otherReason]` against a `Because[testReason]`, and a tag-gated cross-domain assignment test (`compilecheck_test.go`) fails the build if two domains could ever collapse. Pinned by `guarantee_test.go`.
- **The render path is fuzzed.** `fuzz_test.go` builds hand-built, `From`, and `At` carriers from every input and asserts each render is valid UTF-8, idempotent, free of unscrubbed ESC bytes and line breaks, and has one line-start verdict; scrub/redaction/verdict-forgery bytes seed the corpus.

## Tradeoffs

The render is prose only. A structured form is the first roadmap item, and it must derive from the same described parts rather than re-render — the prose contract stays the human path.

## Zero dependencies

stdlib only; reflection is used only to detect nil describables in `From`.

## Status

API is stable modulo the pre-1.0 conventions of an early library. The render contract has changed twice pre-1.0: v0.1.1 added the `Error: ` verdict label at the head of every non-empty render; the next release flattens every line break to a space (so one render is one verdict line) and stops `At` from panicking on file names — ambiguous names now render `<invalid-frame>`. Both are breaking changes to exact-byte expectations, accepted because the one-verdict guarantee and total file-name handling are the stronger contracts. `Is` is likewise narrowed to error-kinded reasons; a non-error reason domain is matched with `errors.As` on the carrier type. Exact-byte pins in this tree were updated in the same commits, not deleted; the contract is byte-stable from the second change onward.

## Roadmap

- **Structured rendering of refusals**: a refusal should be serializable into a machine-readable format. The prose render stays the human path; serialization must not fork the render contract (stable order, whitespace rules), so the structured form derives from the same described parts, not re-rendered. Structured rendering may also carry per-part sensitivity for `Reason` once a schema exists — today the marks cannot cover it.