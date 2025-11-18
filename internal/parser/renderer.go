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

// wrapText wraps text to a maximum width, breaking at word boundaries
func wrapText(text string, width int, indent string) []string {
	if width <= 0 {
		return []string{text}
	}

	// Calculate effective width (accounting for indent on wrapped lines)
	effectiveWidth := width - len(indent)
	if effectiveWidth <= 10 {
		// If indent is too large, don't wrap
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder

	for i, word := range words {
		// Check if adding this word would exceed width
		if currentLine.Len() > 0 {
			testLen := currentLine.Len() + 1 + len(word) // +1 for space
			maxLen := width
			if len(lines) > 0 {
				maxLen = effectiveWidth // Use effective width for continuation lines
			}

			if testLen > maxLen {
				// Line would be too long, start a new line
				lines = append(lines, currentLine.String())
				currentLine.Reset()
				if len(lines) > 0 {
					currentLine.WriteString(indent)
				}
				currentLine.WriteString(word)
			} else {
				// Add word to current line
				currentLine.WriteString(" ")
				currentLine.WriteString(word)
			}
		} else {
			// First word on the line
			if len(lines) > 0 && i > 0 {
				currentLine.WriteString(indent)
			}
			currentLine.WriteString(word)
		}
	}

	// Add the last line
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// Render renders a document to a string
func (r *Renderer) Render(doc *Document) string {
	var b strings.Builder
	linkIndex := 0

	for i, line := range doc.Lines {
		rendered := r.renderLine(line, i, &linkIndex)
		b.WriteString(rendered)
		if i < len(doc.Lines)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// renderLine renders a single line, with optional text wrapping
func (r *Renderer) renderLine(line *Line, lineNum int, linkIndex *int) string {
	cs := r.opts.ColorScheme

	// Line number prefix (if enabled)
	prefix := ""
	if r.opts.ShowLineNumbers {
		prefix = fmt.Sprintf("%4d  ", lineNum+1)
	}

	switch line.Type {
	case LineTypeHeading1:
		return r.renderWrappedLine(prefix+cs.Heading1+"# "+cs.Reset, line.Text, cs.Heading1, cs.Reset, "  ")

	case LineTypeHeading2:
		return r.renderWrappedLine(prefix+cs.Heading2+"## "+cs.Reset, line.Text, cs.Heading2, cs.Reset, "   ")

	case LineTypeHeading3:
		return r.renderWrappedLine(prefix+cs.Heading3+"### "+cs.Reset, line.Text, cs.Heading3, cs.Reset, "    ")

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

		linkPrefix := prefix + style + linkLabel
		// Calculate indent width (without ANSI codes)
		indentWidth := len(prefix) + len(linkLabel)
		indent := strings.Repeat(" ", indentWidth)

		return r.renderWrappedLine(linkPrefix, line.Link.Display, style, cs.Reset, indent)

	case LineTypeListItem:
		bulletPrefix := prefix + cs.ListBullet + "• " + cs.Reset
		// Indent continuation lines to align with text after bullet
		indent := strings.Repeat(" ", len(prefix)+2)
		return r.renderWrappedLine(bulletPrefix, line.Text, cs.Text, cs.Reset, indent)

	case LineTypeQuote:
		quotePrefix := prefix + cs.Quote + "│ "
		// Indent continuation lines with quote bar
		indent := strings.Repeat(" ", len(prefix)) + cs.Quote + "│ " + cs.Reset
		return r.renderWrappedLine(quotePrefix, line.Text, cs.Quote, cs.Reset, indent)

	case LineTypePreformatted:
		// Don't wrap preformatted text
		return prefix + cs.Preformat + line.Text + cs.Reset

	case LineTypePreformatToggle:
		// Don't render the toggle lines themselves
		return ""

	case LineTypeText:
		if line.Text == "" {
			return "" // Empty line
		}
		return r.renderWrappedLine(prefix, line.Text, cs.Text, cs.Reset, prefix)

	default:
		return prefix + line.Text
	}
}

// renderWrappedLine renders a line with text wrapping
func (r *Renderer) renderWrappedLine(linePrefix, text, colorStart, colorEnd, contIndent string) string {
	if r.opts.Width <= 0 {
		// No wrapping
		return linePrefix + colorStart + text + colorEnd
	}

	// Calculate available width for text (accounting for prefix length without ANSI codes)
	// We need to count visible characters only, not ANSI escape codes
	visiblePrefixLen := len(stripANSI(linePrefix))
	availableWidth := r.opts.Width - visiblePrefixLen

	if availableWidth <= 10 {
		// Not enough space to wrap meaningfully
		return linePrefix + colorStart + text + colorEnd
	}

	// Wrap the text
	wrappedLines := wrapText(text, availableWidth, stripANSI(contIndent))

	if len(wrappedLines) == 0 {
		return linePrefix + colorStart + colorEnd
	}

	var result strings.Builder

	// First line uses the original prefix
	result.WriteString(linePrefix)
	result.WriteString(colorStart)
	result.WriteString(wrappedLines[0])
	result.WriteString(colorEnd)

	// Continuation lines use indent
	for i := 1; i < len(wrappedLines); i++ {
		result.WriteString("\n")
		result.WriteString(contIndent)
		result.WriteString(colorStart)
		result.WriteString(wrappedLines[i])
		result.WriteString(colorEnd)
	}

	return result.String()
}

// stripANSI removes ANSI escape codes from a string to get visible length
func stripANSI(s string) string {
	// Simple ANSI stripper for length calculation
	var result strings.Builder
	inEscape := false

	for _, r := range s {
		if r == '\033' {
			inEscape = true
		} else if inEscape && r == 'm' {
			inEscape = false
		} else if !inEscape {
			result.WriteRune(r)
		}
	}

	return result.String()
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
