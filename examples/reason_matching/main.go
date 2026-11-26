// Command reason_matching shows how to match a refusal through ordinary
// wrapping and how to extract the typed carrier for its frame.
package main

import (
	"errors"
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
		What:   gorefuse.At("server.go", errLine, "", ""),
		Cause:  nil,
	}
	wrapped := fmt.Errorf("startup: %w", err)

	fmt.Println(errors.Is(wrapped, errDenied))

	because, ok := errors.AsType[gorefuse.Because[reasonError]](wrapped)
	if ok {
		fmt.Println(because.What.File, because.What.Line)
	}
}
