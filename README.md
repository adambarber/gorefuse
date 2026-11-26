# gorefuse

A typed refusal vocabulary for Go: carry a refusal as an `error` with a reason identity, a `file:line` frame, and one byte-stable render.

`fmt.Errorf` gives wrapping but no guaranteed render shape, no frame, and no reason you can switch on. gorefuse gives you one carrier, one constructor, and one render contract.

## Install

```sh
go get github.com/adambarber/gorefuse
```

Go 1.27 or later; no dependencies outside the standard library.

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

- Use `Because[R]` when a refusal needs a domain-owned identity that a caller can switch on.
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

- **One verdict at line start.** Invalid UTF-8 and ESC become U+FFFD; `\r`, `\n`, `\v`, `\f`, NUL, U+0085, U+2028 and U+2029 become spaces. A hostile slot cannot forge a second line-start `Error:`.
- **Grep-able frame.** `At` accepts every file name; a name that would forge a `file:line` pairing renders as `<invalid-frame>` with no line. A non-positive line never pairs. A `File` ending in `:` followed by non-ASCII digits is prose by contract; only ASCII `:digits` is the pairing grammar.
- **Marked slots are redacted.** A marked `Subject` or `Why` renders as `[redacted]`; a marked `Cause` renders nothing. Marked cause slots still unwrap. Unknown sensitive bits fail closed on every construction path.
- **Reason identity, not reason words.** `Is` matches only an error-kinded reason; a non-error reason domain is matched with `errors.As` on the carrier. A nil target never matches and another reason type never bridges domains.
- **The zero value is an empty error.** A zero `Because[R]` renders `""` and is still a non-nil `error`.
- **Byte-stable.** Exact-byte tests and a fuzz target pin validity, idempotence, a single line-start verdict, no trailing connector, and slot redaction.

## Error and failure contract

- `From[R]` is nil-safe: untyped nil or a typed nil pointer/map/slice/func/chan yields the reason alone.
- `Describe` must not panic; `From` does not recover.
- A `Cause` whose `Error()` panics is the caller's defect; gorefuse does not recover.
- `At` panics only on caller mistakes: negative line, line without a file, or unknown sensitive bits. Hostile file names are data.

## Comparison

- **vs `fmt.Errorf`:** typed reason identity and a fixed, parseable render; you give up free-form lowercase error fragments.
- **vs `cockroachdb/errors`:** no stack traces, network encoding, or Sentry reports; standard-library-only and byte-pinned instead.
- **vs `rotisserie/eris`:** no JSON formatter yet; prose rendering is the only output today.
- **vs `joomcode/errorx`:** the reason is a type parameter rather than a trait system; smaller surface, fewer capabilities.

## Caveats

- The `Error:` prefix means a wrapped chain reads `starting server: Error: server.go:42: ...`. Compose for humans and use `errors.Is` for control flow.
- Redaction is a rendering decision, not a security boundary: only marked slots are hidden, and the marker is guessable.
- Structure is prose-only today; machine-readable output is planned and will derive from the same described parts rather than re-render the string.

## Examples

[`examples/`](examples/README.md) has standalone recipes you can copy into a project: basic refusal, reason matching, custom describers, redaction, cause wrapping, and hostile frames. Each was run with `go run` and its output recorded; the same scenarios are published as executable `Example...` functions on [pkg.go.dev](https://pkg.go.dev/github.com/adambarber/gorefuse).

## Status

Pre-1.0. Unreleased on `main`; the next tag will be `v0.2.0` and the last tag is `v0.1.1`. The pending release changes render bytes: every line break is flattened, `At` is total on file names, and `Is` narrows to error-kinded reasons. Exact-byte pins were updated in the same commits.

[![PkgGoDev](https://pkg.go.dev/badge/github.com/adambarber/gorefuse)](https://pkg.go.dev/github.com/adambarber/gorefuse)

## License

[MIT](LICENSE)
