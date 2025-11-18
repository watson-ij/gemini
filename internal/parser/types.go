package parser

// LineType represents the type of a gemtext line
type LineType int

const (
	// LineTypeText is a regular text line
	LineTypeText LineType = iota

	// LineTypeLink is a link line (=>)
	LineTypeLink

	// LineTypeHeading1 is a level 1 heading (#)
	LineTypeHeading1

	// LineTypeHeading2 is a level 2 heading (##)
	LineTypeHeading2

	// LineTypeHeading3 is a level 3 heading (###)
	LineTypeHeading3

	// LineTypeListItem is a list item (*)
	LineTypeListItem

	// LineTypeQuote is a quote line (>)
	LineTypeQuote

	// LineTypePreformatted is preformatted text (within ``` blocks)
	LineTypePreformatted

	// LineTypePreformatToggle is the toggle line itself (```)
	LineTypePreformatToggle
)

// String returns a string representation of the line type
func (t LineType) String() string {
	switch t {
	case LineTypeText:
		return "Text"
	case LineTypeLink:
		return "Link"
	case LineTypeHeading1:
		return "Heading1"
	case LineTypeHeading2:
		return "Heading2"
	case LineTypeHeading3:
		return "Heading3"
	case LineTypeListItem:
		return "ListItem"
	case LineTypeQuote:
		return "Quote"
	case LineTypePreformatted:
		return "Preformatted"
	case LineTypePreformatToggle:
		return "PreformatToggle"
	default:
		return "Unknown"
	}
}

// IsHeading returns true if the line type is a heading
func (t LineType) IsHeading() bool {
	return t == LineTypeHeading1 || t == LineTypeHeading2 || t == LineTypeHeading3
}

// Line represents a single line in a gemtext document
type Line struct {
	// Type is the type of this line
	Type LineType

	// Raw is the raw line text
	Raw string

	// Text is the processed text content (without markup)
	Text string

	// Link contains link-specific information (only for link lines)
	Link *LinkInfo

	// AltText contains alt text for preformatted blocks
	AltText string
}

// LinkInfo contains information about a link
type LinkInfo struct {
	// URL is the link URL
	URL string

	// Label is the link label/text (may be empty)
	Label string

	// Display is what should be displayed to the user
	// If Label is empty, this is the URL
	Display string
}

// Document represents a parsed gemtext document
type Document struct {
	// Lines contains all lines in the document
	Lines []*Line

	// Links contains all links in the document (for quick access)
	Links []*Line

	// Headings contains all headings in the document (for TOC generation)
	Headings []*Line
}

// NewDocument creates a new empty document
func NewDocument() *Document {
	return &Document{
		Lines:    make([]*Line, 0),
		Links:    make([]*Line, 0),
		Headings: make([]*Line, 0),
	}
}

// AddLine adds a line to the document
func (d *Document) AddLine(line *Line) {
	d.Lines = append(d.Lines, line)

	// Add to links collection if it's a link
	if line.Type == LineTypeLink {
		d.Links = append(d.Links, line)
	}

	// Add to headings collection if it's a heading
	if line.Type.IsHeading() {
		d.Headings = append(d.Headings, line)
	}
}

// GetLink returns the link at the given index (0-based)
func (d *Document) GetLink(index int) *Line {
	if index < 0 || index >= len(d.Links) {
		return nil
	}
	return d.Links[index]
}

// LinkCount returns the number of links in the document
func (d *Document) LinkCount() int {
	return len(d.Links)
}

// HeadingCount returns the number of headings in the document
func (d *Document) HeadingCount() int {
	return len(d.Headings)
}

// LineCount returns the total number of lines in the document
func (d *Document) LineCount() int {
	return len(d.Lines)
}
