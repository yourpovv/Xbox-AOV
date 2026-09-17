package banner

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

const (
	reset          = "\x1b[0m"
	edgeColor      = "\x1b[38;2;243;239;232m"
	titleColor     = "\x1b[38;2;246;240;228m"
	subtitleColor  = "\x1b[38;2;82;176;67m"
	fallbackWidth  = 80
	clearScreenSeq = "\x1b[2J\x1b[H"
)

const (
	bannerTitle    = "yourpov.dev"
	bannerSubtitle = "Xbox AOV"
)

// -11 wont fit uintptr which is why im using 10
const stdOutputHandle = ^uintptr(10)

const (
	utf8CodePage                  = 65001
	enableVirtualTerminalHandling = 0x0004
)

func Show() {
	enableWindowsConsole()
	fmt.Print(clearScreenSeq)
	contentWidth := max(len(bannerTitle), len(bannerSubtitle)) + 4
	margin := strings.Repeat(" ", max((consoleWidth()-contentWidth-2)/2, 0))
	bar := strings.Repeat("─", contentWidth)
	divider := strings.Repeat("─", contentWidth-2)

	fmt.Println()
	fmt.Println(margin + edgeColor + "┌" + bar + "┐" + reset)
	fmt.Println(textRow(margin, bannerLine{text: bannerTitle, color: titleColor}, contentWidth))
	fmt.Println(margin + edgeColor + "│ " + reset + divider + edgeColor + " │" + reset)
	fmt.Println(textRow(margin, bannerLine{text: bannerSubtitle, color: subtitleColor}, contentWidth))
	fmt.Println(margin + edgeColor + "└" + bar + "┘" + reset)
	fmt.Println()
}

type bannerLine struct {
	text  string
	color string
}

func textRow(margin string, line bannerLine, width int) string {
	left := (width - len(line.text)) / 2
	right := width - len(line.text) - left
	return margin + edgeColor + "│" + reset +
		strings.Repeat(" ", left) + line.color + line.text + reset + strings.Repeat(" ", right) +
		edgeColor + "│" + reset
}

func consoleWidth() int {
	if columns, err := strconv.Atoi(os.Getenv("COLUMNS")); err == nil && columns > 0 {
		return columns
	}
	return fallbackWidth
}

func enableWindowsConsole() {
	if runtime.GOOS != "windows" {
		return
	}
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getStdHandle := kernel32.NewProc("GetStdHandle")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	setOutputCP := kernel32.NewProc("SetConsoleOutputCP")

	stdout, _, _ := getStdHandle.Call(uintptr(stdOutputHandle))
	var mode uint32
	ok, _, _ := getConsoleMode.Call(stdout, uintptr(unsafe.Pointer(&mode)))
	if ok != 0 {
		setConsoleMode.Call(stdout, uintptr(mode|enableVirtualTerminalHandling))
	}
	setOutputCP.Call(uintptr(utf8CodePage))
}
