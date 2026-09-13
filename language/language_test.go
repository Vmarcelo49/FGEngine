package language

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(thisFile), "..")
}

func TestLoadLang(t *testing.T) {
	t.Chdir(repoRoot(t))
	en, err := LoadLang(English)
	if err != nil {
		t.Fatalf("load EN: %v", err)
	}
	if en.Lang != English {
		t.Fatalf("lang = %q, want EN", en.Lang)
	}
	if en.GameText["match"] != "Match" {
		t.Fatalf("EN match text = %q", en.GameText["match"])
	}

	br, err := LoadLang(Portuguese)
	if err != nil {
		t.Fatalf("load BR: %v", err)
	}
	if br.GameText["match"] != "Jogar" {
		t.Fatalf("BR match text = %q", br.GameText["match"])
	}
}

func TestImportTOMLRejectsUnknownFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lang.toml")
	content := "lang = \"EN\"\ntypo_field = 1\n\n[game_text]\nmatch = \"Match\"\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := ImportTOML(path); err == nil {
		t.Fatal("expected strict-decode error for unknown field, got nil")
	}
}
