package lookup

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"aov/internal/xbox"
)

const (
	dirPerm  = 0o755
	filePerm = 0o644
)

type Saver struct {
	dir    string
	client *http.Client
}

func NewSaver(dir string, client *http.Client) Saver {
	return Saver{dir: dir, client: client}
}

func (s Saver) lookupDir(gamertag string) string {
	return filepath.Join(s.dir, sanitizeFilename(gamertag))
}

func (s Saver) Save(gamertag string, lines []string) (string, error) {
	folder := s.lookupDir(gamertag)
	if err := os.MkdirAll(folder, dirPerm); err != nil {
		return "", fmt.Errorf("save lookup for %q: %w", gamertag, err)
	}
	path := filepath.Join(folder, sanitizeFilename(gamertag)+".txt")
	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(path, []byte(content), filePerm); err != nil {
		return "", fmt.Errorf("save lookup for %q: %w", gamertag, err)
	}
	return path, nil
}

func (s Saver) SavePfps(lookup string, friends []xbox.Person) {
	for _, friend := range friends {
		if friend.DisplayPicRaw == "" {
			continue
		}
		if err := s.savePfp(lookup, friend); err != nil {
			fmt.Fprintln(os.Stderr, "pfp unavailable for "+friendName(friend)+": "+err.Error())
		}
	}
}

func (s Saver) savePfp(lookup string, friend xbox.Person) error {
	folder := filepath.Join(s.lookupDir(lookup), "pfps")
	if err := os.MkdirAll(folder, dirPerm); err != nil {
		return err
	}
	response, err := s.client.Get(friend.DisplayPicRaw)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected %s", response.Status)
	}
	picture, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(folder, sanitizeFilename(friendName(friend))+".png"), picture, filePerm)
}

func sanitizeFilename(name string) string {
	cleaned := strings.Map(func(r rune) rune {
		if strings.ContainsRune("<>:\"/\\|?*", r) || unicode.IsControl(r) {
			return '_'
		}
		return r
	}, name)
	cleaned = strings.TrimRight(cleaned, " .")
	if cleaned == "" {
		return "lookup"
	}
	return cleaned
}
