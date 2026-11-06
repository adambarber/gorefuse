# gorefuse

The typed refusal vocabulary library: one way to say why a thing was
refused. Sister library to goconf (config validation); follows the same
conventions: birdcage policy gate, strictest merged lint config, TDD.

## Status

Batch 1 landed: `Because[R]` with rendering (frame prefix, reason, subject,
cause tail, fix sentence), `Refusal` and the named kinds (`File`, `Line`,
`Subject`, `Why`), `Unwrap` to `Cause`. Batch 2 landed: `At()` — the single Refusal constructor — `Describable`,
and `From[R](reason, Describable)`. Batch 3 landed. Complete: `Because[R]` (error, unwraps to Cause, renders
frame prefix → reason → subject → cause tail → fix sentence), `Refusal`
with named kinds `File`/`Line`/`Subject`/`Why`, `At()` as the single
constructor, `Describable`, `From[R]`. Compile-time domain wall tested
two ways: `errors.As` never matches across reason domains, and a
tag-gated assignment (`compilecheck_test.go`) is proven to fail the
build by `TestCrossDomainAssignmentDoesNotCompile`.
