# TUI Refactoring Plan for cmd/interactive.go

## Current State Analysis

The `cmd/interactive.go` file implements a sophisticated terminal user interface (TUI) using Bubble Tea. Currently, it uses a monolithic approach with:
- A single large `model` struct containing all UI state
- Large `Update` and `View` methods handling all interactions
- All UI logic centralized in one place

**Current metrics:**
- Single model struct with ~20 fields
- Update method likely 100+ lines
- View method likely 100+ lines
- High cyclomatic complexity

## Proposed Component-Based Architecture

### 1. Component Breakdown

The TUI can be broken down into these logical components:

#### a) HeaderComponent
```go
type HeaderComponent struct {
    title        string
    characterName string
    sessionInfo  string
    styles       headerStyles
}
```
- Manages the top banner/header display
- Shows current character and session information

#### b) MessageListComponent
```go
type MessageListComponent struct {
    messages     []Message
    viewport     viewport.Model
    selectedIdx  int
    styles       messageStyles
}
```
- Manages the scrollable message history
- Handles message selection and scrolling
- Formats messages with appropriate styling

#### c) InputAreaComponent
```go
type InputAreaComponent struct {
    textarea     textarea.Model
    isProcessing bool
    spinner      spinner.Model
    styles       inputStyles
}
```
- Manages user input
- Shows processing state with spinner
- Handles text input events

#### d) StatusBarComponent
```go
type StatusBarComponent struct {
    connectionStatus string
    cacheMetrics     CacheMetrics
    lastError        error
    styles           statusStyles
}
```
- Shows connection status
- Displays cache hit/miss metrics
- Shows error messages

### 2. Implementation Strategy

#### Phase 1: Define Component Interfaces
```go
type Component interface {
    Init() tea.Cmd
    Update(tea.Msg) (Component, tea.Cmd)
    View() string
}
```

#### Phase 2: Extract Components
1. Start with the simplest component (StatusBar)
2. Move relevant fields and logic from main model
3. Create focused Update and View methods
4. Test component in isolation

#### Phase 3: Compose Main Model
```go
type interactiveModel struct {
    header      HeaderComponent
    messageList MessageListComponent
    inputArea   InputAreaComponent
    statusBar   StatusBarComponent
    
    // Shared state
    bot         *services.CharacterBot
    character   *models.Character
    session     *repository.Session
}
```

#### Phase 4: Route Messages
```go
func (m interactiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd
    
    // Route to appropriate component
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if m.inputArea.HasFocus() {
            newInput, cmd := m.inputArea.Update(msg)
            m.inputArea = newInput.(InputAreaComponent)
            cmds = append(cmds, cmd)
        }
    // ... other routing logic
    }
    
    return m, tea.Batch(cmds...)
}
```

### 3. Benefits of Component-Based Approach

1. **Separation of Concerns**: Each component manages its own state and rendering
2. **Testability**: Components can be tested in isolation
3. **Reusability**: Components can be reused in other TUI applications
4. **Maintainability**: Changes to one component don't affect others
5. **Reduced Complexity**: Each component has focused, simple logic

### 4. Migration Path

To minimize disruption:

1. **Keep existing functionality**: Don't change features during refactoring
2. **Incremental approach**: Extract one component at a time
3. **Maintain backwards compatibility**: Keep the same command interface
4. **Test continuously**: Ensure each step maintains functionality

### 5. Example Component Implementation

Here's how the StatusBarComponent might look:

```go
package components

import (
    "fmt"
    "github.com/charmbracelet/bubbles/tea"
    "github.com/charmbracelet/lipgloss"
)

type StatusBarComponent struct {
    width            int
    connectionStatus string
    cacheHits        int
    cacheMisses      int
    lastError        error
    styles           StatusBarStyles
}

type StatusBarStyles struct {
    container   lipgloss.Style
    connected   lipgloss.Style
    disconnected lipgloss.Style
    metrics     lipgloss.Style
    error       lipgloss.Style
}

func NewStatusBar(width int) StatusBarComponent {
    return StatusBarComponent{
        width: width,
        connectionStatus: "connected",
        styles: StatusBarStyles{
            container:    lipgloss.NewStyle().Background(lipgloss.Color("#3c3836")),
            connected:    lipgloss.NewStyle().Foreground(lipgloss.Color("#b8bb26")),
            disconnected: lipgloss.NewStyle().Foreground(lipgloss.Color("#fb4934")),
            metrics:      lipgloss.NewStyle().Foreground(lipgloss.Color("#83a598")),
            error:        lipgloss.NewStyle().Foreground(lipgloss.Color("#fb4934")),
        },
    }
}

func (s StatusBarComponent) Update(msg tea.Msg) (StatusBarComponent, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        s.width = msg.Width
    case ConnectionStatusMsg:
        s.connectionStatus = msg.Status
    case CacheMetricsMsg:
        s.cacheHits = msg.Hits
        s.cacheMisses = msg.Misses
    case ErrorMsg:
        s.lastError = msg.Error
    }
    return s, nil
}

func (s StatusBarComponent) View() string {
    statusStyle := s.styles.connected
    if s.connectionStatus != "connected" {
        statusStyle = s.styles.disconnected
    }
    
    status := statusStyle.Render(fmt.Sprintf("● %s", s.connectionStatus))
    
    metrics := s.styles.metrics.Render(
        fmt.Sprintf("Cache: %d/%d (%.1f%%)", 
            s.cacheHits, 
            s.cacheHits+s.cacheMisses,
            float64(s.cacheHits)/float64(s.cacheHits+s.cacheMisses)*100,
        ),
    )
    
    content := lipgloss.JoinHorizontal(lipgloss.Left, status, " | ", metrics)
    
    if s.lastError != nil {
        errorText := s.styles.error.Render(fmt.Sprintf(" | Error: %v", s.lastError))
        content = lipgloss.JoinHorizontal(lipgloss.Left, content, errorText)
    }
    
    return s.styles.container.Width(s.width).Render(content)
}
```

### 6. Testing Strategy

Each component should have unit tests:

```go
func TestStatusBarComponent(t *testing.T) {
    statusBar := NewStatusBar(80)
    
    // Test connection status update
    statusBar, _ = statusBar.Update(ConnectionStatusMsg{Status: "disconnected"})
    view := statusBar.View()
    assert.Contains(t, view, "disconnected")
    
    // Test cache metrics update
    statusBar, _ = statusBar.Update(CacheMetricsMsg{Hits: 10, Misses: 5})
    view = statusBar.View()
    assert.Contains(t, view, "66.7%")
}
```

## Conclusion

This component-based refactoring will transform the monolithic TUI into a modular, maintainable architecture. While it requires significant effort, the long-term benefits in terms of maintainability, testability, and developer experience make it a worthwhile investment for a developer tool where the TUI is a core feature.