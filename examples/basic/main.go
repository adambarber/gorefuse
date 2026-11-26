// Command basic builds one refusal with a reason, a frame, a subject, and a fix.
// Copy this file into your project and replace the reason type and words.
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
			gorefuse.Subject("listen :8080"),
			gorefuse.Why("bind the port before serving")),
		Cause: nil,
	}

	fmt.Println(err)
}
