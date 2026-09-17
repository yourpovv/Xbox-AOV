package lookup

import (
	"strings"
	"testing"

	"aov/internal/xbox"
)

func TestLinesAlignedFormat(t *testing.T) {
	friends := []xbox.Person{
		{XUID: "2535430838089144", Gamertag: "Magic Spork313", GamerScore: "52005"},
		{XUID: "2535419146258107", Gamertag: "eb4q"},
		{XUID: "not-a-number", DisplayName: "NoGamertag"},
	}

	lines := Lines(friends)

	want := []string{
		"Friends: 3",
		"",
		"Magic Spork313  (xuid: 2535430838089144) (gamerscore: 52005)",
		"eb4q            (xuid: 2535419146258107) (gamerscore: -)",
		"NoGamertag      (xuid: not-a-number) (gamerscore: -)",
	}

	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d:\n%s", len(lines), len(want), strings.Join(lines, "\n"))
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Errorf("line %d:\n got %q\nwant %q", i, lines[i], want[i])
		}
	}
}

func TestLinesEmpty(t *testing.T) {
	lines := Lines(nil)

	if len(lines) != 1 || lines[0] != "Friends: 0" {
		t.Fatalf("got %q, want [\"Friends: 0\"]", lines)
	}
}
