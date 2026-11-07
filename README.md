# gorefuse

A typed refusal vocabulary for Go: one way to say why a thing was refused.

```go
err := gorefuse.At("server.go", 42, gorefuse.Subject("listen :8080"),
	gorefuse.Why("bind the port before serving"))

// renders: "server.go:42: refused: listen :8080. bind the port before serving"
```

- **`Because[R]`** — a refusal paired with a typed reason, carried as an
  `error`: unwraps to its cause, composes through `errors.Is`/`errors.As`,
  and never matches across reason domains — the type parameter is a
  compile-time wall, proven both by `errors.As` behavior and by a
  tag-gated cross-domain assignment test that fails the build.
- **`Refusal`** — named kinds `File`, `Line`, `Subject`, `Why`; a refusal
  carries where and what, not prose fragments.
- **`At()`** — the single constructor: refuses (panics) a File containing
  a newline or a trailing `:digits` pattern, so the `file:line` frame is
  unambiguous — a log parser's `^(\S+):\d+:` match locks onto the true file.
- **`Describable` / `From[R]`** — refusals described by values: a nil
  Describable means no description, for every nil-able kind.
- **A render contract, guaranteed** — output is well-formed UTF-8 with no
  raw carriage return and no ESC byte (each becomes U+FFFD or a space);
  embedded newlines pass through so the frame stays grep-able; bidi
  overrides pass through deliberately; whitespace-only slots contribute
  nothing; the sentence tail honors ASCII and ideographic terminals.
- **stdlib-only** — no dependencies, no reflection beyond the error
  interface.

## The contract

Render order is frame prefix → reason → subject → cause tail → fix
sentence. Every contribution slot obeys one rule: empty-after-trim means
the slot contributes nothing, no connector. The zero `Because[R]` for a
string-kinded reason renders `""` — and is still a valid `error` with
identity, because the reason type is the identity, not the words.

## Status

API is stable modulo the pre-1.0 conventions of an early library.

## Roadmap

- **Structured rendering of refusals**: a refusal should be serializable
  into a machine-readable format. The prose render stays the human path;
  serialization must not fork the render contract (stable order,
  whitespace rules), so the structured form derives from the same
  described parts, not re-rendered.
