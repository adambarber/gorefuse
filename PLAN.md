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
