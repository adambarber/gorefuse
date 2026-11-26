# Recipes

Each directory is a standalone `package main`. Copy `main.go` into your
project, keep or replace the `reasonError` type, and edit the words. Every
recipe is compiled by `go test ./...` and was run with `go run ./examples/<name>`.

| Recipe | Shows | Output |
| --- | --- | --- |
| [`basic`](basic/main.go) | `Because[R]` + `At` | `Error: server.go:42: denied: listen :8080. bind the port before serving` |
| [`reason_matching`](reason_matching/main.go) | `errors.Is`, `errors.As`, wrapping | `true` then `server.go 42` |
| [`custom_describer`](custom_describer/main.go) | `Describable` + `From` | `Error: config.yaml:7: invalid: port must be positive. set a value from 1 to 65535` |
| [`redaction`](redaction/main.go) | `SubjectSensitive`, `WhySensitive` | `Error: server.go:42: denied: [redacted]. [redacted]` |
| [`cause`](cause/main.go) | `Cause` tail and `errors.Is` through `Unwrap` | `Error: server.go:42: denied: bind: address in use` then `true` |
| [`hostile_frame`](hostile_frame/main.go) | forged file name and embedded line break | `Error: <invalid-frame>: denied: a Error: forged` |

Run one locally:

```sh
go run ./examples/basic
```

The same scenarios are published as executable `Example...` functions on
[pkg.go.dev](https://pkg.go.dev/github.com/adambarber/gorefuse).
