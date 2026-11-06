package gorefuse

import (
	"strings"
	"testing"
)

func TestAtRefusesNegativeLine(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("At accepted a negative Line; want a refusal")
		}

		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic is %T, want string", r)
		}

		for _, want := range []string{"-2", string(testFile)} {
			if !strings.Contains(msg, want) {
				t.Fatalf("panic %q does not echo %q", msg, want)
			}
		}
	}()

	At(testFile, -2, testSubject, "")
}
