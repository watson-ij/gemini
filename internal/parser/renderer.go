package parser

import (
	"fmt"
	"strings"
)

// RenderOptions contains options for rendering a gemtext document
type RenderOptions struct {
	// Width is the maximum width for text wrapping (0 = no wrapping)
	Width int

	// ShowLineNumbers shows line numbers
	ShowLineNumbers bool

	// NumberLinks adds numeric labels to links [1], [2], etc.
	NumberLinks bool

	// HighlightedLink is the index of the currently highlighted link (-1 = none)
	HighlightedLink int

	// ColorScheme contains the color/style codes for different elements
	// These would be lipgloss styles in the real implementation
	ColorScheme *ColorScheme
}

// ColorScheme contains ANSI color codes or lipgloss styles for different elements
type ColorScheme struct {
	// For now, we'll use simple string markers
	// In the real TUI, these would be lipgloss.Style objects
	Heading1   string
	Heading2   string
	Heading3   string
	Link       string
	LinkActive string
	ListBullet string
	Quote      string
	Preformat  string
	Text       string
	Reset      string
}

// DefaultColorScheme returns a default color scheme
func DefaultColorScheme() *ColorScheme {
	return &ColorScheme{
		Heading1:   "\033[1;36m", // Bold Cyan
		Heading2:   "\033[36m",   // Cyan
		Heading3:   "\033[34m",   // Blue
		Link:       "\033[32m",   // Green
		LinkActive: "\033[1;42m", // Bold on Green background
		ListBullet: "\033[33m",   // Yellow
		Quote:      "\033[35m",   // Magenta
		Preformat:  "\033[37m",   // White
		Text:       "",           // Default
		Reset:      "\033[0m",    // Reset
	}
}

// Renderer renders gemtext documents to styled text
type Renderer struct {
	opts *RenderOptions
}

// NewRenderer creates a new renderer with the given options
func NewRenderer(opts *RenderOptions) *Renderer {
	if opts == nil {
		opts = &RenderOptions{
			Width:           80,
			ShowLineNumbers: false,
			NumberLinks:     true,
			HighlightedLink: -1,
			ColorScheme:     DefaultColorScheme(),
		}
	}

	if opts.ColorScheme == nil {
		opts.ColorScheme = DefaultColorScheme()
	}

	return &Renderer{opts: opts}
}

// Render renders a document to a string
func (r *Renderer) Render(doc *Document) string {
	var b strings.Builder
	linkIndex := 0

	for i, line := range doc.Lines {
		rendered := r.renderLine(line, i, &linkIndex)
		b.WriteString(rendered)
		b.WriteString("\n")
	}

	return b.String()
}

// renderLine renders a single line
func (r *Renderer) renderLine(line *Line, lineNum int, linkIndex *int) string {
	cs := r.opts.ColorScheme

	// Line number prefix (if enabled)
	prefix := ""
	if r.opts.ShowLineNumbers {
		prefix = fmt.Sprintf("%4d  ", lineNum+1)
	}

	switch line.Type {
	case LineTypeHeading1:
		return prefix + cs.Heading1 + "# " + line.Text + cs.Reset

	case LineTypeHeading2:
		return prefix + cs.Heading2 + "## " + line.Text + cs.Reset

	case LineTypeHeading3:
		return prefix + cs.Heading3 + "### " + line.Text + cs.Reset

	case LineTypeLink:
		linkNum := *linkIndex
		*linkIndex++

		// Determine if this link is highlighted
		style := cs.Link
		if r.opts.HighlightedLink == linkNum {
			style = cs.LinkActive
		}

		// Add link number if enabled
		linkLabel := ""
		if r.opts.NumberLinks {
			linkLabel = fmt.Sprintf("[%d] ", linkNum+1)
		}

		return prefix + style + linkLabel + line.Link.Display + cs.Reset

	case LineTypeListItem:
		return prefix + cs.ListBullet + "• " + cs.Reset + line.Text

	case LineTypeQuote:
		return prefix + cs.Quote + "│ " + line.Text + cs.Reset

	case LineTypePreformatted:
		return prefix + cs.Preformat + line.Text + cs.Reset

	case LineTypePreformatToggle:
		// Don't render the toggle lines themselves
		return ""

	case LineTypeText:
		if line.Text == "" {
			return "" // Empty line
		}
		return prefix + cs.Text + line.Text + cs.Reset

	default:
		return prefix + line.Text
	}
}

// RenderToPlainText renders a document to plain text (no colors)
func (r *Renderer) RenderToPlainText(doc *Document) string {
	// Temporarily remove colors
	originalScheme := r.opts.ColorScheme
	r.opts.ColorScheme = &ColorScheme{} // All empty strings

	result := r.Render(doc)

	r.opts.ColorScheme = originalScheme
	return result
}

// GetLinkAtLine returns the link index at a given line number, or -1 if none
func GetLinkAtLine(doc *Document, lineNum int) int {
	if lineNum < 0 || lineNum >= len(doc.Lines) {
		return -1
	}

	line := doc.Lines[lineNum]
	if line.Type != LineTypeLink {
		return -1
	}

	// Find the index of this link
	for i, link := range doc.Links {
		if link == line {
			return i
		}
	}

	return -1
}

// GetLineForLink returns the line number for a given link index
func GetLineForLink(doc *Document, linkIndex int) int {
	if linkIndex < 0 || linkIndex >= len(doc.Links) {
		return -1
	}

	targetLink := doc.Links[linkIndex]

	for i, line := range doc.Lines {
		if line == targetLink {
			return i
		}
	}

	return -1
}
