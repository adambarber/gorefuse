# Changelog

## Unreleased

No changes since v0.2.0.

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
