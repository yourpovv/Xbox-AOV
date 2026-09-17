package xbox

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"aov/internal/config"
)

// zero value is unusable
type Session struct {
	settings      config.Config
	client        *http.Client
	authorization string
}

func Authenticate(settings config.Config, client *http.Client) (Session, error) {
	session := Session{settings: settings, client: client}
	var tokens authTokens
	headers := map[string]string{"Authorization": settings.APIToken}
	if err := session.fetchJSON(settings.AuthURL, headers, &tokens); err != nil {
		return Session{}, fmt.Errorf("authenticate with x-bot.live: %w", err)
	}
	if tokens.UserHash == "" || tokens.XSTSToken == "" {
		return Session{}, fmt.Errorf("authenticate with x-bot.live: response missed userHash or XSTSToken")
	}
	session.authorization = "XBL3.0 x=" + tokens.UserHash + ";" + tokens.XSTSToken
	return session, nil
}

func (s Session) ResolveXUID(gamertag string) (string, error) {
	var profile profileResponse
	endpoint := fmt.Sprintf(s.settings.ProfileURLFormat, url.PathEscape(gamertag))
	if err := s.fetchJSON(endpoint, s.authHeaders(s.settings.ProfileContract), &profile); err != nil {
		var missing *notFoundError
		if errors.As(err, &missing) {
			return "", fmt.Errorf("resolve gamertag %q to XUID: no Xbox profile found", gamertag)
		}
		return "", fmt.Errorf("resolve gamertag %q to XUID: %w", gamertag, err)
	}
	if len(profile.Users) == 0 || profile.Users[0].ID == "" {
		return "", fmt.Errorf("resolve gamertag %q to XUID: no profile user returned", gamertag)
	}
	return profile.Users[0].ID, nil
}

func (s Session) FetchFriends(xuid string) ([]Person, error) {
	var list friendList
	endpoint := fmt.Sprintf(s.settings.SocialURLFormat, url.PathEscape(xuid))
	if err := s.fetchJSON(endpoint, s.authHeaders(s.settings.SocialContract), &list); err != nil {
		return nil, fmt.Errorf("fetch friends of XUID %s: %w", xuid, err)
	}
	return list.People, nil
}

func (s Session) authHeaders(contractVersion string) map[string]string {
	return map[string]string{
		"X-XBL-Contract-Version": contractVersion,
		"Authorization":          s.authorization,
		"Accept-Language":        "en-US",
	}
}
