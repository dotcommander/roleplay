package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InputArea handles user input and processing state
type InputArea struct {
	textarea     textarea.Model
	spinner      spinner.Model
	isProcessing bool
	error        error
	focused      bool
	width        int
	styles       inputAreaStyles
}

type inputAreaStyles struct {
	error      lipgloss.Style
	processing lipgloss.Style
}

// NewInputArea creates a new input area component
func NewInputArea(width int) *InputArea {
	ta := textarea.New()
	ta.Placeholder = ""
	ta.Prompt = "> "  // Simple clean prompt like aider-go
	ta.CharLimit = 500
	// Minimal width adjustment for clean single-line footer
	ta.SetWidth(width - 3)
	ta.SetHeight(1)
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	// Minimal styling for clean look
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(GruvboxFg)
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(GruvboxFg)
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(GruvboxGray)
	ta.BlurredStyle.Prompt = lipgloss.NewStyle().Foreground(GruvboxGray)
	ta.BlurredStyle.Text = lipgloss.NewStyle().Foreground(GruvboxGray)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(GruvboxAqua)

	return &InputArea{
		textarea: ta,
		spinner:  s,
		width:    width,
		styles: inputAreaStyles{
			error: lipgloss.NewStyle().
				Foreground(GruvboxRed).
				Bold(true),
			processing: lipgloss.NewStyle().
				Foreground(GruvboxGray),
		},
	}
}

// Update handles messages for the input area
func (i *InputArea) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		i.width = msg.Width
		// Minimal width adjustment for clean single-line footer
		i.textarea.SetWidth(msg.Width - 3)

	case ProcessingStateMsg:
		i.isProcessing = msg.IsProcessing
		if msg.IsProcessing {
			i.textarea.Blur()
		} else {
			i.textarea.Focus()
		}

	case spinner.TickMsg:
		if i.isProcessing {
			var cmd tea.Cmd
			i.spinner, cmd = i.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Update textarea if not processing
	if !i.isProcessing && i.focused {
		var cmd tea.Cmd
		i.textarea, cmd = i.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

// View renders the input area as a clean sticky footer
func (i *InputArea) View() string {
	if i.error != nil {
		return i.styles.error.Render(fmt.Sprintf("Error: %v", i.error))
	}

	if i.isProcessing {
		spinnerText := i.styles.processing.Render("thinking...")
		return fmt.Sprintf("%s %s", i.spinner.View(), spinnerText)
	}

	// Clean single line footer, no extra padding
	return i.textarea.View()
}

// Focus sets the input area as focused
func (i *InputArea) Focus() {
	i.focused = true
	if !i.isProcessing {
		i.textarea.Focus()
	}
}

// Blur removes focus from the input area
func (i *InputArea) Blur() {
	i.focused = false
	i.textarea.Blur()
}

// IsFocused returns whether the input area is focused
func (i *InputArea) IsFocused() bool {
	return i.focused
}

// SetSize updates the size of the input area
func (i *InputArea) SetSize(width, _ int) {
	i.width = width
	// Minimal width adjustment for clean single-line footer
	i.textarea.SetWidth(width - 3)
}

// Value returns the current input value
func (i *InputArea) Value() string {
	return i.textarea.Value()
}

// SetValue sets the input value
func (i *InputArea) SetValue(s string) {
	i.textarea.SetValue(s)
}

// Reset clears the input
func (i *InputArea) Reset() {
	i.textarea.Reset()
}

// CursorEnd moves cursor to end of input
func (i *InputArea) CursorEnd() {
	i.textarea.CursorEnd()
}

// Init returns initialization commands
func (i *InputArea) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		i.spinner.Tick,
	)
}
