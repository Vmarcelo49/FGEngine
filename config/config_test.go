package config

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

// The default config must survive a marshal/unmarshal round-trip unchanged.
func TestDefaultConfigRoundTrip(t *testing.T) {
	want := loadDefaultConfig()

	data, err := toml.Marshal(want)
	if err != nil {
		t.Fatalf("marshal default config: %v", err)
	}
	if !strings.Contains(string(data), "window_width") {
		t.Fatalf("marshalled config lost key names:\n%s", data)
	}

	var got Config
	decoder := toml.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&got); err != nil {
		t.Fatalf("unmarshal default config: %v", err)
	}
	if got != want {
		t.Fatalf("round-trip mismatch:\nwant %+v\ngot  %+v", want, got)
	}
}

// Strict decoding (as used by LoadConfigFile) must reject unknown fields
// so typos surface instead of silently no-op'ing.
func TestStrictDecodeRejectsUnknownFields(t *testing.T) {
	data := []byte("window_width = 1600\nwindow_height = 900\ndeadzone = 0.3\nlanguage = \"EN\"\ntypo_field = 1\n")

	var config Config
	decoder := toml.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&config); err == nil {
		t.Fatal("expected strict decode error for unknown field, got nil")
	}
}
