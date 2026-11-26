package gorefuse

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

var errFuzzCause = errors.New("cause")

// FuzzRedactionAcrossEncoders pins that a marked slot never survives any
// output surface, not just the prose renderer.
func FuzzRedactionAcrossEncoders(f *testing.F) {
	for _, seed := range []string{"", "hunter2", Redacted, "a\nb", "token\\key"} {
		f.Add([]byte(seed))
	}

	f.Fuzz(func(t *testing.T, value []byte) {
		planted := secretEncoderProbe(fmt.Sprintf("SECRET-%x", value))

		b := Because[string]{
			Reason: "denied",
			Cause:  fmt.Errorf("%w %s", errFuzzCause, planted),
			What: At(File("f.go"), 1,
				Subject(planted), Why(planted),
				SubjectSensitive|WhySensitive|CauseSensitive),
		}

		encoded, err := json.Marshal(b)
		if err != nil {
			t.Fatalf("json.Marshal: %v", err)
		}

		surfaces := map[encoderSurface]string{
			"Error":       b.Error(),
			"Verdict":     b.Verdict(),
			"GoString":    fmt.Sprintf("%+v", b),
			"JSON":        string(encoded),
			"RefusalJSON": string(jsonMustMarshal(t, b.What)),
		}

		for name, out := range surfaces {
			if strings.Contains(out, string(planted)) {
				t.Fatalf("%s echoed a marked value: %s", name, out)
			}
		}
	})
}

type secretEncoderProbe string

type encoderSurface string

type encoderOutput string

func jsonMustMarshal(t *testing.T, v Refusal) encoderOutput {
	t.Helper()

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal Refusal: %v", err)
	}

	return encoderOutput(b)
}
