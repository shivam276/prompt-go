package tui

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const maxClipboardSize = 100 * 1024 // 100KB limit for OSC 52

type copyFeedbackMsg struct{}
type saveSuccessMsg struct{ path string }
type saveErrorMsg struct{ err error }

// CopyToClipboard returns a tea.Cmd that copies content to clipboard using OSC 52
func CopyToClipboard(content string) tea.Cmd {
	return func() tea.Msg {
		// Truncate if too large
		if len(content) > maxClipboardSize {
			content = content[:maxClipboardSize] + "\n\n[... truncated for clipboard]"
		}

		// Encode to base64
		encoded := base64.StdEncoding.EncodeToString([]byte(content))

		// Generate OSC 52 sequence (using both ST terminators for compatibility)
		osc52 := fmt.Sprintf("\033]52;c;%s\033\\", encoded)

		// Print the escape sequence to terminal
		fmt.Print(osc52)

		return copyFeedbackMsg{}
	}
}

// SaveToFile saves content to a file in the current directory
func SaveToFile(content string) tea.Cmd {
	return func() tea.Msg {
		// Generate filename with timestamp
		filename := fmt.Sprintf("prompt_%d.txt", time.Now().Unix())

		// Save to current directory
		path := filepath.Join(".", filename)
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			return saveErrorMsg{err: err}
		}

		return saveSuccessMsg{path: path}
	}
}

// ShowCopyFeedback returns a tea.Cmd that shows feedback for 2 seconds
func ShowCopyFeedback() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return hideCopyFeedbackMsg{}
	})
}

type hideCopyFeedbackMsg struct{}
