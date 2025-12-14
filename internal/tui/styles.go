package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("205") // Pink/Purple
	secondaryColor = lipgloss.Color("240") // Gray
	successColor   = lipgloss.Color("42")  // Green
	errorColor     = lipgloss.Color("196") // Red
	accentColor    = lipgloss.Color("86")  // Cyan
)

// TitleStyle returns the style for the app title
func TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		Padding(0, 1)
}

// SubtitleStyle returns the style for subtitle text
func SubtitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(secondaryColor).
		Italic(true)
}

// FieldLabelStyle returns the style for field labels
func FieldLabelStyle(focused bool) lipgloss.Style {
	if focused {
		return lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)
	}
	return lipgloss.NewStyle().
		Foreground(secondaryColor)
}

// ErrorStyle returns the style for error messages
func ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		Padding(0, 1)
}

// TipStyle returns the style for tips
func TipStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(accentColor).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentColor)
}

// HelpStyle returns the style for help text
func HelpStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(secondaryColor).
		Italic(true)
}

// StatusBarStyle returns the style for status bar
func StatusBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true)
}

// ContainerStyle returns the style for content containers
func ContainerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(0, 1)
}

// CenterView centers content both horizontally and vertically in the terminal
func CenterView(content string, width, height int) string {
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// CursorStyle returns the style for cursor (same color as header)
func CursorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(primaryColor)
}
