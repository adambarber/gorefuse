package gorefuse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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

		b := Because[encoderReason]{
			Reason: encoderReason("denied"),
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
			"GoSharp":     fmt.Sprintf("%#v", b),
			"JSON":        string(encoded),
			"RefusalFmt":  string(refusalFmt(b.What)),
			"RefusalJSON": string(jsonMustMarshal(t, b.What)),
			"RefusalSlog": string(slogMustText(t, b.What)),
			"SlogJSON":    string(slogMustJSON(t, b)),
		}

		for name, out := range surfaces {
			if strings.Contains(out, string(planted)) {
				t.Fatalf("%s echoed a marked value: %s", name, out)
			}
		}
	})
}

type secretEncoderProbe string

type encoderReason string

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

func slogMustJSON(t *testing.T, v Because[encoderReason]) encoderOutput {
	t.Helper()

	var out bytes.Buffer

	logger := slog.New(slog.NewJSONHandler(&out, nil))
	logger.Info("refusal", "error", v)

	return encoderOutput(out.String())
}

func refusalFmt(r Refusal) encoderOutput {
	return encoderOutput(fmt.Sprintf("%v %+v %#v", r, r, r))
}

func slogMustText(t *testing.T, r Refusal) encoderOutput {
	t.Helper()

	var out bytes.Buffer

	logger := slog.New(slog.NewTextHandler(&out, nil))
	logger.Info("refusal", "what", r)

	return encoderOutput(out.String())
}

func TestStructuredFrameGuard(t *testing.T) {
	b := Refuse(encoderReason("denied"),
		At(File("evil\n10.0.0.1:9"), 7, Subject("s"), Why("")))

	out := string(slogMustJSON(t, b))
	if !strings.Contains(out, "<invalid-frame>") {
		t.Fatalf("structured output lost the frame guard: %s", out)
	}

	if strings.Contains(out, "10.0.0.1") {
		t.Fatalf("structured output carried a hostile file: %s", out)
	}
}
