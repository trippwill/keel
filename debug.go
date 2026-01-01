package keel

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// DefaultDebugProvider returns a debug content provider that fits the content box.
// It truncates lines to ContentWidth and caps the number of lines to ContentHeight.
func DefaultDebugProvider[KID KeelID](id KID, info FrameInfo) (string, error) {
	if info.ContentWidth <= 0 || info.ContentHeight <= 0 {
		return "", nil
	}

	compact := formatCompactDebug(id, info)
	if info.ContentHeight == 1 {
		return truncateDebugLine(compact, info.ContentWidth), nil
	}
	if info.ContentHeight == 2 {
		lines := []string{
			truncateDebugLine(fmt.Sprintf("id:%v", id), info.ContentWidth),
			truncateDebugLine(compact, info.ContentWidth),
		}
		return strings.Join(lines, "\n"), nil
	}

	lines := []string{
		fmt.Sprintf("id:%v", id),
		fmt.Sprintf("alloc:%dx%d", info.Width, info.Height),
		fmt.Sprintf("frame:%dx%d", info.FrameWidth, info.FrameHeight),
		fmt.Sprintf("content:%dx%d", info.ContentWidth, info.ContentHeight),
		fmt.Sprintf("fit:%s", info.Fit.String()),
	}

	maxLines := min(info.ContentHeight, len(lines))

	for i := range maxLines {
		lines[i] = truncateDebugLine(lines[i], info.ContentWidth)
	}

	return strings.Join(lines[:maxLines], "\n"), nil
}

func truncateDebugLine(s string, width int) string {
	if width <= 0 {
		return ""
	}
	return ansi.Truncate(s, width, "")
}

func formatCompactDebug[KID KeelID](id KID, info FrameInfo) string {
	parts := []string{
		fmt.Sprintf("id:%v", id),
		fmt.Sprintf("a:%dx%d", info.Width, info.Height),
		fmt.Sprintf("f:%dx%d", info.FrameWidth, info.FrameHeight),
		fmt.Sprintf("c:%dx%d", info.ContentWidth, info.ContentHeight),
		fmt.Sprintf("ft:%s", info.Fit.String()),
	}

	var b strings.Builder
	b.WriteString(parts[0])
	for _, part := range parts[1:] {
		b.WriteString("|")
		b.WriteString(part)
	}
	return b.String()
}
