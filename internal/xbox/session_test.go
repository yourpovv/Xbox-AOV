package xbox

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"aov/internal/config"
)

func TestResolveXUIDNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"code":28,"StatusCode":404}`))
	}))
	defer server.Close()

	settings := config.Config{
		ProfileURLFormat: server.URL + "/users/gt(%s)/profile/settings",
		ProfileContract:  "2",
	}
	session := Session{settings: settings, client: server.Client()}

	_, err := session.ResolveXUID("364b")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if want := `resolve gamertag "364b" to XUID: no Xbox profile found`; err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
	if strings.Contains(err.Error(), `"code":28`) {
		t.Errorf("transport detail leaked into user-facing error: %q", err.Error())
	}
}
