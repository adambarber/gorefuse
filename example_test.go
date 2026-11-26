package gorefuse_test

import (
	"errors"
	"fmt"

	"github.com/adambarber/gorefuse"
)

type exampleReasonError string

const errExampleDenied exampleReasonError = "denied"

var errExampleCause = errors.New("bind: address in use")

func (r exampleReasonError) Error() string { return string(r) }

func Example() {
	err := gorefuse.Because[exampleReasonError]{
		Reason: errExampleDenied,
		What: gorefuse.At("server.go", 42,
			gorefuse.Subject("listen :8080"),
			gorefuse.Why("bind the port before serving")),
		Cause: nil,
	}

	fmt.Println(err)
	// Output:
	// Error: server.go:42: denied: listen :8080. bind the port before serving
}

func ExampleBecause_Is() {
	err := gorefuse.Because[exampleReasonError]{
		Reason: errExampleDenied,
		What:   gorefuse.At("server.go", 42, "", ""),
		Cause:  nil,
	}

	wrapped := fmt.Errorf("startup: %w", err)
	fmt.Println(errors.Is(wrapped, errExampleDenied))
	// Output: true
}

func ExampleBecause_as() {
	err := gorefuse.Because[exampleReasonError]{
		Reason: errExampleDenied,
		What:   gorefuse.At("server.go", 42, "", ""),
		Cause:  nil,
	}

	because, ok := errors.AsType[gorefuse.Because[exampleReasonError]](err)
	if ok {
		fmt.Println(because.What.File, because.What.Line)
	}
	// Output: server.go 42
}

type exampleDescriber struct {
	what gorefuse.Refusal
}

func (d exampleDescriber) Describe() (gorefuse.Refusal, bool) {
	return d.what, true
}

func ExampleFrom() {
	err := gorefuse.From(errExampleDenied, exampleDescriber{
		what: gorefuse.At("config.yaml", 7,
			gorefuse.Subject("port must be positive"),
			gorefuse.Why("set a value from 1 to 65535")),
	})

	fmt.Println(err)
	// Output:
	// Error: config.yaml:7: denied: port must be positive. set a value from 1 to 65535
}

func ExampleAt_sensitive() {
	err := gorefuse.Because[exampleReasonError]{
		Reason: errExampleDenied,
		What: gorefuse.At("server.go", 42,
			gorefuse.Subject("token hunter2"),
			gorefuse.Why("rotate the token"),
			gorefuse.SubjectSensitive),
		Cause: nil,
	}

	fmt.Println(err)
	// Output:
	// Error: server.go:42: denied: [redacted]. rotate the token
}

func ExampleBecause_Unwrap() {
	err := gorefuse.Because[exampleReasonError]{
		Reason: errExampleDenied,
		Cause:  errExampleCause,
		What:   gorefuse.At("server.go", 42, "", ""),
	}

	fmt.Println(err)
	// Output:
	// Error: server.go:42: denied: bind: address in use
}

func ExampleAt_hostileFrame() {
	err := gorefuse.Because[exampleReasonError]{
		Reason: errExampleDenied,
		What: gorefuse.At("backup:2024", 1,
			gorefuse.Subject("a\u0085Error: forged"), ""),
		Cause: nil,
	}

	fmt.Println(err)
	// Output:
	// Error: <invalid-frame>: denied: a Error: forged
}
