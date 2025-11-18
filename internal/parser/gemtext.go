package parser

import (
	"bufio"
	"io"
	"strings"
)

// Parser parses gemtext documents
type Parser struct {
	// options could go here if needed
}

// NewParser creates a new gemtext parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses a gemtext document from a reader
func (p *Parser) Parse(r io.Reader) (*Document, error) {
	doc := NewDocument()
	scanner := bufio.NewScanner(r)

	inPreformat := false
	var preformatAltText string

	for scanner.Scan() {
		rawLine := scanner.Text()
		line := p.parseLine(rawLine, &inPreformat, &preformatAltText)
		doc.AddLine(line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return doc, nil
}

// ParseString parses a gemtext document from a string
func (p *Parser) ParseString(s string) (*Document, error) {
	return p.Parse(strings.NewReader(s))
}

// parseLine parses a single line of gemtext
func (p *Parser) parseLine(raw string, inPreformat *bool, preformatAltText *string) *Line {
	line := &Line{
		Raw: raw,
	}

	// Check for preformat toggle
	if strings.HasPrefix(raw, "```") {
		*inPreformat = !*inPreformat
		line.Type = LineTypePreformatToggle

		if *inPreformat {
			// Starting a preformat block - capture alt text
			*preformatAltText = strings.TrimSpace(strings.TrimPrefix(raw, "```"))
			line.AltText = *preformatAltText
		} else {
			// Ending a preformat block
			line.AltText = *preformatAltText
			*preformatAltText = ""
		}

		return line
	}

	// If we're in a preformat block, everything is preformatted text
	if *inPreformat {
		line.Type = LineTypePreformatted
		line.Text = raw
		line.AltText = *preformatAltText
		return line
	}

	// Check for link line
	if strings.HasPrefix(raw, "=>") {
		line.Type = LineTypeLink
		line.Link = parseLink(raw)
		line.Text = line.Link.Display
		return line
	}

	// Check for heading lines
	if strings.HasPrefix(raw, "#") {
		if strings.HasPrefix(raw, "###") {
			line.Type = LineTypeHeading3
			line.Text = strings.TrimSpace(strings.TrimPrefix(raw, "###"))
			return line
		}
		if strings.HasPrefix(raw, "##") {
			line.Type = LineTypeHeading2
			line.Text = strings.TrimSpace(strings.TrimPrefix(raw, "##"))
			return line
		}
		if strings.HasPrefix(raw, "#") {
			line.Type = LineTypeHeading1
			line.Text = strings.TrimSpace(strings.TrimPrefix(raw, "#"))
			return line
		}
	}

	// Check for list item
	if strings.HasPrefix(raw, "* ") {
		line.Type = LineTypeListItem
		line.Text = strings.TrimPrefix(raw, "* ")
		return line
	}

	// Check for quote
	if strings.HasPrefix(raw, ">") {
		line.Type = LineTypeQuote
		// Quote can have optional space after >
		text := strings.TrimPrefix(raw, ">")
		line.Text = strings.TrimPrefix(text, " ")
		return line
	}

	// Default: regular text line
	line.Type = LineTypeText
	line.Text = raw
	return line
}

// parseLink parses a link line and extracts URL and label
// Format: => <URL> [<LABEL>]
func parseLink(raw string) *LinkInfo {
	// Remove the => prefix
	content := strings.TrimPrefix(raw, "=>")
	content = strings.TrimSpace(content)

	if content == "" {
		return &LinkInfo{
			URL:     "",
			Label:   "",
			Display: "",
		}
	}

	// Split on whitespace to separate URL from label
	parts := strings.SplitN(content, " ", 2)

	url := parts[0]
	label := ""

	if len(parts) > 1 {
		label = strings.TrimSpace(parts[1])
	}

	display := label
	if display == "" {
		display = url
	}

	return &LinkInfo{
		URL:     url,
		Label:   label,
		Display: display,
	}
}

// Parse is a convenience function that creates a parser and parses a document
func Parse(r io.Reader) (*Document, error) {
	p := NewParser()
	return p.Parse(r)
}

// ParseString is a convenience function that creates a parser and parses a string
func ParseString(s string) (*Document, error) {
	p := NewParser()
	return p.ParseString(s)
}
