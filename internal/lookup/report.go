package lookup

import (
	"fmt"

	"aov/internal/xbox"
)

const nameGap = 2

type friendStats struct {
	name       string
	xuid       string
	gamerscore string
}

func Lines(friends []xbox.Person) []string {
	rows := make([]friendStats, 0, len(friends))
	for _, friend := range friends {
		rows = append(rows, statsOf(friend))
	}
	lines := make([]string, 0, len(rows)+2)
	lines = append(lines, fmt.Sprintf("Friends: %d", len(rows)))
	if len(rows) == 0 {
		return lines
	}
	lines = append(lines, "")
	width := nameColumnWidth(rows)
	for _, row := range rows {
		lines = append(lines, formatRow(row, width))
	}
	return lines
}

func Print(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}

func statsOf(friend xbox.Person) friendStats {
	gamerscore := friend.GamerScore
	if gamerscore == "" {
		gamerscore = "-"
	}
	return friendStats{
		name:       friendName(friend),
		xuid:       friend.XUID,
		gamerscore: gamerscore,
	}
}

func friendName(friend xbox.Person) string {
	if friend.Gamertag != "" {
		return friend.Gamertag
	}
	if friend.DisplayName != "" {
		return friend.DisplayName
	}
	return friend.XUID
}

func nameColumnWidth(rows []friendStats) int {
	width := 0
	for _, row := range rows {
		width = max(width, len(row.name))
	}
	return width + nameGap
}

func formatRow(row friendStats, width int) string {
	return fmt.Sprintf("%-*s(xuid: %s) (gamerscore: %s)", width, row.name, row.xuid, row.gamerscore)
}
