// Command redaction marks the slots that carry secret-classified text.
package main

import (
	"fmt"

	"github.com/adambarber/gorefuse"
)

type reasonError string

const errDenied reasonError = "denied"

const errLine gorefuse.Line = 42

func (r reasonError) Error() string { return string(r) }

func main() {
	err := gorefuse.Because[reasonError]{
		Reason: errDenied,
		What: gorefuse.At("server.go", errLine,
			gorefuse.Subject("token hunter2"),
			gorefuse.Why("rotate the token"),
			gorefuse.SubjectSensitive,
			gorefuse.WhySensitive),
		Cause: nil,
	}

	fmt.Println(err)
}
