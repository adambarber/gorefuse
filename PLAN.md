# gorefuse

The typed refusal vocabulary library: one way to say why a thing was
refused. Sister library to goconf (config validation); follows the same
conventions: birdcage policy gate, strictest merged lint config, TDD.

## Status

Batch 1 landed: `Because[R]` with rendering (frame prefix, reason, subject,
cause tail, fix sentence), `Refusal` and the named kinds (`File`, `Line`,
`Subject`, `Why`), `Unwrap` to `Cause`. Remaining: `Describable` + `At()`
constructor (Batch 2), compile-time guarantees polish (Batch 3).
