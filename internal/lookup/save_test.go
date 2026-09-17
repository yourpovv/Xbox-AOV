package lookup

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"aov/internal/xbox"
)

func TestSavePfp(t *testing.T) {
	picture := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(picture)
	}))
	defer server.Close()

	dir := t.TempDir()
	saver := NewSaver(dir, server.Client())
	friend := xbox.Person{Gamertag: "Some One", DisplayPicRaw: server.URL + "/pic.png"}

	if err := saver.savePfp("e64b", friend); err != nil {
		t.Fatalf("savePfp returned an error: %v", err)
	}

	saved, err := os.ReadFile(filepath.Join(dir, "e64b", "pfps", "Some One.png"))
	if err != nil {
		t.Fatalf("saved pfp not found: %v", err)
	}
	if string(saved) != string(picture) {
		t.Errorf("saved pfp content mismatch: got %v, want %v", saved, picture)
	}
}

func TestSavePfpServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	saver := NewSaver(t.TempDir(), server.Client())
	friend := xbox.Person{Gamertag: "Broken", DisplayPicRaw: server.URL + "/pic.png"}

	if err := saver.savePfp("e64b", friend); err == nil {
		t.Fatal("expected an error for a non-200 response, got nil")
	}
}
