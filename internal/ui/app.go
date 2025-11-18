package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/watson-ij/gemini/internal/parser"
	"github.com/watson-ij/gemini/internal/protocol"
)

// AppMode represents the current mode of the application
type AppMode int

const (
	// ModeBrowse is the normal browsing mode
	ModeBrowse AppMode = iota

	// ModeAddressBar is when the address bar is focused
	ModeAddressBar

	// ModeHelp is when the help screen is displayed
	ModeHelp

	// ModeBookmarks is when the bookmarks sidebar is displayed
	ModeBookmarks
)

// Model is the main application model
type Model struct {
	// UI state
	mode       AppMode
	width      int
	height     int
	ready      bool
	err        error
	statusMsg  string

	// Components
	viewport    viewport.Model
	addressBar  textinput.Model
	help        help.Model
	keys        KeyMap

	// Content
	currentURL  string
	document    *parser.Document
	rawContent  string
	loading     bool
	selectedLink int  // Currently selected link index (-1 = none)

	// Protocol
	client *protocol.Client

	// Navigation history
	history  []string  // URLs visited
	historyPos int     // Current position in history

	// Styles
	styles Styles
}

// Styles contains all the lipgloss styles for the UI
type Styles struct {
	TitleBar       lipgloss.Style
	AddressBar     lipgloss.Style
	AddressBarFocused lipgloss.Style
	StatusBar      lipgloss.Style
	StatusBarError lipgloss.Style
	StatusBarInfo  lipgloss.Style
	HelpBar        lipgloss.Style
}

// DefaultStyles returns the default styles
func DefaultStyles() Styles {
	titleBarStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Bold(true).
		Padding(0, 1)

	addressBarStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	addressBarFocusedStyle := addressBarStyle.Copy().
		BorderForeground(lipgloss.Color("69"))

	statusBarStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	statusBarErrorStyle := statusBarStyle.Copy().
		Background(lipgloss.Color("196"))

	statusBarInfoStyle := statusBarStyle.Copy().
		Background(lipgloss.Color("33"))

	helpBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(0, 1)

	return Styles{
		TitleBar:          titleBarStyle,
		AddressBar:        addressBarStyle,
		AddressBarFocused: addressBarFocusedStyle,
		StatusBar:         statusBarStyle,
		StatusBarError:    statusBarErrorStyle,
		StatusBarInfo:     statusBarInfoStyle,
		HelpBar:           helpBarStyle,
	}
}

// NewModel creates a new application model
func NewModel(startURL string) Model {
	// Create text input for address bar
	ti := textinput.New()
	ti.Placeholder = "gemini://..."
	ti.Prompt = "URL: "
	ti.CharLimit = 1024
	if startURL != "" {
		ti.SetValue(startURL)
	}

	// Create Gemini client
	client := protocol.NewClient()

	// Create viewport
	vp := viewport.New(80, 20)
	vp.KeyMap = viewport.KeyMap{
		PageDown: key.NewBinding(key.WithKeys("pgdown", "space")),
		PageUp:   key.NewBinding(key.WithKeys("pgup", "shift+space")),
		Down:     key.NewBinding(key.WithKeys("j", "down")),
		Up:       key.NewBinding(key.WithKeys("k", "up")),
	}

	m := Model{
		mode:         ModeBrowse,
		addressBar:   ti,
		viewport:     vp,
		help:         help.New(),
		keys:         DefaultKeyMap(),
		client:       client,
		currentURL:   startURL,
		selectedLink: -1,
		history:      []string{},
		historyPos:   -1,
		styles:       DefaultStyles(),
	}

	return m
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	// If we have a start URL, load it
	if m.currentURL != "" {
		return m.loadURL(m.currentURL)
	}
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			// Initialize viewport with correct size
			headerHeight := 4 // Title + address bar
			footerHeight := 2 // Status bar
			m.viewport = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
			m.viewport.KeyMap = viewport.KeyMap{
				PageDown: key.NewBinding(key.WithKeys("pgdown", "space")),
				PageUp:   key.NewBinding(key.WithKeys("pgup", "shift+space")),
				Down:     key.NewBinding(key.WithKeys("j", "down")),
				Up:       key.NewBinding(key.WithKeys("k", "up")),
			}
			m.ready = true
		} else {
			headerHeight := 4
			footerHeight := 2
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - headerHeight - footerHeight
		}

		if m.document != nil {
			m.renderDocument()
		}

	case pageLoadedMsg:
		m.loading = false
		m.document = msg.doc
		m.rawContent = msg.raw
		m.selectedLink = -1
		m.statusMsg = fmt.Sprintf("Loaded %d lines, %d links", msg.doc.LineCount(), msg.doc.LinkCount())
		m.renderDocument()

	case errorMsg:
		m.loading = false
		m.err = msg.err
		m.statusMsg = fmt.Sprintf("Error: %v", msg.err)

	case tea.KeyMsg:
		// Handle mode-specific keys first
		switch m.mode {
		case ModeAddressBar:
			return m.updateAddressBar(msg)

		case ModeHelp:
			if key.Matches(msg, m.keys.Help) || msg.String() == "esc" {
				m.mode = ModeBrowse
			}
			return m, nil
		}

		// Global keys (browse mode)
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Help):
			m.mode = ModeHelp
			return m, nil

		case key.Matches(msg, m.keys.FocusAddress):
			m.mode = ModeAddressBar
			m.addressBar.Focus()
			return m, textinput.Blink

		case key.Matches(msg, m.keys.FollowLink):
			if m.document != nil && m.selectedLink >= 0 && m.selectedLink < len(m.document.Links) {
				link := m.document.Links[m.selectedLink]
				url := link.Link.URL
				// Resolve relative URLs
				if !strings.HasPrefix(url, "gemini://") {
					url = m.resolveURL(url)
				}
				return m, m.loadURL(url)
			}

		case key.Matches(msg, m.keys.NextLink):
			if m.document != nil && m.document.LinkCount() > 0 {
				m.selectedLink = (m.selectedLink + 1) % m.document.LinkCount()
				m.renderDocument()
				// Scroll to show the selected link
				lineNum := parser.GetLineForLink(m.document, m.selectedLink)
				if lineNum >= 0 {
					m.viewport.GotoTop()
					m.viewport.LineDown(lineNum)
				}
			}

		case key.Matches(msg, m.keys.PrevLink):
			if m.document != nil && m.document.LinkCount() > 0 {
				m.selectedLink--
				if m.selectedLink < 0 {
					m.selectedLink = m.document.LinkCount() - 1
				}
				m.renderDocument()
				lineNum := parser.GetLineForLink(m.document, m.selectedLink)
				if lineNum >= 0 {
					m.viewport.GotoTop()
					m.viewport.LineDown(lineNum)
				}
			}

		case key.Matches(msg, m.keys.Reload):
			if m.currentURL != "" {
				return m, m.loadURL(m.currentURL)
			}

		case key.Matches(msg, m.keys.Back):
			if m.historyPos > 0 {
				m.historyPos--
				url := m.history[m.historyPos]
				m.currentURL = url
				m.addressBar.SetValue(url)
				return m, m.loadURL(url)
			}

		case key.Matches(msg, m.keys.Forward):
			if m.historyPos < len(m.history)-1 {
				m.historyPos++
				url := m.history[m.historyPos]
				m.currentURL = url
				m.addressBar.SetValue(url)
				return m, m.loadURL(url)
			}

		case key.Matches(msg, m.keys.Home):
			m.viewport.GotoTop()

		case key.Matches(msg, m.keys.End):
			m.viewport.GotoBottom()
		}
	}

	// Update viewport
	if m.mode == ModeBrowse {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// updateAddressBar handles updates when in address bar mode
func (m Model) updateAddressBar(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "enter":
		m.mode = ModeBrowse
		m.addressBar.Blur()
		url := m.addressBar.Value()
		if url != "" {
			return m, m.loadURL(url)
		}
		return m, nil

	case "esc":
		m.mode = ModeBrowse
		m.addressBar.Blur()
		m.addressBar.SetValue(m.currentURL)
		return m, nil
	}

	m.addressBar, cmd = m.addressBar.Update(msg)
	return m, cmd
}

// View renders the application
func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	// Different views for different modes
	switch m.mode {
	case ModeHelp:
		return m.helpView()
	default:
		return m.browseView()
	}
}

// browseView renders the main browsing view
func (m Model) browseView() string {
	// Title bar
	title := m.styles.TitleBar.Render("📡 Gemini Browser")

	// Address bar
	addressStyle := m.styles.AddressBar
	if m.mode == ModeAddressBar {
		addressStyle = m.styles.AddressBarFocused
	}
	address := addressStyle.Render(m.addressBar.View())

	// Content
	content := m.viewport.View()

	// Status bar
	statusStyle := m.styles.StatusBar
	if m.err != nil {
		statusStyle = m.styles.StatusBarError
	} else if m.loading {
		statusStyle = m.styles.StatusBarInfo
	}

	statusLeft := m.statusMsg
	if m.loading {
		statusLeft = "Loading..."
	}

	linkCount := 0
	if m.document != nil {
		linkCount = m.document.LinkCount()
	}

	statusRight := fmt.Sprintf("Link %d/%d | %d%% | ? for help",
		m.selectedLink+1,
		linkCount,
		int(float64(m.viewport.YOffset)/float64(max(1, len(strings.Split(m.viewport.View(), "\n"))-1))*100))

	statusPadding := m.width - lipgloss.Width(statusLeft) - lipgloss.Width(statusRight)
	if statusPadding < 0 {
		statusPadding = 0
	}

	status := statusStyle.Render(statusLeft + strings.Repeat(" ", statusPadding) + statusRight)

	// Help
	helpText := m.styles.HelpBar.Render("↑/↓: scroll | tab: next link | enter: follow | ctrl+l: address | ctrl+q: quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		address,
		content,
		status,
		helpText,
	)
}

// helpView renders the help screen
func (m Model) helpView() string {
	title := m.styles.TitleBar.Render("Help - Press ? or ESC to close")

	helpContent := `
Gemini Browser - Keyboard Shortcuts

Navigation:
  ↑/k, ↓/j       Scroll up/down
  PgUp/PgDn      Page up/down
  Space          Page down
  g, G           Go to top/bottom
  h, l           Scroll left/right

Links:
  Tab            Next link
  Shift+Tab      Previous link
  Enter          Follow selected link

URL Navigation:
  Ctrl+L         Focus address bar
  Ctrl+R         Reload current page
  Alt+←          Go back
  Alt+→          Go forward

Other:
  Ctrl+B         Toggle bookmarks (TODO)
  Ctrl+D         Bookmark page (TODO)
  Ctrl+F         Find in page (TODO)
  ?              Show this help
  Ctrl+Q         Quit

Press ? or ESC to close this help screen.
`

	content := lipgloss.NewStyle().
		Padding(1, 2).
		Render(helpContent)

	return lipgloss.JoinVertical(lipgloss.Left, title, content)
}

// renderDocument renders the current document to the viewport
func (m *Model) renderDocument() {
	if m.document == nil {
		m.viewport.SetContent("No document loaded")
		return
	}

	renderer := parser.NewRenderer(&parser.RenderOptions{
		Width:           m.viewport.Width,
		NumberLinks:     true,
		HighlightedLink: m.selectedLink,
		ColorScheme:     parser.DefaultColorScheme(),
	})

	content := renderer.Render(m.document)
	m.viewport.SetContent(content)
}

// loadURL loads a URL and returns a command
func (m *Model) loadURL(url string) tea.Cmd {
	// Add to history
	if url != m.currentURL {
		// Trim history after current position
		m.history = m.history[:m.historyPos+1]
		m.history = append(m.history, url)
		m.historyPos = len(m.history) - 1
	}

	m.currentURL = url
	m.addressBar.SetValue(url)
	m.loading = true
	m.err = nil

	return func() tea.Msg {
		resp, err := m.client.Get(url)
		if err != nil {
			return errorMsg{err: err}
		}
		defer resp.Close()

		if !resp.Status.IsSuccess() {
			return errorMsg{err: fmt.Errorf("status %d: %s", resp.Status, resp.Meta)}
		}

		body, err := resp.ReadBody()
		if err != nil {
			return errorMsg{err: err}
		}

		doc, err := parser.ParseString(string(body))
		if err != nil {
			return errorMsg{err: err}
		}

		return pageLoadedMsg{
			doc: doc,
			raw: string(body),
		}
	}
}

// resolveURL resolves a relative URL against the current URL
func (m *Model) resolveURL(relativeURL string) string {
	// Simple implementation - would need proper URL resolution
	if strings.HasPrefix(relativeURL, "/") {
		// Absolute path
		parts := strings.SplitN(m.currentURL, "/", 4)
		if len(parts) >= 3 {
			return parts[0] + "//" + parts[2] + relativeURL
		}
	}
	// For now, just return as-is
	return relativeURL
}

// Messages
type pageLoadedMsg struct {
	doc *parser.Document
	raw string
}

type errorMsg struct {
	err error
}

// Helper function
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
