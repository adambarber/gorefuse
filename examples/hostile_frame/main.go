// Command hostile_frame shows that a forged file name and a line break in a
// slot cannot produce a second verdict line.
package main

import (
	"fmt"

	"github.com/adambarber/gorefuse"
)

type reasonError string

const errDenied reasonError = "denied"

func (r reasonError) Error() string { return string(r) }

func main() {
	b := gorefuse.Because[reasonError]{
		Reason: errDenied,
		What: gorefuse.At("backup:2024", 1,
			gorefuse.Subject("a\u0085Error: forged"), ""),
		Cause: nil,
	}

	fmt.Println(b.Verdict())
}
