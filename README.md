# gorefuse

A typed refusal vocabulary for Go: carry a refusal as an `error` with a reason identity, a `file:line` frame, and one stable prose render.

`fmt.Errorf` gives wrapping but no guaranteed render shape, no frame, and no reason you can switch on. gorefuse gives you one carrier, a few plain constructors, and one render contract.

## Install

```sh
go get github.com/adambarber/gorefuse
```

Go 1.21 or later; no dependencies outside the standard library.

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
	b := gorefuse.Because[Reason]{
		Reason: Denied,
		What: gorefuse.At("server.go", 42,
			gorefuse.Subject("listen :8080"),
			gorefuse.Why("bind the port before serving")),
	}

	fmt.Println(b.Verdict())
	fmt.Println(errors.Is(fmt.Errorf("starting server: %w", b), Denied))
}
// Output:
// Error: server.go:42: denied: listen :8080. bind the port before serving
// true
```

The reason type is the identity: two libraries' reasons are two incompatible types, and `errors.Is` matches through `Because` and ordinary wrapping when the reason is an error. `Error()` is the unlabelled, composable content; `Verdict()` adds the `Error:` label for log emission.

## At a glance

- Use `Because[R]` when a refusal needs a domain-owned identity that a caller can switch on.
- Use `Error()` for composable prose and `Verdict()` when emitting a single labelled log line.
- Use `errors.Is` against an error-kinded reason; use `errors.As` on the carrier when you need the frame or a non-error reason domain.
- Use `At(file, line, subject, why)` to build the frame; it is the one constructor that keeps it well formed.
- Use `Describable` and `From[R]` when the domain says what it refused and the caller supplies the reason.
- Use `Sensitive` marks on `Subject`, `Why`, or `Cause` to redact those slots in the render. `File` and `Reason` are never redacted — keep secrets out of them.
- Use `Unwrap` to keep an underlying cause reachable through `errors.Is` and `errors.As`.

## Why

A sentinel plus a custom type gives identity; `fmt.Errorf("%w", err)` gives wrapping; a hand-written struct gives fields. What none of them gives is a single documented shape for what a refusal looks like when it reaches a log.

gorefuse is a small, opinionated answer: a refusal has a where (`File`, `Line`), a what (`Subject`), and a fix (`Why`); the reason is a type parameter, so two domains cannot collapse; and the render is fixed and fuzz-pinned. It is a house style, not a new capability — if plain sentinels and lowercase error strings already fit your project, you do not need it.

## The render contract

Render order is fixed: **verdict label → frame → reason → subject → cause tail → fix sentence**. Every slot is scrubbed and then trimmed; an empty-after-scrub slot contributes nothing and no connector.

- **One verdict at line start in `Verdict()`.** `Error()` carries the same bytes without the label, so wrapping reads `starting server: server.go:42: ...`. Invalid UTF-8 and ESC become U+FFFD; `\r`, `\n`, `\v`, `\f`, NUL, U+0085, U+2028 and U+2029 become spaces. A hostile slot cannot forge a second line-start `Error:`.
- **Grep-able frame.** `At` accepts every file name; a name that would forge a `file:line` pairing renders as `<invalid-frame>` with no line in prose, JSON, and slog. A non-positive line never pairs. A `File` ending in `:` followed by non-ASCII digits is prose by contract; only ASCII `:digits` is the pairing grammar.
- **Marked slots are redacted.** A marked `Subject` or `Why` renders as `[redacted]`; a marked `Cause` renders nothing. Marked cause slots still unwrap. Unknown sensitive bits fail closed on every construction path.
- **Reason identity, not reason words.** `Is` matches only an error-kinded reason; a non-error reason domain is matched with `errors.As` on the carrier. A nil target never matches and another reason type never bridges domains.
- **The zero value is an empty error.** A zero `Because[R]` renders `""` and is still a non-nil `error`.
- **Pinned prose shape.** Exact-byte tests and a fuzz target pin validity, idempotence, a single line-start verdict, no trailing connector, and slot redaction. Order and shape are stable; bytes may change before v1.0.

## Error and failure contract

- `From[R]` is nil-safe: untyped nil or a typed nil pointer/map/slice/func/chan yields the reason alone.
- `Describe` must not panic; `From` does not recover.
- A `Cause` whose `Error()` panics is the caller's defect; gorefuse does not recover.
- `At` panics only on caller mistakes: negative line, line without a file, or unknown sensitive bits. Hostile file names are data.

## Construction and composition

Use `Refuse(reason, what)`, `Wrap(reason, what, cause)`, or `Join(reason, what, causes...)` instead of spelling the struct literal. `Join` stores `errors.Join(causes...)` as the cause, so `errors.Is` and `errors.As` traverse every member through the standard library. A marked cause suppresses the whole joined tail in prose, JSON, and structured logs.

## What redaction does and does not do

`Sensitive` is a classification **you** supply. A marked `Subject` or `Why` renders as `[redacted]`; a marked `Cause` renders nothing but still unwraps. The token is a stable rendering marker, not an authenticated seal: unmarked text that literally contains `[redacted]` renders identically.

```go
b := gorefuse.Because[Reason]{
	Reason: Denied,
	What:   gorefuse.At("server.go", 42, "[redacted]", ""),
}
fmt.Println(b.Verdict())
// Output:
// Error: server.go:42: denied: [redacted]
```

To prove how a value was classified, match the carrier with `errors.As` and inspect `What.Sensitive`; the string alone cannot tell you. `fmt`'s default verbs, `%#v`, `MarshalJSON`, and `slog.LogValuer` all apply the same slot decisions, including on a bare `Refusal`, so a marked subject or fix becomes `[redacted]` everywhere.

## Comparison

- **vs `fmt.Errorf`:** typed reason identity and a fixed, parseable render; you give up free-form lowercase error fragments.
- **vs `cockroachdb/errors`:** no stack traces, network encoding, or Sentry reports; standard-library-only and byte-pinned instead.
- **vs `rotisserie/eris`:** JSON and `slog` output carry the same slot decisions; no formatter plugins or stack encoding.
- **vs `joomcode/errorx`:** the reason is a type parameter rather than a trait system; smaller surface, fewer capabilities.
- **vs `shiwano/errdef`:** definition identity and redaction across encoders are close prior art; gorefuse adds the refusal-as-remediation framing and a byte-pinned prose contract.
- **vs `upspin.io/errors`:** the fixed `file:line` frame and reason vocabulary are similar in spirit; gorefuse makes the reason domain a compile-time type parameter.

## Caveats

- The label is added by `Verdict()`, not by `Error()`. Wrapping reads `starting server: server.go:42: ...`; call `Verdict()` when emitting a single labelled log line.
- Redaction is a rendering decision, not a security boundary: only marked slots are hidden, and the marker is guessable.
- Machine-readable output is limited to `encoding/json` and `log/slog`; there is no protobuf or Sentry integration.

## Examples

[`examples/`](examples/README.md) has standalone recipes you can copy into a project: basic refusal, reason matching, custom describers, redaction, cause wrapping, and hostile frames. Each was run with `go run` and its output recorded; the same scenarios are published as executable `Example...` functions on [pkg.go.dev](https://pkg.go.dev/github.com/adambarber/gorefuse).

## Status

Pre-1.0. Current release: `v0.2.1`. It adds machine-readable output, plain constructors, a documented multi-cause story, and lowers the Go floor to 1.21.

[![PkgGoDev](https://pkg.go.dev/badge/github.com/adambarber/gorefuse)](https://pkg.go.dev/github.com/adambarber/gorefuse) [![Release](https://img.shields.io/github/release/adambarber/gorefuse.svg?color=%23007ec6)](https://github.com/adambarber/gorefuse/releases/latest) [![test](https://github.com/adambarber/gorefuse/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/adambarber/gorefuse/actions/workflows/test.yml) [![lint](https://github.com/adambarber/gorefuse/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/adambarber/gorefuse/actions/workflows/lint.yml) [![coverage](https://raw.githubusercontent.com/adambarber/gorefuse/badges/.badges/main/coverage.svg)](.testcoverage.yml)

The coverage badge is written by CI to the `badges` branch; it appears after that branch exists on the remote and the first `main` run completes.

## License

[MIT](LICENSE)
