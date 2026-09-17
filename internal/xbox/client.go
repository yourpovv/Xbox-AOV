package xbox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (s Session) fetchJSON(requestURL string, headers map[string]string, target any) error {
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return fmt.Errorf("build request for %s: %w", requestURL, err)
	}

	applyBrowserHeaders(request, s.settings.UserAgent)
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := s.client.Do(request)
	if err != nil {
		return fmt.Errorf("get %s: %w", requestURL, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("get %s: %w", requestURL, err)
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return &notFoundError{requestURL: requestURL}
		}
		if response.StatusCode == http.StatusForbidden && strings.Contains(string(body), "Attention Required") {
			return fmt.Errorf("get %s: blocked by Cloudflare bot protection", requestURL)
		}
		return fmt.Errorf("get %s: unexpected %s: %s", requestURL, response.Status, strings.TrimSpace(string(body)))
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("decode response from %s: %w", requestURL, err)
	}
	return nil
}

// notFoundError marks an HTTP 404 so callers can translate a missing entity
// into domain language so we don't leak transport details
type notFoundError struct {
	requestURL string
}

func (e *notFoundError) Error() string {
	return fmt.Sprintf("get %s: unexpected 404 Not Found", e.requestURL)
}

// x-bot.live rejects the default Go client at the Cloudflare edge
func applyBrowserHeaders(request *http.Request, userAgent string) {
	defaults := [][2]string{
		{"User-Agent", userAgent},
		{"Accept", "application/json, text/plain, */*"},
		{"Accept-Language", "en-US,en;q=0.9"},
	}
	for _, header := range defaults {
		key, value := header[0], header[1]
		if request.Header.Get(key) == "" {
			request.Header.Set(key, value)
		}
	}
}
