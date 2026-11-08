# gorefuse

A typed refusal vocabulary for Go: one way to say **why** a thing was refused — with a `file:line` frame, a reason type as identity, and a guaranteed render shape.

```bash
go get github.com/adambarber/gorefuse
```

```go
package main

import (
	"fmt"

	"github.com/adambarber/gorefuse"
)

type Reason string

const Denied Reason = "denied"

func main() {
	var err error = gorefuse.Because[Reason]{
		Reason: Denied,
		What: gorefuse.At("server.go", 42, gorefuse.Subject("listen :8080"),
			gorefuse.Why("bind the port before serving")),
	}
	fmt.Println(err)
}

// Output:
// Error: server.go:42: denied: listen :8080. bind the port before serving
```

At a glance — use gorefuse if you want:

- a refusal carried as an `error` with a typed reason as its identity: `Because[R]`
- one constructor that builds the refusal and guarantees the frame: `At`
- named slots instead of prose fragments: `Refusal` with `File`, `Line`, `Subject`, `Why`
- refusals described by values, nil-safe: `Describable` and `From[R]`
- a render you can parse and grep, byte for byte, forever: the render contract below (v0.1.1 added the `Error: ` verdict label — a breaking change to exact-byte expectations, accepted pre-1.0)

## Why

`fmt.Errorf("...: %w", err)` gives you wrapping but no guaranteed render shape and no `file:line` frame that stays grep-able; `errors.Is` finds a cause but says nothing about what the error *looks like* in a log. gorefuse is the missing half: a refusal is a small struct — where, what, why — with a reason type that makes two domains two distinct types at compile time, rendered by one contract.

What it does **not** do: no structured/JSON rendering yet (roadmap), no retry or classification policy, no terminal-color output.

## The render contract

Render order is fixed: **verdict label → frame prefix → reason → subject → cause tail → fix sentence**. Every render begins with the `Error: ` verdict label — nothing about a log line should leave the reader guessing whether it marks a failure — and the label obeys the same skip-if-empty rule as every slot: the zero `Because[R]` for a string-kinded reason renders `""`, label and all, and is still a valid `error` with identity, because the reason type is the identity, not the words. Every contribution slot — reason, subject, cause, `File`, `Why` — obeys one rule: trimmed once, and empty-after-trim means the slot contributes nothing, no connector.

- **Well-formed UTF-8, always.** Invalid UTF-8 and the ESC byte become U+FFFD, so a hostile slot cannot break a UTF-8 log encoder or forge terminal output; each `\r` becomes a space. Pinned by `scrub_test.go` and `edge_test.go`.
- **Marked-sensitive slots render as `[redacted]`.** gorefuse cannot know what a secret is — the caller classifies; a slot marked with `Sensitive` renders as the fixed token `[redacted]` regardless of content, a marked `Cause` renders nothing at all (the cause is a tail, so it has no single-token form), and the token itself is never scrubbed or transformed. The marks cover Subject, Why, and Cause — **never `Reason`**, which is of arbitrary caller type `R` and outside what gorefuse can redact: never put secret material in `Reason`. This is defense-in-depth, not the primary wall. Pinned by `redact_test.go`.
- **The `[redacted]` token is not authenticated.** Unmarked content containing the literal `[redacted]` renders identically to a real redaction; redaction is a rendering decision of this library, not a cryptographic guarantee. Pinned by `TestRedactedTokenIsByteExact`.
- **Grep-able frame.** Embedded newlines in a slot pass through, so a multi-line reason keeps its shape under the `file:line` prefix; bidi overrides such as U+202E pass through deliberately — they are valid Unicode and a grapheme-level display concern, not something a refusal library sanitizes.
- **Unambiguous frame on every path.** `At` panics on a `File` containing a newline or ending in trailing `:digits`; the render path independently renders a fixed `<invalid-frame>` placeholder (with no `Line`) for any `File` that would forge the pairing — and for a colon-terminated `File` or unknown `Sensitive` bits — so even a hand-built `Refusal` — e.g. through `From`/`Describe` — cannot produce a false `file:line` pairing. A non-positive `Line` never pairs. Pinned by `at_test.go`.
- **`At` panics on a forged or malformed frame.** The five conditions: a newline-bearing `File`, a trailing `:digits` `File`, a colon-terminated `File` (a dangling `:` in the frame), a negative `Line`, and a whitespace-only `File` with a non-zero `Line` — plus unknown `Sensitive` bits. A `File` ending in `:` followed by non-ASCII digits is prose by contract: only ASCII `:digits` is the pairing grammar. `At` stores the `File` trimmed, so the stored value equals the rendered frame modulo the scrub stage (a stored dirty byte renders as U+FFFD; storage keeps the caller's value). Pinned by `at_test.go`.
- **errors.Is / errors.As unwrap to `Cause`** via `Unwrap`. Pinned by `because_test.go`.
- **errors.Is against the reason.** `Because[R]` carries an `Is(target error) bool`, and it is on the carrier itself so `errors.Is(err, someReason)` matches through `Because` and any wrapping. When `R` is an error type, matching follows `errors.Is` semantics on the reason value; when `R` is not an error, matching is any-equality, so a non-error reason matches only an equal non-error target. A nil target never matches, and a `Because[otherR]` never bridges domains — the compile-time wall holds through `Is` too. Pinned by `is_test.go`.
- **The compile-time wall.** The type parameter `R` is the reason domain: `errors.As` never matches `Because[otherReason]` against a `Because[testReason]`, and a tag-gated cross-domain assignment test (`compilecheck_test.go`) fails the build if two domains could ever collapse. Pinned by `guarantee_test.go`.

## Tradeoffs

The render is prose only. A structured form is the first roadmap item, and it must derive from the same described parts rather than re-render — the prose contract stays the human path.

## Zero dependencies

stdlib only; reflection is used only to detect nil describables in `From`.

## Status

API is stable modulo the pre-1.0 conventions of an early library. The render contract changed once pre-1.0: v0.1.1 added the `Error: ` verdict label at the head of every non-empty render — a breaking change to exact-byte expectations, accepted because nothing in v0.1.0 marked the render as an error and the label is what log readers lock onto. Exact-byte pins in this tree were updated in the same commits, not deleted.

## Roadmap

- **Structured rendering of refusals**: a refusal should be serializable into a machine-readable format. The prose render stays the human path; serialization must not fork the render contract (stable order, whitespace rules), so the structured form derives from the same described parts, not re-rendered. Structured rendering may also carry per-part sensitivity for `Reason` once a schema exists — today the marks cannot cover it.