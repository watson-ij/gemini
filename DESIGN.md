# Gemini Terminal Client - Design Document

**Project:** Gemini TUI Client
**Language:** Go
**Version:** 1.0.0
**Last Updated:** 2025-11-18

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Protocol Overview](#protocol-overview)
3. [Architecture](#architecture)
4. [User Interface Design](#user-interface-design)
5. [Core Features](#core-features)
6. [Technical Specifications](#technical-specifications)
7. [Implementation Plan](#implementation-plan)
8. [Data Storage](#data-storage)

---

## Executive Summary

This document outlines the design for a full-featured Gemini protocol terminal client built in Go. The client will provide a modern, power-user-focused TUI (Text User Interface) for browsing Geminispace with features like tabbed browsing, bookmarks, history, theming, and more.

**Design Goals:**
- **Simplicity**: Honor the Gemini protocol's philosophy of simplicity
- **Power User**: Vim-inspired keybindings and efficient navigation
- **Beautiful**: Modern TUI with theming support
- **Feature Rich**: Bookmarks, tabs, history, search, subscriptions
- **Secure**: Proper TLS and certificate management (TOFU)

---

## Protocol Overview

### What is Gemini?

Gemini is a lightweight internet protocol designed as an alternative to HTTP, emphasizing simplicity and privacy. It sits between Gopher and the modern web in complexity.

**Key Characteristics:**
- Application-layer protocol over TCP
- Default port: **1965**
- **TLS required** (v1.2 or higher)
- Single request-response per connection
- Uses custom markup language: **Gemtext**

### Request Format

```
<URL><CR><LF>
```

Example:
```
gemini://example.com/page\r\n
```

### Response Format

```
<STATUS><SPACE><META><CR><LF>
[<CONTENT>]
```

Example:
```
20 text/gemini\r\n
# Welcome to Geminispace!
```

### Status Codes

| Range | Category | Description |
|-------|----------|-------------|
| **1x** | INPUT | Server requests input from client |
| | 10 | Input required |
| | 11 | Sensitive input (password) |
| **2x** | SUCCESS | Request successful |
| | 20 | Success (only defined code) |
| **3x** | REDIRECT | Resource moved |
| | 30 | Temporary redirect |
| | 31 | Permanent redirect |
| **4x** | TEMPORARY FAILURE | Temporary server-side error |
| | 40 | Temporary failure |
| | 41 | Server unavailable |
| | 42 | CGI error |
| | 44 | Slow down (rate limiting) |
| **5x** | PERMANENT FAILURE | Permanent error |
| | 50 | Permanent failure |
| | 51 | Not found |
| | 52 | Gone |
| | 59 | Bad request |
| **6x** | CLIENT CERTIFICATE | Certificate required/invalid |
| | 60 | Client certificate required |
| | 61 | Certificate not authorized |
| | 62 | Certificate not valid |

**Default Behavior:** Undefined codes (e.g., 14, 22) should be treated as their category default (10, 20, etc.)

### Gemtext Format

Gemtext is a line-oriented markup format with 6 line types:

```gemtext
# Heading level 1
## Heading level 2
### Heading level 3

=> gemini://example.com Link text here
=> gemini://example.com/image.jpg

* List item one
* List item two

> Quote text

``` preformatted toggle (optional alt text)
Preformatted content
Code, ASCII art, etc.
```

Regular text line (paragraph)
```

**Line Type Rules:**
- Text lines: Regular paragraphs
- Link lines: Start with `=>` then whitespace, then URL, optionally followed by link text
- Heading lines: Start with 1-3 `#` followed by space
- List items: Start with `*` and space
- Quotes: Start with `>` and optional space
- Preformatted: Toggle with `\`\`\`` (three backticks), optional alt text after

**Compliance:**
- MUST handle: text lines, link lines, preformat toggles
- SHOULD handle: headings, list items, quotes (for better UX)

### TLS & TOFU

**TOFU (Trust On First Use):**
1. First connection: Accept and store certificate
2. Subsequent connections: Verify against stored cert
3. Certificate change: Warn user, require explicit acceptance

**Certificate Requirements:**
- TLS 1.2 or higher
- Self-signed certificates are common and acceptable
- Client certificates for authentication (optional)

---

## Architecture

### Technology Stack

**Core Libraries:**
- **charmbracelet/bubbletea** - TUI framework (Elm architecture)
- **charmbracelet/bubbles** - Reusable UI components
- **charmbracelet/lipgloss** - Terminal styling
- **makeworld-the-better-one/go-gemini** - Gemini protocol (or custom)
- **pelletier/go-toml** - Configuration parsing
- **mattn/go-sqlite3** - History database
- **golang.org/x/text** - Character encoding

### Project Structure

```
gemini-client/
├── main.go                 # Application entry point
├── go.mod
├── go.sum
├── README.md
├── DESIGN.md              # This document
├── LICENSE
│
├── cmd/
│   └── gemini/            # CLI entry point
│       └── main.go
│
├── internal/
│   ├── protocol/          # Gemini protocol implementation
│   │   ├── client.go      # HTTP client, TLS connection
│   │   ├── request.go     # Request builder
│   │   ├── response.go    # Response parser
│   │   ├── status.go      # Status code handling
│   │   └── tofu.go        # Certificate trust management
│   │
│   ├── parser/            # Gemtext parsing
│   │   ├── gemtext.go     # Parse gemtext to AST
│   │   ├── renderer.go    # Render AST for TUI
│   │   └── types.go       # AST node types
│   │
│   ├── storage/           # Data persistence
│   │   ├── bookmarks.go   # Bookmark CRUD operations
│   │   ├── history.go     # History tracking
│   │   ├── cache.go       # Page caching
│   │   ├── config.go      # Configuration management
│   │   └── db.go          # Database initialization
│   │
│   ├── ui/                # TUI components (Bubble Tea)
│   │   ├── app.go         # Main application model
│   │   ├── browser.go     # Page viewing component
│   │   ├── tabs.go        # Tab management
│   │   ├── sidebar.go     # Bookmarks/history sidebar
│   │   ├── addressbar.go  # URL input field
│   │   ├── statusbar.go   # Bottom status bar
│   │   ├── modal.go       # Modal dialogs
│   │   ├── help.go        # Help screen
│   │   └── keys.go        # Keybinding definitions
│   │
│   └── theme/             # Theming system
│       ├── theme.go       # Theme interface and loader
│       ├── default.go     # Default theme
│       └── styles.go      # Style definitions
│
├── pkg/                   # Public packages (if any)
│
├── assets/
│   └── themes/            # Built-in themes
│       ├── default.json
│       ├── dark.json
│       └── monokai.json
│
└── testdata/              # Test fixtures
    └── pages/
        └── *.gmi
```

### Bubble Tea Architecture

The application follows **The Elm Architecture**:

```
┌─────────────────────────────────────────┐
│              Model (State)              │
│  - Current URL                          │
│  - Open tabs                            │
│  - Bookmarks                            │
│  - History                              │
│  - UI state                             │
└─────────────────────────────────────────┘
         │                    ▲
         │                    │
         ▼                    │
┌──────────────┐      ┌──────────────┐
│   View()     │      │   Update()   │
│              │      │              │
│ Render UI    │      │ Handle msgs  │
│ based on     │      │ Update model │
│ model state  │      │ Return cmds  │
└──────────────┘      └──────────────┘
         │                    ▲
         │                    │
         ▼                    │
   Terminal Display      User Input / Events
```

**Key Components:**
- **Model**: Application state (tabs, bookmarks, current page, etc.)
- **Update**: Event handler (keyboard, network responses, etc.)
- **View**: Render function (converts model to terminal output)
- **Commands**: Side effects (network requests, file I/O)

---

## User Interface Design

### Layout

```
┌─────────────────────────────────────────────────────────────────────┐
│ Tab 1: Example    Tab 2: Home    Tab 3: News    [+]           [Help] │
├─────────────────────────────────────────────────────────────────────┤
│ 🔍 gemini://example.com/page                              [🔖] [⚙️]  │
├──────┬──────────────────────────────────────────────────────────────┤
│      │  # Welcome to Example Page                                   │
│  📑  │                                                                │
│ Book │  This is a normal text line in gemtext format.               │
│ mark │                                                                │
│  s   │  ## Links                                                     │
│      │  [1] => About This Site                                       │
│ [▼]  │  [2] => Latest News                                           │
│ Home │  [3] => Documentation                                         │
│ Work │                                                                │
│ Blog │  ## Code Example                                              │
│      │  ```go                                                        │
│ [▼]  │  func main() {                                                │
│ Hist │      fmt.Println("Hello, Gemini!")                           │
│ ory  │  }                                                             │
│      │  ```                                                           │
│ Today│                                                                │
│ Week │  * Feature one                                                │
│ Month│  * Feature two                                                │
│      │  > "Simplicity is the ultimate sophistication"               │
│      │                                                                │
├──────┴──────────────────────────────────────────────────────────────┤
│ ✓ Loaded | Link 1/12 | Scroll 45% | UTF-8        [Ctrl+H: Help]    │
└─────────────────────────────────────────────────────────────────────┘
```

### Components

#### 1. Tab Bar (Top)
- Display all open tabs with truncated titles
- Active tab highlighted
- Close button per tab (hover/focus)
- New tab button `[+]`
- Help/settings indicator

#### 2. Address Bar
- Current URL display
- URL input mode (Ctrl+L)
- Status indicators:
  - 🔄 Loading
  - ✓ Success
  - ⚠️ Warning
  - ✗ Error
- Quick actions:
  - Bookmark toggle
  - Settings menu

#### 3. Sidebar (Collapsible)
- **Bookmarks Section**
  - Tree view with folders
  - Drag-and-drop organization (if supported)
  - Quick search/filter
- **History Section**
  - Today
  - This week
  - This month
  - All history
- Toggle: `Ctrl+B`
- Resizable width

#### 4. Content Area (Main)
- Rendered gemtext content
- Visual styling:
  - H1: Large, bold, colored
  - H2: Medium, bold
  - H3: Normal, bold
  - Links: Highlighted with numbers `[1]`, `[2]`
  - Lists: Bullet points
  - Quotes: Left border, italic
  - Code: Monospace, syntax highlighting
- Scrollable
- Link selection with Tab/Shift+Tab

#### 5. Status Bar (Bottom)
- Left side:
  - Connection status
  - Current link (when selected)
  - Scroll position
- Right side:
  - Character encoding
  - Quick help reminder

### Modal Dialogs

**Input Modal** (Status 10/11):
```
┌────────────────────────────────────┐
│         Input Required             │
├────────────────────────────────────┤
│ Enter search query:                │
│ [________________________]         │
│                                    │
│        [OK]      [Cancel]          │
└────────────────────────────────────┘
```

**Certificate Warning Modal**:
```
┌────────────────────────────────────┐
│      Certificate Changed!          │
├────────────────────────────────────┤
│ The certificate for                │
│ example.com has changed.           │
│                                    │
│ Previous: ABC123...                │
│ Current:  DEF456...                │
│                                    │
│ This could be a security risk.     │
│                                    │
│  [Accept Once]  [Accept Forever]   │
│           [Reject]                 │
└────────────────────────────────────┘
```

### Color Scheme (Default Theme)

| Element | Foreground | Background | Style |
|---------|-----------|------------|-------|
| H1 | Bright Cyan | Default | Bold |
| H2 | Cyan | Default | Bold |
| H3 | Blue | Default | Bold |
| Link | Green | Default | Underline |
| Selected Link | Black | Green | Bold |
| List Bullet | Yellow | Default | - |
| Quote | Magenta | Default | Italic |
| Code Block | White | Dark Gray | - |
| Active Tab | Black | White | Bold |
| Inactive Tab | Gray | Default | - |
| Status Bar | White | Blue | - |
| Error | Red | Default | Bold |

---

## Core Features

### Phase 1: Essential Features (MVP)

#### 1. URL Navigation
- **Address bar** for URL entry
- **Link following** with keyboard/mouse
- **Back/forward** navigation
- **Home page** configuration
- **Redirect handling** (30, 31)

#### 2. Gemtext Rendering
- Parse and render all gemtext elements
- Proper text wrapping
- Link numbering for quick access
- Syntax highlighting for code blocks
- ANSI color code support

#### 3. TLS & Certificate Management
- TLS 1.2+ connection
- **TOFU implementation**:
  - Store certificates on first visit
  - Verify on subsequent visits
  - Warn on certificate changes
- Certificate acceptance/rejection UI
- Certificate storage (JSON)

#### 4. Bookmarks
- Add bookmark (Ctrl+D)
- Remove bookmark
- Organize in folders
- Bookmark sidebar
- Quick access
- Persistence (JSON file)

#### 5. History
- Track all visited pages
- Per-tab history (back/forward)
- Global history view
- History sidebar (today/week/month)
- Clear history option
- Persistence (SQLite)

### Phase 2: Advanced Features

#### 6. Tabbed Browsing
- Multiple tabs
- New tab (Ctrl+T)
- Close tab (Ctrl+W)
- Switch tabs:
  - Ctrl+Tab / Ctrl+Shift+Tab
  - Ctrl+1-9 (jump to tab)
- Session restoration
- Per-tab navigation history

#### 7. Downloads
- Detect non-gemtext MIME types
- Download prompt with destination
- Progress indicator
- Download history
- Configurable download directory

#### 8. Search
- **Find in page** (Ctrl+F)
  - Highlight matches
  - Navigate between matches (n/N)
- **Geminispace search**
  - Built-in search engine (geminispace.info)
  - Custom search engine configuration

#### 9. Input Handling
- Modal dialog for status 10 (input)
- Password field for status 11 (sensitive input)
- Form submission handling
- Input history

#### 10. Client Certificates
- Generate client certificates
- Associate certificates with domains
- Certificate management UI
- Certificate selection prompt

### Phase 3: Polish & Extra Features

#### 11. Subscriptions/Feeds
- Subscribe to pages
- Check for updates
- Notification on changes
- Atom/RSS feed support

#### 12. Theming
- Multiple built-in themes
- Theme selector UI
- Custom theme support
- Theme import/export (JSON)

#### 13. Caching
- Cache page content
- Configurable cache size
- Cache expiration
- Clear cache option

#### 14. Configuration
- Config file (TOML)
- Settings:
  - Home page
  - Default search engine
  - Download directory
  - Theme
  - Keybindings
  - Cache settings
  - Proxy settings

#### 15. Proxy Support
- HTTP-to-Gemini proxy
- Configurable proxy servers
- Per-domain proxy rules

#### 16. Streaming
- Stream large files instead of download-first
- Progress indicator
- Audio streaming support

---

## Technical Specifications

### Keyboard Shortcuts

```
Navigation:
  Ctrl+L              Focus address bar
  Ctrl+Enter          Navigate to URL in address bar
  Alt+Left            Back
  Alt+Right           Forward
  Ctrl+R              Reload current page
  Ctrl+H              Go to home page

Tabs:
  Ctrl+T              New tab
  Ctrl+W              Close current tab
  Ctrl+Tab            Next tab
  Ctrl+Shift+Tab      Previous tab
  Ctrl+1-9            Jump to tab 1-9

Links & Navigation:
  Tab                 Next link
  Shift+Tab           Previous link
  Enter               Follow selected link
  1-9, 0              Follow link by number (if 1-10 links on page)
  j/k                 Scroll down/up (vim-style)
  Space               Page down
  Shift+Space         Page up
  g                   Go to top of page
  G                   Go to bottom of page

Bookmarks & History:
  Ctrl+D              Bookmark current page
  Ctrl+B              Toggle bookmarks sidebar
  Ctrl+Shift+B        Manage bookmarks
  Ctrl+Shift+H        Show history

Content:
  Ctrl+F              Find in page
  n                   Next search result
  N                   Previous search result
  Ctrl+S              Save page to disk
  Ctrl+U              View source (raw gemtext)

Other:
  ?                   Show help screen
  Ctrl+,              Open settings
  Ctrl+Q              Quit application
  Esc                 Cancel/close modal or input
```

### Configuration File

**Location:** `~/.config/gemini-client/config.toml`

```toml
[general]
home_page = "gemini://geminiprotocol.net"
download_dir = "~/Downloads"
cache_enabled = true
cache_size_mb = 100
cache_ttl_hours = 24

[appearance]
theme = "default"
show_line_numbers = false
wrap_text = true
max_width = 100

[network]
timeout_seconds = 30
max_redirects = 5
user_agent = "gemini-client/1.0"

[network.proxy]
enabled = false
url = ""

[search]
default_engine = "gemini://geminispace.info/search"

[keybindings]
quit = "Ctrl+Q"
new_tab = "Ctrl+T"
close_tab = "Ctrl+W"
# ... (customizable)

[certificates]
tofu_enabled = true
warn_on_change = true
client_cert_dir = "~/.config/gemini-client/certificates"
```

### Data Storage

#### Directory Structure

```
~/.config/gemini-client/
├── config.toml              # Main configuration
├── bookmarks.json           # Bookmarks
├── history.db               # SQLite history database
├── cache/                   # Page cache
│   ├── <hash1>.gmi
│   ├── <hash2>.gmi
│   └── ...
├── certificates/            # TOFU certificate storage
│   ├── known_hosts.json     # Trusted certificates
│   └── client/              # Client certificates
│       ├── example.com.crt
│       └── example.com.key
└── themes/                  # User themes
    └── mytheme.json
```

#### Bookmarks Schema (JSON)

```json
{
  "version": "1.0",
  "bookmarks": [
    {
      "id": "uuid-1",
      "title": "Project Gemini",
      "url": "gemini://geminiprotocol.net",
      "tags": ["official", "docs"],
      "created": "2025-11-18T12:00:00Z",
      "folder": "Tech"
    },
    {
      "id": "uuid-2",
      "title": "My Capsule",
      "url": "gemini://example.com",
      "tags": ["personal"],
      "created": "2025-11-18T13:00:00Z",
      "folder": ""
    }
  ],
  "folders": [
    {
      "name": "Tech",
      "parent": ""
    },
    {
      "name": "Blogs",
      "parent": ""
    }
  ]
}
```

#### History Schema (SQLite)

```sql
CREATE TABLE history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT NOT NULL,
    title TEXT,
    visited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    visit_count INTEGER DEFAULT 1
);

CREATE INDEX idx_url ON history(url);
CREATE INDEX idx_visited_at ON history(visited_at);
```

#### Known Hosts Schema (JSON)

```json
{
  "version": "1.0",
  "hosts": {
    "example.com": {
      "fingerprint": "SHA256:abc123...",
      "first_seen": "2025-11-18T12:00:00Z",
      "last_seen": "2025-11-18T14:30:00Z",
      "trust": "permanent"
    }
  }
}
```

### Error Handling

#### Network Errors
- Connection timeout → Retry with exponential backoff
- TLS handshake failure → Show certificate error
- DNS failure → Show user-friendly error

#### Protocol Errors
- Invalid status code → Show raw response
- Malformed response → Show error, log details
- Redirect loop → Detect after N redirects, abort

#### User Errors
- Invalid URL → Show error, keep in address bar for correction
- Bookmark already exists → Show warning, allow update
- No space for cache → Clear old entries, warn user

---

## Implementation Plan

### Phase 1: Core Protocol (Week 1-2)

**Goals:**
- Working Gemini protocol client
- TLS connection with TOFU
- Basic gemtext parsing

**Tasks:**
1. Project setup (Go module, directory structure)
2. Implement `protocol` package:
   - TLS connection handler
   - Request builder
   - Response parser
   - Status code constants
3. Implement `protocol/tofu` package:
   - Certificate storage
   - Certificate verification
   - Change detection
4. Implement `parser` package:
   - Gemtext lexer/parser
   - AST representation
   - Basic renderer
5. Unit tests for all packages
6. CLI tool for testing (simple, non-TUI)

**Deliverable:** CLI tool that can fetch and parse Gemini pages

### Phase 2: Basic TUI (Week 2-3)

**Goals:**
- Working TUI with Bubble Tea
- Basic navigation
- Content viewing

**Tasks:**
1. Set up Bubble Tea application structure
2. Implement `ui/app` package (main model)
3. Implement `ui/browser` component:
   - Content rendering
   - Scrolling
   - Link selection
4. Implement `ui/addressbar` component
5. Implement `ui/statusbar` component
6. Wire up navigation (load URL, follow links)
7. Basic styling with Lip Gloss

**Deliverable:** TUI app that can navigate Geminispace

### Phase 3: Essential Features (Week 3-4)

**Goals:**
- Bookmarks
- History
- Input handling

**Tasks:**
1. Implement `storage/bookmarks` package
2. Implement `storage/history` package (SQLite)
3. Implement `ui/sidebar` component
4. Implement `ui/modal` component (for input)
5. Wire up Ctrl+D (bookmark), Ctrl+B (sidebar)
6. Handle status codes 10/11 (input prompts)
7. Persistence (save/load bookmarks and history)

**Deliverable:** Usable browser with bookmarks and history

### Phase 4: Advanced Features (Week 4-5)

**Goals:**
- Tabs
- Downloads
- Search

**Tasks:**
1. Implement `ui/tabs` component
2. Refactor state to support multiple tabs
3. Implement download handling (non-gemtext MIME types)
4. Implement `ui/modal` for find-in-page
5. Add search engine support
6. Client certificate management
7. Configuration file loading

**Deliverable:** Feature-complete browser

### Phase 5: Polish & Extras (Week 5-6)

**Goals:**
- Theming
- Caching
- Documentation

**Tasks:**
1. Implement `theme` package
2. Create multiple built-in themes
3. Implement `storage/cache` package
4. Add feed subscription support
5. Performance optimization
6. Comprehensive testing
7. User documentation (README, help screen)
8. Developer documentation (code comments, architecture docs)

**Deliverable:** Production-ready v1.0.0

---

## Testing Strategy

### Unit Tests
- All protocol handling
- Gemtext parser (edge cases)
- TOFU certificate logic
- Bookmark/history CRUD
- Configuration loading

### Integration Tests
- End-to-end navigation flows
- TLS handshake with test server
- Cache behavior
- Session persistence

### Manual Testing
- Real Geminispace browsing
- Various capsules (different gemtext styles)
- Performance with large pages
- Error scenarios

---

## Future Enhancements (v2.0+)

- Mouse support (click links, scroll)
- Split panes (view multiple pages)
- Vim-style marks for quick navigation
- Macro recording/playback
- Plugin system
- Gemini-to-HTTP proxy mode
- Gopher protocol support
- Finger protocol support
- RSS reader integration
- Markdown export from gemtext

---

## References

- [Gemini Protocol Specification](https://geminiprotocol.net/docs/protocol-specification.gmi)
- [Gemtext Specification](https://geminiprotocol.net/docs/gemtext-specification.gmi)
- [Awesome Gemini](https://github.com/kr1sp1n/awesome-gemini)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Amfora Client](https://github.com/makew0rld/amfora) (inspiration)

---

**Document Version:** 1.0.0
**Status:** Draft → **Approved** → Implementation
**Next Review:** After Phase 2 completion
