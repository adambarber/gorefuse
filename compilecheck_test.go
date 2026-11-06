//go:build compilecheck

package gorefuse

// This file must fail to compile: Because of two different reason types are
// distinct types, and assigning one to the other is the leak the generic
// exists to prevent. `go test -tags=compilecheck` expects this build to fail.
func crossDomainAssign() {
	a := Because[testReason]{Reason: testRefused}
	var _ Because[otherReason] = a
}
