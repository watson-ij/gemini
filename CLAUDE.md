# CLAUDE.md - AI Assistant Guide

This document helps AI assistants (like Claude) navigate and contribute to the Gemini Browser project effectively.

## Project Overview

**Gemini Browser** is a modern, terminal-based browser for the Gemini protocol built with Go and Bubble Tea. It provides a feature-rich TUI (Text User Interface) for browsing Geminispace with keyboard-driven navigation, beautiful rendering, and essential features like bookmarks and history.

- **Language**: Go 1.24+
- **UI Framework**: Bubble Tea (Elm-architecture TUI framework)
- **Protocol**: Gemini (gemini://) - a lightweight, privacy-focused internet protocol
- **Architecture**: Clean modular design with separation of protocol, parsing, UI, and storage

## Quick Navigation

### Essential Files
- `main.go` - Application entry point (starts Bubble Tea program)
- `README.md` - User-facing documentation with features and usage
- `DESIGN.md` - Comprehensive technical design document
- `config.example.toml` - Example configuration file

### Key Directories
```
internal/
├── protocol/       # Gemini protocol implementation (TLS, requests, TOFU)
├── parser/         # Gemtext parsing and rendering
├── ui/            # Bubble Tea UI components and application model
└── config/        # Configuration management
```

## Architecture

### The Elm Architecture (Bubble Tea)
The application follows the Model-Update-View pattern:

1. **Model** (`internal/ui/app.go`) - Application state (current page, history, UI state)
2. **Update** (`internal/ui/app.go`) - Message handler (keyboard events, network responses)
3. **View** (`internal/ui/app.go`) - Rendering logic (converts model to terminal output)

### Key Components

#### Protocol Layer (`internal/protocol/`)
- `client.go` - Gemini protocol client with TLS support
- `response.go` - Response parsing and status code handling
- `status.go` - Gemini status code definitions (1x-6x ranges)
- `tofu.go` - Trust On First Use certificate verification

**Important**: Gemini uses TLS by default (port 1965) with TOFU certificate verification. Self-signed certs are common and acceptable.

#### Parser Layer (`internal/parser/`)
- `gemtext.go` - Gemtext lexer/parser (converts gemtext to AST)
- `types.go` - AST node types (headings, links, lists, quotes, preformatted)
- `renderer.go` - Renders AST to styled terminal output with syntax highlighting

**Gemtext Line Types**:
- Headings: `#`, `##`, `###`
- Links: `=> URL [optional text]`
- Lists: `* item`
- Quotes: `> text`
- Preformatted: ` ``` ` toggles
- Text: everything else

#### UI Layer (`internal/ui/`)
- `app.go` - Main Bubble Tea model and update/view logic
- `keys.go` - Keyboard binding definitions

**Key UI Features**:
- Address bar (Ctrl+L to focus)
- Content viewport with scrolling
- Link selection and navigation
- Status bar
- Help screen (?)

#### Configuration (`internal/config/`)
- `config.go` - TOML configuration loading
- Config location: `~/.config/gemini-client/config.toml`
- Supports display settings (wrap_width, line numbers)

## Common Development Tasks

### Adding a New Keyboard Shortcut
1. Define the key binding in `internal/ui/keys.go`
2. Handle the key press in the `Update()` function in `internal/ui/app.go`
3. Document in README.md keyboard shortcuts section

### Modifying Gemtext Rendering
1. Parser logic: `internal/parser/gemtext.go`
2. AST types: `internal/parser/types.go`
3. Styling/rendering: `internal/parser/renderer.go`
4. Use lipgloss for terminal styling

### Adding Protocol Features
1. Protocol implementation: `internal/protocol/client.go`
2. Response handling: `internal/protocol/response.go`
3. Status codes: `internal/protocol/status.go`

### Adding UI Components
1. Create component in `internal/ui/`
2. Add to main model in `app.go`
3. Update `Update()` to handle component messages
4. Update `View()` to render component

## Testing

### Running Tests
```bash
go test ./...                    # All tests
go test ./internal/parser/...    # Parser tests only
go test -v ./...                 # Verbose output
```

### Test Files
- `internal/parser/gemtext_test.go` - Parser tests
- `internal/parser/renderer_test.go` - Renderer tests

### Manual Testing
```bash
go build -o gemini-browser .
./gemini-browser gemini://geminiprotocol.net
```

**Test Sites**:
- `gemini://geminiprotocol.net` - Official Gemini site
- `gemini://warmedal.se/~antenna/` - Gemini feed aggregator
- `gemini://geminispace.info` - Gemini search

## Code Style & Conventions

### General Go Practices
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names
- Add comments for exported functions
- Keep functions focused and small

### Project-Specific Conventions
- UI state goes in the main model (`internal/ui/app.go`)
- Protocol logic stays separate from UI logic
- Use Bubble Tea commands for async operations (network requests)
- Style with lipgloss, not raw ANSI codes

### Error Handling
- Return errors, don't panic (except in main)
- Provide user-friendly error messages in UI
- Log technical details for debugging

## Recent Changes & Active Features

### Recently Added (check git log for latest)
- Page navigation keys ('n' and 'p') - See commit a324a6b
- Link scroll behavior fixes - See commit 5a1859d
- View initialization fixes

### TODO/Coming Soon (from README.md)
- Tabbed browsing
- Bookmarks with folder organization
- History tracking and search
- Find in page
- Downloads
- Client certificate management
- Multiple themes

See DESIGN.md Phase 2-5 for detailed implementation plans.

## Git Workflow

### Branch Naming
- Feature branches: `claude/add-<feature>-<session-id>`
- Current branch: `claude/add-clau-feature-0124gJGasYVTih4AfoFuVkCg`

### Commit Guidelines
1. Use clear, descriptive commit messages
2. Focus on "why" rather than "what"
3. Reference issue numbers if applicable
4. Follow existing commit message style (see `git log`)

### Pushing Changes
```bash
git push -u origin claude/add-clau-feature-0124gJGasYVTih4AfoFuVkCg
```

**Important**: Branch must start with 'claude/' and match session ID

## Dependencies

### Core Libraries
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - Reusable UI components
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/BurntSushi/toml` - Configuration parsing

### Adding Dependencies
```bash
go get <package>
go mod tidy
```

## Building & Running

### Development Build
```bash
go build -o gemini-browser .
./gemini-browser [optional-gemini-url]
```

### With Race Detector
```bash
go build -race -o gemini-browser .
```

### Clean Build
```bash
go clean
go build -o gemini-browser .
```

## Debugging Tips

### UI Issues
- Check `View()` function in `internal/ui/app.go`
- Verify lipgloss styles in renderer
- Test with different terminal sizes

### Protocol Issues
- Enable verbose logging in protocol client
- Check TLS handshake with test server
- Verify response parsing in `internal/protocol/response.go`

### Parser Issues
- Add test cases in `internal/parser/gemtext_test.go`
- Check AST construction in `internal/parser/gemtext.go`
- Verify rendering output in `internal/parser/renderer.go`

## Understanding the Gemini Protocol

### Key Concepts
- **TOFU (Trust On First Use)**: Like SSH, accept cert on first visit, verify on subsequent visits
- **Status Codes**: 1x=Input, 2x=Success, 3x=Redirect, 4x=Temp Failure, 5x=Perm Failure, 6x=Client Cert
- **Gemtext**: Line-oriented markup (not like Markdown - no inline formatting)

### Common Status Codes
- `10` - Input required (show modal)
- `11` - Sensitive input (password field)
- `20` - Success (render content)
- `30/31` - Redirect (follow automatically)
- `51` - Not found
- `60` - Client certificate required

## Resources

### Documentation
- [Gemini Protocol Spec](https://geminiprotocol.net/docs/protocol-specification.gmi)
- [Gemtext Spec](https://geminiprotocol.net/docs/gemtext-specification.gmi)
- [Bubble Tea Docs](https://github.com/charmbracelet/bubbletea)

### Similar Projects (for inspiration)
- [Amfora](https://github.com/makew0rld/amfora) - Feature-rich Gemini client
- [Lagrange](https://gmi.skyjake.fi/lagrange/) - GUI Gemini client

## When Making Changes

### Before You Start
1. Read the relevant section in DESIGN.md
2. Check existing tests
3. Understand the affected components
4. Consider backward compatibility

### Making Changes
1. Keep changes focused and atomic
2. Update tests if behavior changes
3. Update documentation (README, DESIGN, comments)
4. Test manually with real Gemini sites
5. Ensure code follows project conventions

### Security Considerations
- Validate URLs before making requests
- Handle TLS certificate changes carefully (security implications)
- Sanitize user input (especially for status 10/11 prompts)
- Don't store sensitive data in plaintext

## Quick Reference

### File Lookup
| Need to... | Look in... |
|------------|-----------|
| Add keyboard shortcut | `internal/ui/keys.go`, `internal/ui/app.go` |
| Modify gemtext parsing | `internal/parser/gemtext.go` |
| Change text styling | `internal/parser/renderer.go` |
| Add protocol feature | `internal/protocol/client.go` |
| Modify UI layout | `internal/ui/app.go` (View function) |
| Add configuration option | `internal/config/config.go`, `config.example.toml` |
| Handle new status code | `internal/protocol/status.go`, `internal/ui/app.go` |

### Command Cheat Sheet
```bash
go build .                   # Build
go test ./...               # Test
go run . <url>              # Run with URL
go mod tidy                 # Clean dependencies
git status                  # Check changes
git log --oneline -10       # Recent commits
```

## Getting Help

If you're stuck:
1. Check DESIGN.md for architectural decisions
2. Review similar existing code
3. Look at test cases for examples
4. Check the Gemini protocol spec
5. Review Bubble Tea documentation

---

**Last Updated**: 2025-11-18
**Project Version**: v0.1.0 (MVP)
**Maintainer**: watson-ij
