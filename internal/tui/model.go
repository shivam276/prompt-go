package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"promptgo/internal/enhancer"
)

type hideSaveFeedbackMsg struct{}
type printAndExitMsg struct{}

type appState int

const (
	stateInput appState = iota
	stateResult
)

type focusedField int

const (
	fieldTask focusedField = iota
	fieldDetails
	fieldSecretWord
)

type Model struct {
	// State
	state   appState
	focused focusedField

	// Input fields (inputView)
	taskInput    textarea.Model
	detailsInput textarea.Model
	secretInput  textinput.Model

	// Result data (resultView)
	enhancedPrompt string
	tip            string
	resultViewport viewport.Model

	// UI state
	width        int
	height       int
	err          string
	copyFeedback bool
	saveFeedback string
}

// NewModel creates a new TUI model
func NewModel() Model {
	// Configure task textarea
	task := textarea.New()
	task.Placeholder = "Describe what you want to build..."
	task.Focus()
	task.CharLimit = 5000
	task.SetHeight(6)
	task.SetWidth(80)
	task.ShowLineNumbers = false
	task.FocusedStyle.CursorLine = CursorStyle()
	task.Cursor.Style = CursorStyle()

	// Configure details textarea
	details := textarea.New()
	details.Placeholder = "Additional context like libraries, architecture, constraints... (optional)"
	details.CharLimit = 5000
	details.SetHeight(6)
	details.SetWidth(80)
	details.ShowLineNumbers = false
	details.FocusedStyle.CursorLine = CursorStyle()
	details.Cursor.Style = CursorStyle()

	// Configure secret word input
	secret := textinput.New()
	secret.Placeholder = "your secret word"
	secret.CharLimit = 50
	secret.Width = 80
	secret.Cursor.Style = CursorStyle()
	secret.TextStyle = CursorStyle()

	// Create viewport for results
	vp := viewport.New(80, 20)

	return Model{
		state:          stateInput,
		focused:        fieldTask,
		taskInput:      task,
		detailsInput:   details,
		secretInput:    secret,
		resultViewport: vp,
		width:          80,
		height:         24,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// State-specific handling
		switch m.state {
		case stateInput:
			return m.updateInput(msg)
		case stateResult:
			return m.updateResult(msg)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Set a max content width for better centering
		maxContentWidth := 80
		contentWidth := msg.Width - 4
		if contentWidth > maxContentWidth {
			contentWidth = maxContentWidth
		}

		// Update input fields
		m.taskInput.SetWidth(contentWidth)
		m.detailsInput.SetWidth(contentWidth)
		m.secretInput.Width = contentWidth

		// Update result viewport
		m.resultViewport.Width = contentWidth
		m.resultViewport.Height = msg.Height - 15 // Leave room for header/footer

		return m, nil

	case copyFeedbackMsg:
		m.copyFeedback = true
		return m, ShowCopyFeedback()

	case hideCopyFeedbackMsg:
		m.copyFeedback = false
		return m, nil

	case saveSuccessMsg:
		m.saveFeedback = fmt.Sprintf("✓ Saved to %s", msg.path)
		return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
			return hideSaveFeedbackMsg{}
		})

	case saveErrorMsg:
		m.saveFeedback = fmt.Sprintf("✗ Error: %v", msg.err)
		return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
			return hideSaveFeedbackMsg{}
		})

	case hideSaveFeedbackMsg:
		m.saveFeedback = ""
		return m, nil

	case printAndExitMsg:
		// Print the enhanced prompt to stdout, then quit
		fmt.Println("\n" + strings.Repeat("=", 80))
		fmt.Println("ENHANCED PROMPT")
		fmt.Println(strings.Repeat("=", 80) + "\n")
		fmt.Println(m.enhancedPrompt)
		fmt.Println("\n" + strings.Repeat("=", 80))
		return m, tea.Quit
	}

	return m, tea.Batch(cmds...)
}

// updateInput handles input view updates
func (m Model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyTab:
		// Cycle focus forward
		m.blurAll()
		switch m.focused {
		case fieldTask:
			m.focused = fieldDetails
			m.detailsInput.Focus()
		case fieldDetails:
			m.focused = fieldSecretWord
			m.secretInput.Focus()
		case fieldSecretWord:
			m.focused = fieldTask
			m.taskInput.Focus()
		}
		return m, nil

	case tea.KeyShiftTab:
		// Cycle focus backward
		m.blurAll()
		switch m.focused {
		case fieldTask:
			m.focused = fieldSecretWord
			m.secretInput.Focus()
		case fieldDetails:
			m.focused = fieldTask
			m.taskInput.Focus()
		case fieldSecretWord:
			m.focused = fieldDetails
			m.detailsInput.Focus()
		}
		return m, nil

	case tea.KeyCtrlE:
		// Trigger enhancement
		return m.enhance()
	}

	// Delegate to focused field
	switch m.focused {
	case fieldTask:
		m.taskInput, cmd = m.taskInput.Update(msg)
	case fieldDetails:
		m.detailsInput, cmd = m.detailsInput.Update(msg)
	case fieldSecretWord:
		m.secretInput, cmd = m.secretInput.Update(msg)
	}

	return m, cmd
}

// updateResult handles result view updates
func (m Model) updateResult(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "c":
		// Copy to clipboard
		return m, CopyToClipboard(m.enhancedPrompt)

	case "s":
		// Save to file
		return m, SaveToFile(m.enhancedPrompt)

	case "p":
		// Print to terminal and exit
		return m, func() tea.Msg {
			return printAndExitMsg{}
		}

	case "r":
		// Reset to input view
		m.state = stateInput
		m.taskInput.SetValue("")
		m.detailsInput.SetValue("")
		m.secretInput.SetValue("")
		m.focused = fieldTask
		m.err = ""
		m.copyFeedback = false
		m.saveFeedback = ""
		m.taskInput.Focus()
		return m, nil

	case "q":
		return m, tea.Quit
	}

	// Delegate to viewport for scrolling
	m.resultViewport, cmd = m.resultViewport.Update(msg)
	return m, cmd
}

// enhance triggers the enhancement process
func (m Model) enhance() (tea.Model, tea.Cmd) {
	// Validate
	if strings.TrimSpace(m.taskInput.Value()) == "" {
		m.err = "Task is required"
		return m, nil
	}
	if strings.TrimSpace(m.secretInput.Value()) == "" {
		m.err = "Secret word is required"
		return m, nil
	}

	// Enhance
	input := enhancer.Input{
		Task:       m.taskInput.Value(),
		Details:    m.detailsInput.Value(),
		SecretWord: m.secretInput.Value(),
	}
	output := enhancer.Enhance(input)

	// Transition to result view
	m.state = stateResult
	m.enhancedPrompt = output.EnhancedPrompt
	m.tip = output.Tip
	m.resultViewport.SetContent(output.EnhancedPrompt)
	m.resultViewport.GotoTop()
	m.err = ""

	return m, nil
}

// blurAll blurs all input fields
func (m *Model) blurAll() {
	m.taskInput.Blur()
	m.detailsInput.Blur()
	m.secretInput.Blur()
}

// View renders the view
func (m Model) View() string {
	var content string
	switch m.state {
	case stateInput:
		content = m.viewInput()
	case stateResult:
		content = m.viewResult()
	default:
		return ""
	}

	// Center the content in the terminal
	return CenterView(content, m.width, m.height)
}

// viewInput renders the input view
func (m Model) viewInput() string {
	var b strings.Builder

	// Title
	b.WriteString(TitleStyle().Render("🐹 PromptGo"))
	b.WriteString("\n")
	b.WriteString(SubtitleStyle().Render("Stop letting AI write garbage Go code"))
	b.WriteString("\n\n")

	// Task field
	b.WriteString(FieldLabelStyle(m.focused == fieldTask).Render("What do you want to build?"))
	b.WriteString("\n")
	b.WriteString(m.taskInput.View())
	b.WriteString("\n\n")

	// Details field
	b.WriteString(FieldLabelStyle(m.focused == fieldDetails).Render("Additional details (optional):"))
	b.WriteString("\n")
	b.WriteString(m.detailsInput.View())
	b.WriteString("\n\n")

	// Secret word field
	b.WriteString(FieldLabelStyle(m.focused == fieldSecretWord).Render("Secret word:"))
	b.WriteString(" ")
	b.WriteString(m.secretInput.View())
	b.WriteString("\n\n")

	// Error message
	if m.err != "" {
		b.WriteString(ErrorStyle().Render("❌ " + m.err))
		b.WriteString("\n\n")
	}

	// Help
	b.WriteString(HelpStyle().Render("[Tab] Next field   [Shift+Tab] Prev   [Ctrl+E] Enhance   [Ctrl+C] Quit"))
	b.WriteString("\n")

	return b.String()
}

// viewResult renders the result view
func (m Model) viewResult() string {
	var b strings.Builder

	// Title
	b.WriteString(TitleStyle().Render("🐹 PromptGo - Enhanced Prompt"))
	b.WriteString("\n\n")

	// Viewport with enhanced prompt
	b.WriteString(ContainerStyle().Render(m.resultViewport.View()))
	b.WriteString("\n\n")

	// Tip
	tipText := fmt.Sprintf("💡 %s", m.tip)
	b.WriteString(TipStyle().Render(tipText))
	b.WriteString("\n\n")

	// Help
	b.WriteString(HelpStyle().Render("[p] Print & exit (copyable!)   [s] Save to file   [r] Start over   [q] Quit"))
	b.WriteString("\n")

	// Feedback messages
	if m.copyFeedback {
		b.WriteString("\n")
		b.WriteString(StatusBarStyle().Render("✓ Copied to clipboard!"))
	}
	if m.saveFeedback != "" {
		b.WriteString("\n")
		b.WriteString(StatusBarStyle().Render(m.saveFeedback))
	}

	return b.String()
}
