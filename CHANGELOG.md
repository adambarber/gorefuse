# Changelog

## v0.2.1

- `fmt` verbs (including `%#v`), `MarshalJSON`, and `slog.LogValuer` carry the same slot decisions as the prose renderer, on both the carrier and a bare `Refusal`: marked Subject and Why become `[redacted]`, and a marked Cause is omitted.
- The frame guard now applies to structured output too; a hostile file renders `<invalid-frame>` with line 0 in JSON and slog.
- Added `Refuse`, `Wrap`, and `Join` constructors; `Join` composes causes with `errors.Join` and keeps standard `errors.Is`/`errors.As` traversal.
- Lowered the Go floor from 1.27 to 1.21; the shipped package needs only generics and stdlib basics.
- The render contract is documented as a pinned prose shape, not immutable bytes before v1.0.

## v0.2.0

Breaking changes from v0.1.1:

- The verdict label moved from `Error()` to `Verdict()`. Wrapping now reads `starting server: server.go:42: ...` instead of embedding a second `Error:` mid-string.
- `Error()` flattens `\r`, `\n`, `\v`, `\f`, NUL, U+0085, U+2028 and U+2029 to spaces. Earlier versions emitted embedded line breaks.
- NUL-only slots are empty after scrubbing; they no longer emit dangling `": "` or `". "` connectors.
- Unknown `Sensitive` bits fail closed on hand-built and `From` carriers as well as `At`: the frame becomes `<invalid-frame>` and populated Subject/Why slots render as `[redacted]`.
- A `File` containing carriage return renders `<invalid-frame>`, matching the other line breaks.
- `Is` matches only error-kinded reasons; equality against non-error reason values no longer matches. Use `errors.As` on the carrier for those domains.
- `At` no longer panics on hostile file names. Ambiguous names render `<invalid-frame>` instead.

Other changes: the lint gate was repaired, the README was rewritten, and the MIT license was standardized. A typed-nil error reason renders as empty; `From` and reason rendering share one nil-value check. The README states that `File` is never redacted and that redaction happens after scrubbing.
