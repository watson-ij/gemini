package ui

import "github.com/charmbracelet/bubbles/key"

// KeyMap contains all the key bindings for the application
type KeyMap struct {
	// Navigation
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Home     key.Binding
	End      key.Binding

	// Link navigation
	NextLink     key.Binding
	PrevLink     key.Binding
	FollowLink   key.Binding
	NumberedLink key.Binding // For 1-9, 0 to follow numbered links

	// URL navigation
	FocusAddress key.Binding
	Navigate     key.Binding
	Back         key.Binding
	Forward      key.Binding
	Reload       key.Binding
	GoHome       key.Binding

	// Tabs
	NewTab    key.Binding
	CloseTab  key.Binding
	NextTab   key.Binding
	PrevTab   key.Binding
	JumpTab1  key.Binding
	JumpTab2  key.Binding
	JumpTab3  key.Binding
	JumpTab4  key.Binding
	JumpTab5  key.Binding
	JumpTab6  key.Binding
	JumpTab7  key.Binding
	JumpTab8  key.Binding
	JumpTab9  key.Binding

	// Bookmarks & History
	BookmarkPage   key.Binding
	ToggleSidebar  key.Binding
	ShowHistory    key.Binding

	// Other
	Find  key.Binding
	Help  key.Binding
	Quit  key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Navigation
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("→/l", "right"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "shift+space"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "space"),
			key.WithHelp("pgdn/space", "page down"),
		),
		Home: key.NewBinding(
			key.WithKeys("g", "home"),
			key.WithHelp("g/home", "top"),
		),
		End: key.NewBinding(
			key.WithKeys("G", "end"),
			key.WithHelp("G/end", "bottom"),
		),

		// Link navigation
		NextLink: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next link"),
		),
		PrevLink: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev link"),
		),
		FollowLink: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "follow link"),
		),

		// URL navigation
		FocusAddress: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp("ctrl+l", "address bar"),
		),
		Navigate: key.NewBinding(
			key.WithKeys("ctrl+enter"),
			key.WithHelp("ctrl+enter", "navigate"),
		),
		Back: key.NewBinding(
			key.WithKeys("alt+left"),
			key.WithHelp("alt+←", "back"),
		),
		Forward: key.NewBinding(
			key.WithKeys("alt+right"),
			key.WithHelp("alt+→", "forward"),
		),
		Reload: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "reload"),
		),
		GoHome: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", "home"),
		),

		// Tabs
		NewTab: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "new tab"),
		),
		CloseTab: key.NewBinding(
			key.WithKeys("ctrl+w"),
			key.WithHelp("ctrl+w", "close tab"),
		),
		NextTab: key.NewBinding(
			key.WithKeys("ctrl+tab"),
			key.WithHelp("ctrl+tab", "next tab"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("ctrl+shift+tab"),
			key.WithHelp("ctrl+shift+tab", "prev tab"),
		),
		JumpTab1: key.NewBinding(key.WithKeys("ctrl+1")),
		JumpTab2: key.NewBinding(key.WithKeys("ctrl+2")),
		JumpTab3: key.NewBinding(key.WithKeys("ctrl+3")),
		JumpTab4: key.NewBinding(key.WithKeys("ctrl+4")),
		JumpTab5: key.NewBinding(key.WithKeys("ctrl+5")),
		JumpTab6: key.NewBinding(key.WithKeys("ctrl+6")),
		JumpTab7: key.NewBinding(key.WithKeys("ctrl+7")),
		JumpTab8: key.NewBinding(key.WithKeys("ctrl+8")),
		JumpTab9: key.NewBinding(key.WithKeys("ctrl+9")),

		// Bookmarks & History
		BookmarkPage: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "bookmark"),
		),
		ToggleSidebar: key.NewBinding(
			key.WithKeys("ctrl+b"),
			key.WithHelp("ctrl+b", "bookmarks"),
		),
		ShowHistory: key.NewBinding(
			key.WithKeys("ctrl+shift+h"),
			key.WithHelp("ctrl+shift+h", "history"),
		),

		// Other
		Find: key.NewBinding(
			key.WithKeys("ctrl+f", "/"),
			key.WithHelp("ctrl+f", "find"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "ctrl+q"),
			key.WithHelp("ctrl+q", "quit"),
		),
	}
}

// ShortHelp returns a short help text
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Help,
		k.Quit,
		k.FocusAddress,
		k.ToggleSidebar,
	}
}

// FullHelp returns all key bindings for the help screen
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PageUp, k.PageDown},
		{k.Home, k.End, k.NextLink, k.PrevLink},
		{k.FocusAddress, k.Back, k.Forward, k.Reload},
		{k.NewTab, k.CloseTab, k.NextTab, k.BookmarkPage},
		{k.Find, k.Help, k.Quit},
	}
}
