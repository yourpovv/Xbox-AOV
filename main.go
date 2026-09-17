package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"aov/internal/banner"
	"aov/internal/config"
	"aov/internal/lookup"
	"aov/internal/xbox"
)

func main() {
	if err := runLookup(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

var stdin = bufio.NewReader(os.Stdin)

func runLookup() error {
	settings, err := config.Load("config.json")
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: time.Duration(settings.RequestTimeoutSeconds) * time.Second}
	session, err := xbox.Authenticate(settings, client)
	if err != nil {
		return err
	}
	saver := lookup.NewSaver(settings.LookupsDir, client)
	for {
		banner.Show()
		gamertag, err := readGamertag()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if gamertag == "" {
			continue
		}
		if err := lookupPlayer(session, saver, gamertag); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}
	}
}

func lookupPlayer(session xbox.Session, saver lookup.Saver, gamertag string) error {
	xuid, err := session.ResolveXUID(gamertag)
	if err != nil {
		return err
	}
	friends, err := session.FetchFriends(xuid)
	if err != nil {
		return err
	}
	lines := lookup.Lines(friends)
	lookup.Print(lines)
	export, err := confirmExport()
	if err != nil {
		return err
	}
	if !export {
		return nil
	}
	path, err := saver.Save(gamertag, lines)
	if err != nil {
		return err
	}
	fmt.Println("Saved to " + path)
	saver.SavePfps(gamertag, friends)
	return nil
}

func readGamertag() (string, error) {
	gamertag, err := readLine("Enter player gamertag: ")
	if err != nil {
		return "", fmt.Errorf("read gamertag: %w", err)
	}
	return gamertag, nil
}

func confirmExport() (bool, error) {
	answer, err := readLine("\r\nExport Data? Y/N: ")
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false, nil
		}
		return false, fmt.Errorf("read export answer: %w", err)
	}
	answer = strings.ToLower(answer)
	return answer == "y" || answer == "yes", nil
}

func readLine(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	line, err := stdin.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
