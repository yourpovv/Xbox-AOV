package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	APIToken              string `json:"apiToken"`
	RequestTimeoutSeconds int    `json:"requestTimeoutSeconds"`
	AuthURL               string `json:"authURL"`
	ProfileURLFormat      string `json:"profileURLFormat"`
	ProfileContract       string `json:"profileContract"`
	SocialURLFormat       string `json:"socialURLFormat"`
	SocialContract        string `json:"socialContract"`
	UserAgent             string `json:"userAgent"`
	LookupsDir            string `json:"lookupsDir"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("load config %s: %w", path, err)
	}
	var settings Config
	if err := json.Unmarshal(data, &settings); err != nil {
		return Config{}, fmt.Errorf("decode config %s: %w", path, err)
	}
	if missing := missingSettings(settings); len(missing) > 0 {
		return Config{}, fmt.Errorf("decode config %s: missing %s", path, strings.Join(missing, ", "))
	}
	return settings, nil
}

func missingSettings(settings Config) []string {
	var missing []string
	if settings.APIToken == "" {
		missing = append(missing, "apiToken (sign in at x-bot.live and click Reveal API Key)")
	}
	if settings.RequestTimeoutSeconds <= 0 {
		missing = append(missing, "requestTimeoutSeconds")
	}
	if settings.AuthURL == "" {
		missing = append(missing, "authURL")
	}
	if settings.ProfileURLFormat == "" {
		missing = append(missing, "profileURLFormat")
	}
	if settings.ProfileContract == "" {
		missing = append(missing, "profileContract")
	}
	if settings.SocialURLFormat == "" {
		missing = append(missing, "socialURLFormat")
	}
	if settings.SocialContract == "" {
		missing = append(missing, "socialContract")
	}
	if settings.UserAgent == "" {
		missing = append(missing, "userAgent")
	}
	if settings.LookupsDir == "" {
		missing = append(missing, "lookupsDir")
	}
	return missing
}
