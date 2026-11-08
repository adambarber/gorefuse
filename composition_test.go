package gorefuse

import (
	"testing"
)

// confFileSpell is the consuming-library File spelling the composition model
// uses: the path and the line live in the frame, the field name in the reason.
type confFileSpell string

// TestCompositionPinsOwnerModel pins the composition contract: the location
// phrase rides the reason, the fix rides the Why, and the fixed render order
// delivers the model byte for byte; the second test pins the frame:line form.
func TestCompositionPinsOwnerModel(t *testing.T) {
	b := Because[confFileSpell]{
		Reason: `unknown field "extra" in /path/app.json, line 1`,
		Cause:  nil,
		What:   At("", 0, "", Why("remove that line from the file, or the app must declare the field in its schema.")),
	}

	assertRender(t, b, renderText(`Error: unknown field "extra" in /path/app.json, line 1. remove that line from the file, or the app must declare the field in its schema.`))
}

// TestCompositionFrameCarriesPathLine pins the frame half of the composition
// contract: the frame prefix is `path:line:`, so a consuming library's reason
// reads as prose after it — the grep-able pairing stays machine-parseable.
func TestCompositionFrameCarriesPathLine(t *testing.T) {
	b := From(confFileSpell(`unknown field "extra"`), testDescriber{
		what:   At("/path/app.json", 1, Subject(""), Why("")),
		misses: false,
	})

	got := b.Error()
	want := `Error: /path/app.json:1: unknown field "extra"`

	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}
