# gorefuse

A typed refusal vocabulary for Go: carry a refusal as an `error` with a reason type as its identity, a `file:line` frame, and one byte-stable render.

`fmt.Errorf` gives you wrapping but no guaranteed render shape, no frame, and no reason you can switch on. gorefuse gives you one carrier, one constructor, and one render contract.

## Install

```sh
go get github.com/adambarber/gorefuse
```

Go 1.27 or later is required. The package has no dependencies outside the standard library.

## Example

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
		What: gorefuse.At("server.go", 42,
			gorefuse.Subject("listen :8080"),
			gorefuse.Why("bind the port before serving")),
	}

	fmt.Println(err)
	fmt.Println(errors.Is(fmt.Errorf("starting server: %w", err), Denied))
}
// Output:
// Error: server.go:42: denied: listen :8080. bind the port before serving
// true
```

The reason type is the identity: two libraries' reasons are two incompatible types, and `errors.Is` matches through `Because` and ordinary wrapping when the reason is an error.

## At a glance

- `Because[R]` carries one refusal of reason type `R` as an `error`.
- `errors.Is(err, Denied)` switches on the reason; `errors.As(err, &because)` extracts the carrier when you need its frame.
- `At(file, line, subject, why)` is the one construction path that keeps the frame well formed.
- `Describable` and `From[R]` let another package describe what it refused while you supply the reason.
- `gorefuse.Sensitive` marks slots that hold classified content; marked slots render as `[redacted]`.
- `Error()` is byte-stable and fuzz-pinned: valid UTF-8, one line, one line-start verdict.

## Why

The common Go patterns cover part of this. A sentinel plus a custom type gives identity; `fmt.Errorf("%w", err)` gives wrapping; a hand-written struct gives fields. What they do not give is a single, documented shape for what a refusal looks like when it reaches a log.

gorefuse is a small, opinionated answer: a refusal has a where (`File`, `Line`), a what (`Subject`), and a fix (`Why`); the reason is a type parameter, so two domains cannot collapse; and the render is fixed and fuzz-pinned. It is a house style, not a new capability — if you are happy with plain sentinels and lowercase error strings, you do not need it.

## The render contract

Render order is fixed: **verdict label → frame → reason → subject → cause tail → fix sentence**. Every slot — reason, subject, cause, `File`, `Why` — is trimmed once, and an empty slot contributes nothing and no connector.

- **One verdict per line.** Every non-empty render starts with `Error: `. Invalid UTF-8 and the ESC byte become U+FFFD; `\r`, `\n`, `\v`, `\f`, NUL, U+2028 and U+2029 become spaces. A hostile slot cannot forge a second line-start verdict. Bidi overrides pass through deliberately — they are valid Unicode and a display concern, not a refusal-library concern.
- **Grep-able frame.** `At` accepts every `File` name; a name that would forge a `file:line` pairing renders as `<invalid-frame>` with no line, so even a hand-built `Refusal` cannot smuggle a false frame through. A non-positive line never pairs.

```go
err := gorefuse.Because[Reason]{
	Reason: Denied,
	What: gorefuse.At("backup:2024", 42,
		gorefuse.Subject("line one\nline two"),
		gorefuse.Why("do the thing")),
}
fmt.Println(err)
// Output:
// Error: <invalid-frame>: denied: line one line two. do the thing
```

- **`At` panics only on caller mistakes.** Three conditions: a negative `Line`, a `Line` without a `File`, and unknown `Sensitive` bits. Hostile file names are data, not programmer errors. A `File` ending in `:` followed by non-ASCII digits is prose by contract; only ASCII `:digits` is the pairing grammar.
- **Marked-sensitive slots render as `[redacted]`.** The caller classifies what is sensitive; gorefuse guarantees the marked case. A marked `Cause` renders nothing, because a tail has no single-token form. The token itself is never scrubbed. `Reason` is never redacted — it is of arbitrary caller type `R` — so keep secrets out of it.
- **The `[redacted]` token is not authenticated.** Unmarked text that happens to contain `[redacted]` renders identically. Redaction is a rendering decision, not a cryptographic guarantee.
- **Reason identity, not reason words.** `Unwrap` exposes the `Cause` to `errors.Is` and `errors.As`. `Is` matches only an error-kinded reason, following `errors.Is` on the reason value; a non-error reason domain is matched with `errors.As` on the carrier type. A nil target never matches, and `Because[OtherReason]` never bridges domains.
- **The zero value is an empty error.** The zero `Because[R]` for a string-kinded reason renders `""` and is still a non-nil `error`, because the reason type is the identity, not the words.
- **A nested `Because` keeps its own label.** `Because{Cause: innerBecause}` renders two `Error:` labels, but only the outer one sits at line start, so line-start grep still sees one verdict per render.
- **Compile-time wall.** `errors.As` never matches `Because[OtherReason]` against `Because[Reason]`; a tag-gated build test fails if two domains could collapse.
- **Fuzz-pinned.** The fuzz target builds hand-built, `From`, and `At` carriers from arbitrary bytes and asserts each render is valid UTF-8, idempotent, free of unscrubbed ESC bytes and line breaks, and has one line-start verdict.

## Error and failure contract

- `From[R]` is nil-safe: an untyped nil or a typed nil pointer/map/slice/func/chan yields the reason alone. A nil nested in a struct describer is the describer's own contract.
- `Describe` must not panic; `From` does not recover.
- A `Cause` whose `Error()` panics is the caller's defect; gorefuse does not recover.

## Comparison

- **vs `fmt.Errorf`:** typed reason identity and a fixed, parseable render; you give up free-form lowercase error fragments.
- **vs `cockroachdb/errors`:** no stack traces, network encoding, or Sentry reports; stdlib-only and byte-pinned instead.
- **vs `rotisserie/eris`:** no JSON formatter yet; prose rendering is the only output today.
- **vs `joomcode/errorx`:** the reason is a type parameter rather than a trait system; smaller surface, fewer capabilities.

## Caveats

- The `Error:` prefix means a wrapped chain reads `starting server: Error: server.go:42: ...`. Compose for humans and use `errors.Is` for control flow.
- Redaction is a rendering decision, not a security boundary: only marked slots are hidden, and the marker is guessable.
- Structure is prose-only today; machine-readable output is planned and will derive from the same described parts rather than re-render the string.

## Roadmap

A structured, machine-readable form is the first item. It must derive from the same described parts rather than re-render the prose, so the two forms cannot drift. Once a schema exists, it may carry per-part sensitivity for `Reason`, which the current marks cannot cover.

## Status

Pre-1.0. `v0.1.1` added the `Error:` verdict label at the head of every non-empty render. The untagged commit `96cc429` then flattened every line break to a space, made `At` total on file names, and narrowed `Is` to error-kinded reasons. Exact-byte pins were updated in the same commits; the render is byte-stable since `96cc429`.

[![PkgGoDev](https://pkg.go.dev/badge/github.com/adambarber/gorefuse)](https://pkg.go.dev/github.com/adambarber/gorefuse)
[![Go Report Card](https://goreportcard.com/badge/github.com/adambarber/gorefuse)](https://goreportcard.com/report/github.com/adambarber/gorefuse)

## License

[MIT](LICENSE)
