package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write test config: %v", err)
	}
	return path
}

func TestLoadValid(t *testing.T) {
	path := writeConfig(t, `{
		"apiToken": "test-key",
		"requestTimeoutSeconds": 15,
		"authURL": "https://example.com/auth",
		"profileURLFormat": "https://example.com/profile/%s",
		"profileContract": "2",
		"socialURLFormat": "https://example.com/social/%s",
		"socialContract": "5",
		"userAgent": "test-agent",
		"lookupsDir": "Lookups"
	}`)

	settings, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned an error: %v", err)
	}
	if settings.LookupsDir != "Lookups" {
		t.Errorf("LookupsDir = %q, want %q", settings.LookupsDir, "Lookups")
	}
}

func TestLoadMissingFile(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "absent.json")); err == nil {
		t.Fatal("expected an error for a missing file, got nil")
	}
}

func TestLoadBadJSON(t *testing.T) {
	path := writeConfig(t, `{"apiToken":`)

	if _, err := Load(path); err == nil {
		t.Fatal("expected an error for malformed JSON, got nil")
	}
}

func TestLoadMissingFields(t *testing.T) {
	path := writeConfig(t, `{}`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected an error for missing fields, got nil")
	}
	for _, field := range []string{"apiToken", "requestTimeoutSeconds", "authURL", "lookupsDir"} {
		if !strings.Contains(err.Error(), field) {
			t.Errorf("error %q does not name missing field %q", err.Error(), field)
		}
	}
}
