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
	b := gorefuse.Because[exampleReasonError]{
		Reason: errExampleDenied,
		What: gorefuse.At("server.go", 42,
			gorefuse.Subject("listen :8080"),
			gorefuse.Why("bind the port before serving")),
		Cause: nil,
	}

	fmt.Println(b.Verdict())
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

	var because gorefuse.Because[exampleReasonError]
	if errors.As(err, &because) {
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
	b := gorefuse.From(errExampleDenied, exampleDescriber{
		what: gorefuse.At("config.yaml", 7,
			gorefuse.Subject("port must be positive"),
			gorefuse.Why("set a value from 1 to 65535")),
	})

	fmt.Println(b.Verdict())
	// Output:
	// Error: config.yaml:7: denied: port must be positive. set a value from 1 to 65535
}

func ExampleAt_sensitive() {
	b := gorefuse.Because[exampleReasonError]{
		Reason: errExampleDenied,
		What: gorefuse.At("server.go", 42,
			gorefuse.Subject("token hunter2"),
			gorefuse.Why("rotate the token"),
			gorefuse.SubjectSensitive),
		Cause: nil,
	}

	fmt.Println(b.Verdict())
	// Output:
	// Error: server.go:42: denied: [redacted]. rotate the token
}

func ExampleBecause_Unwrap() {
	b := gorefuse.Because[exampleReasonError]{
		Reason: errExampleDenied,
		Cause:  errExampleCause,
		What:   gorefuse.At("server.go", 42, "", ""),
	}

	fmt.Println(b.Verdict())
	// Output:
	// Error: server.go:42: denied: bind: address in use
}

func ExampleAt_unmarkedRedactedToken() {
	b := gorefuse.Because[exampleReasonError]{
		Reason: errExampleDenied,
		What:   gorefuse.At("server.go", 42, "[redacted]", ""),
		Cause:  nil,
	}

	fmt.Println(b.Verdict())
	// Output:
	// Error: server.go:42: denied: [redacted]
}

func ExampleAt_hostileFrame() {
	b := gorefuse.Because[exampleReasonError]{
		Reason: errExampleDenied,
		What: gorefuse.At("backup:2024", 1,
			gorefuse.Subject("a\u0085Error: forged"), ""),
		Cause: nil,
	}

	fmt.Println(b.Verdict())
	// Output:
	// Error: <invalid-frame>: denied: a Error: forged
}
