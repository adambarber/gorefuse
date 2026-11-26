// Command cause shows how an underlying failure renders as a cause tail and
// stays reachable through Unwrap.
package main

import (
	"errors"
	"fmt"

	"github.com/adambarber/gorefuse"
)

type reasonError string

const errDenied reasonError = "denied"

var errBind = errors.New("bind: address in use")

const errLine gorefuse.Line = 42

func (r reasonError) Error() string { return string(r) }

func main() {
	b := gorefuse.Because[reasonError]{
		Reason: errDenied,
		Cause:  errBind,
		What:   gorefuse.At("server.go", errLine, "", ""),
	}

	fmt.Println(b.Verdict())
	fmt.Println(errors.Is(b, errBind))
}
