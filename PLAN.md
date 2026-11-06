# gorefuse

The typed refusal vocabulary library: one way to say why a thing was
refused. Sister library to goconf (config validation); follows the same
conventions: birdcage policy gate, strictest merged lint config, TDD.

## Status

Batch 1 landed: `Because[R]` with rendering (frame prefix, reason, subject,
cause tail, fix sentence), `Refusal` and the named kinds (`File`, `Line`,
`Subject`, `Why`), `Unwrap` to `Cause`. Batch 2 landed: `At()` — the single Refusal constructor — `Describable`,
and `From[R](reason, Describable)`. Remaining: compile-time guarantee
polish (Batch 3).
