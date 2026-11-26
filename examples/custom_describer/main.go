// Command custom_describer lets a domain type say what it refused while the
// caller supplies the reason identity.
package main

import (
	"fmt"

	"github.com/adambarber/gorefuse"
)

type reasonError string

const errInvalid reasonError = "invalid"

const errLine gorefuse.Line = 7

func (r reasonError) Error() string { return string(r) }

type domainRefusal struct {
	what gorefuse.Refusal
}

func (d domainRefusal) Describe() (gorefuse.Refusal, bool) {
	return d.what, true
}

func main() {
	b := gorefuse.From(errInvalid, domainRefusal{
		what: gorefuse.At("config.yaml", errLine,
			gorefuse.Subject("port must be positive"),
			gorefuse.Why("set a value from 1 to 65535")),
	})

	fmt.Println(b.Verdict())
}
