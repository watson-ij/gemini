package parser

import (
	"strings"
	"testing"
)

func TestParseText(t *testing.T) {
	input := "This is a normal text line"
	doc, err := ParseString(input)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(doc.Lines) != 1 {
		t.Fatalf("Expected 1 line, got %d", len(doc.Lines))
	}

	line := doc.Lines[0]
	if line.Type != LineTypeText {
		t.Errorf("Expected LineTypeText, got %v", line.Type)
	}

	if line.Text != input {
		t.Errorf("Expected text %q, got %q", input, line.Text)
	}
}

func TestParseHeadings(t *testing.T) {
	tests := []struct {
		input    string
		lineType LineType
		text     string
	}{
		{"# Heading 1", LineTypeHeading1, "Heading 1"},
		{"## Heading 2", LineTypeHeading2, "Heading 2"},
		{"### Heading 3", LineTypeHeading3, "Heading 3"},
		{"#No space", LineTypeHeading1, "No space"}, // Should still parse
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			doc, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("ParseString failed: %v", err)
			}

			if len(doc.Lines) != 1 {
				t.Fatalf("Expected 1 line, got %d", len(doc.Lines))
			}

			line := doc.Lines[0]
			if line.Type != tt.lineType {
				t.Errorf("Expected %v, got %v", tt.lineType, line.Type)
			}

			if line.Text != tt.text {
				t.Errorf("Expected text %q, got %q", tt.text, line.Text)
			}
		})
	}
}

func TestParseLink(t *testing.T) {
	tests := []struct {
		input   string
		url     string
		label   string
		display string
	}{
		{"=> gemini://example.com", "gemini://example.com", "", "gemini://example.com"},
		{"=> gemini://example.com Example Site", "gemini://example.com", "Example Site", "Example Site"},
		{"=>gemini://example.com", "gemini://example.com", "", "gemini://example.com"},
		{"=> gemini://example.com   Lots   of   spaces  ", "gemini://example.com", "Lots   of   spaces", "Lots   of   spaces"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			doc, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("ParseString failed: %v", err)
			}

			if len(doc.Lines) != 1 {
				t.Fatalf("Expected 1 line, got %d", len(doc.Lines))
			}

			line := doc.Lines[0]
			if line.Type != LineTypeLink {
				t.Errorf("Expected LineTypeLink, got %v", line.Type)
			}

			if line.Link == nil {
				t.Fatal("Link info is nil")
			}

			if line.Link.URL != tt.url {
				t.Errorf("Expected URL %q, got %q", tt.url, line.Link.URL)
			}

			if line.Link.Label != tt.label {
				t.Errorf("Expected label %q, got %q", tt.label, line.Link.Label)
			}

			if line.Link.Display != tt.display {
				t.Errorf("Expected display %q, got %q", tt.display, line.Link.Display)
			}
		})
	}
}

func TestParseListItem(t *testing.T) {
	input := "* List item"
	doc, err := ParseString(input)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(doc.Lines) != 1 {
		t.Fatalf("Expected 1 line, got %d", len(doc.Lines))
	}

	line := doc.Lines[0]
	if line.Type != LineTypeListItem {
		t.Errorf("Expected LineTypeListItem, got %v", line.Type)
	}

	if line.Text != "List item" {
		t.Errorf("Expected text %q, got %q", "List item", line.Text)
	}
}

func TestParseQuote(t *testing.T) {
	tests := []struct {
		input string
		text  string
	}{
		{">Quote", "Quote"},
		{"> Quote with space", "Quote with space"},
		{">", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			doc, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("ParseString failed: %v", err)
			}

			if len(doc.Lines) != 1 {
				t.Fatalf("Expected 1 line, got %d", len(doc.Lines))
			}

			line := doc.Lines[0]
			if line.Type != LineTypeQuote {
				t.Errorf("Expected LineTypeQuote, got %v", line.Type)
			}

			if line.Text != tt.text {
				t.Errorf("Expected text %q, got %q", tt.text, line.Text)
			}
		})
	}
}

func TestParsePreformatted(t *testing.T) {
	input := `Normal text
` + "```" + `go
func main() {
    fmt.Println("Hello")
}
` + "```" + `
More normal text`

	doc, err := ParseString(input)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	// Should have: text, toggle, preformat*3, toggle, text
	expectedTypes := []LineType{
		LineTypeText,
		LineTypePreformatToggle,
		LineTypePreformatted,
		LineTypePreformatted,
		LineTypePreformatted,
		LineTypePreformatToggle,
		LineTypeText,
	}

	if len(doc.Lines) != len(expectedTypes) {
		t.Fatalf("Expected %d lines, got %d", len(expectedTypes), len(doc.Lines))
	}

	for i, expectedType := range expectedTypes {
		if doc.Lines[i].Type != expectedType {
			t.Errorf("Line %d: expected %v, got %v", i, expectedType, doc.Lines[i].Type)
		}
	}

	// Check alt text
	if doc.Lines[1].AltText != "go" {
		t.Errorf("Expected alt text 'go', got %q", doc.Lines[1].AltText)
	}
}

func TestDocumentLinks(t *testing.T) {
	input := `# Welcome
=> gemini://example.com Link 1
Some text
=> gemini://another.com Link 2
* List item
=> gemini://third.com Link 3`

	doc, err := ParseString(input)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if doc.LinkCount() != 3 {
		t.Errorf("Expected 3 links, got %d", doc.LinkCount())
	}

	expectedURLs := []string{
		"gemini://example.com",
		"gemini://another.com",
		"gemini://third.com",
	}

	for i, expectedURL := range expectedURLs {
		link := doc.GetLink(i)
		if link == nil {
			t.Fatalf("Link %d is nil", i)
		}
		if link.Link.URL != expectedURL {
			t.Errorf("Link %d: expected URL %q, got %q", i, expectedURL, link.Link.URL)
		}
	}
}

func TestDocumentHeadings(t *testing.T) {
	input := `# H1
Text
## H2
More text
### H3`

	doc, err := ParseString(input)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if doc.HeadingCount() != 3 {
		t.Errorf("Expected 3 headings, got %d", doc.HeadingCount())
	}

	expectedTypes := []LineType{LineTypeHeading1, LineTypeHeading2, LineTypeHeading3}
	expectedTexts := []string{"H1", "H2", "H3"}

	for i, heading := range doc.Headings {
		if heading.Type != expectedTypes[i] {
			t.Errorf("Heading %d: expected type %v, got %v", i, expectedTypes[i], heading.Type)
		}
		if heading.Text != expectedTexts[i] {
			t.Errorf("Heading %d: expected text %q, got %q", i, expectedTexts[i], heading.Text)
		}
	}
}

func TestEmptyDocument(t *testing.T) {
	doc, err := ParseString("")
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if doc.LineCount() != 0 {
		t.Errorf("Expected 0 lines, got %d", doc.LineCount())
	}
}

func TestMultilineDocument(t *testing.T) {
	input := strings.Join([]string{
		"# Welcome",
		"",
		"This is a paragraph.",
		"",
		"=> gemini://example.com Example",
		"* Item 1",
		"* Item 2",
		"> A quote",
	}, "\n")

	doc, err := ParseString(input)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if doc.LineCount() != 8 {
		t.Errorf("Expected 8 lines, got %d", doc.LineCount())
	}

	if doc.LinkCount() != 1 {
		t.Errorf("Expected 1 link, got %d", doc.LinkCount())
	}

	if doc.HeadingCount() != 1 {
		t.Errorf("Expected 1 heading, got %d", doc.HeadingCount())
	}
}
