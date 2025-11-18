# 📡 Gemini Browser

A modern, feature-rich terminal-based browser for the [Gemini protocol](https://geminiprotocol.net), built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/go-1.21+-blue)
![License](https://img.shields.io/badge/license-MIT-blue)

## Features

### Core Features ✨

- **Full Gemini Protocol Support**
  - TLS 1.2+ with TOFU (Trust On First Use) certificate verification
  - All status codes (input, success, redirect, errors, client certificates)
  - Automatic redirect following
  - Proper gemtext parsing and rendering

- **Beautiful TUI**
  - Syntax-highlighted gemtext rendering
  - Smart text wrapping with configurable width
  - Link numbering and navigation
  - Vim-inspired keyboard shortcuts
  - Responsive design with lipgloss styling
  - Color-coded headings, links, quotes, and code blocks

- **Navigation**
  - Address bar with URL history
  - Back/forward navigation
  - Link selection and following
  - Keyboard-driven browsing

- **Developer-Friendly**
  - Clean, modular architecture
  - Well-documented code
  - Easy to extend and customize

### Coming Soon 🚧

- [ ] Tabbed browsing
- [ ] Bookmarks with folder organization
- [ ] History tracking and search
- [ ] Find in page
- [ ] Downloads
- [ ] Client certificate management
- [ ] Multiple themes
- [ ] Page caching
- [ ] Subscriptions/feeds

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/watson-ij/gemini.git
cd gemini

# Build the browser
go build -o gemini-browser .

# Run it
./gemini-browser gemini://geminiprotocol.net
```

### Using Go Install

```bash
go install github.com/watson-ij/gemini@latest
```

## Usage

### Basic Usage

```bash
# Start with default homepage
./gemini-browser

# Start with a specific URL
./gemini-browser gemini://geminiprotocol.net

# Any gemini:// URL works
./gemini-browser gemini://warmedal.se/~antenna/
```

### Keyboard Shortcuts

#### Navigation
- `↑/k`, `↓/j` - Scroll up/down
- `PgUp`/`PgDn` - Page up/down
- `Space` - Page down
- `Shift+Space` - Page up
- `g` - Go to top
- `G` - Go to bottom

#### Links
- `Tab` - Select next link
- `Shift+Tab` - Select previous link
- `Enter` - Follow selected link

#### URL Navigation
- `Ctrl+L` - Focus address bar
- `Enter` - Navigate to URL (when in address bar)
- `Esc` - Cancel address bar editing
- `Ctrl+R` - Reload current page
- `Alt+←` - Go back in history
- `Alt+→` - Go forward in history

#### Other
- `?` - Show help screen
- `Ctrl+Q` - Quit application

## Architecture

The project follows a clean, modular architecture:

```
gemini/
├── internal/
│   ├── protocol/      # Gemini protocol implementation (TLS, TOFU, status codes)
│   ├── parser/        # Gemtext parser and renderer
│   ├── ui/            # Bubble Tea TUI components
│   ├── storage/       # Bookmarks, history, cache (TODO)
│   └── theme/         # Theming system (TODO)
├── cmd/gemini/        # CLI entry point
├── DESIGN.md          # Comprehensive design document
└── README.md          # This file
```

### Key Components

- **Protocol Package** (`internal/protocol/`)
  - Client with TLS support
  - Request/response handling
  - TOFU certificate verification
  - Status code definitions

- **Parser Package** (`internal/parser/`)
  - Gemtext lexer and parser
  - AST representation
  - Renderer with syntax highlighting

- **UI Package** (`internal/ui/`)
  - Bubble Tea application model
  - Viewport for content display
  - Address bar and status bar
  - Keyboard bindings

## Configuration

The application can be configured via a TOML configuration file located at:
`~/.config/gemini-client/config.toml`

If the configuration file doesn't exist, sensible defaults are used.

### Configuration Options

```toml
[display]
# Maximum width for text wrapping (in characters)
# Set to 0 to use the full terminal width
# Default: 100
wrap_width = 100

# Show line numbers in the margin
# Default: false
show_line_numbers = false
```

### Creating a Configuration File

A sample configuration file is provided in `config.example.toml`. To use it:

```bash
# Create config directory if it doesn't exist
mkdir -p ~/.config/gemini-client

# Copy the example config
cp config.example.toml ~/.config/gemini-client/config.toml

# Edit to your preferences
nano ~/.config/gemini-client/config.toml
```

You can also view the configuration path and settings from within the browser by pressing `?` to open the help screen.

## Development

### Prerequisites

- Go 1.21 or higher
- Terminal with ANSI color support

### Building

```bash
# Build the application
go build -o gemini-browser .

# Run tests
go test ./...

# Build with race detector
go build -race -o gemini-browser .
```

### Project Structure

See [DESIGN.md](DESIGN.md) for comprehensive documentation on:
- Protocol specification
- Architecture decisions
- Implementation details
- Feature roadmap
- Data storage schemas

### Contributing

Contributions are welcome! Please:

1. Check [DESIGN.md](DESIGN.md) for the project roadmap
2. Open an issue to discuss major changes
3. Follow the existing code style
4. Add tests for new features
5. Update documentation as needed

## Gemini Protocol

The Gemini protocol is a modern, lightweight internet protocol that:

- Emphasizes simplicity and privacy
- Uses TLS by default (port 1965)
- Has its own lightweight markup language (gemtext)
- Is heavier than Gopher but lighter than HTTP/HTML
- Has no cookies, no tracking, no JavaScript

Learn more at [geminiprotocol.net](https://geminiprotocol.net)

## Resources

- **Official Gemini Sites**
  - [Gemini Protocol](gemini://geminiprotocol.net) - Official specs
  - [Antenna](gemini://warmedal.se/~antenna/) - Gemini feed aggregator
  - [Gemini Search](gemini://geminispace.info) - Search engine

- **Related Projects**
  - [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
  - [Amfora](https://github.com/makew0rld/amfora) - Another great Gemini client (inspiration)
  - [Awesome Gemini](https://github.com/kr1sp1n/awesome-gemini) - Curated list of Gemini resources

## Roadmap

### v0.1.0 (Current - MVP)
- [x] Basic Gemini protocol client
- [x] TLS with TOFU
- [x] Gemtext parser and renderer
- [x] TUI with Bubble Tea
- [x] Link navigation
- [x] Address bar
- [x] Back/forward history
- [x] Smart text wrapping
- [x] Configuration file support

### v0.2.0 (Next Release)
- [ ] Tabbed browsing
- [ ] Bookmarks
- [ ] Persistent history
- [ ] Find in page
- [ ] Downloads

### v0.3.0 (Future)
- [ ] Client certificates
- [ ] Theming system
- [ ] Page caching
- [ ] Subscriptions

### v1.0.0 (Stable)
- [ ] All planned features complete
- [ ] Comprehensive testing
- [ ] Performance optimization
- [ ] Full documentation

See [DESIGN.md](DESIGN.md) for detailed implementation plans.

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Solderpunk and the Gemini community for creating the protocol
- [Charm Bracelet](https://charm.sh) for the amazing TUI libraries
- All Gemini client developers for inspiration

## Support

- 🐛 Report bugs: [GitHub Issues](https://github.com/watson-ij/gemini/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/watson-ij/gemini/discussions)
- 📖 Documentation: [DESIGN.md](DESIGN.md)

---

**Happy browsing in Geminispace! 🚀**
