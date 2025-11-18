package parser

import (
	"strings"
	"testing"
)

func TestWrapText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		indent   string
		expected []string
	}{
		{
			name:     "no wrapping needed",
			text:     "short text",
			width:    50,
			indent:   "",
			expected: []string{"short text"},
		},
		{
			name:     "simple wrapping",
			text:     "this is a very long line that should be wrapped at word boundaries",
			width:    30,
			indent:   "",
			expected: []string{
				"this is a very long line that",
				"should be wrapped at word",
				"boundaries",
			},
		},
		{
			name:     "wrapping with indent",
			text:     "this is a very long line that should be wrapped with indentation",
			width:    30,
			indent:   "  ",
			expected: []string{
				"this is a very long line that",
				"  should be wrapped with",
				"  indentation",
			},
		},
		{
			name:     "no wrapping when width is 0",
			text:     "this should not be wrapped",
			width:    0,
			indent:   "",
			expected: []string{"this should not be wrapped"},
		},
		{
			name:     "empty text",
			text:     "",
			width:    50,
			indent:   "",
			expected: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapText(tt.text, tt.width, tt.indent)

			if len(result) != len(tt.expected) {
				t.Errorf("wrapText() returned %d lines, expected %d\nGot: %v\nExpected: %v",
					len(result), len(tt.expected), result, tt.expected)
				return
			}

			for i, line := range result {
				if line != tt.expected[i] {
					t.Errorf("wrapText() line %d = %q, expected %q", i, line, tt.expected[i])
				}
			}
		})
	}
}

func TestRenderWrapping(t *testing.T) {
	// Create a simple document with a long paragraph
	doc := &Document{
		Lines: []*Line{
			{
				Type: LineTypeText,
				Text: "This is a very long paragraph that should be wrapped at word boundaries when rendered with a limited width setting.",
			},
			{
				Type: LineTypeHeading1,
				Text: "This is a very long heading that should also be wrapped properly",
			},
		},
	}

	// Render with wrapping
	renderer := NewRenderer(&RenderOptions{
		Width:       50,
		NumberLinks: false,
		ColorScheme: &ColorScheme{}, // No colors for testing
	})

	result := renderer.Render(doc)

	// Check that the result contains newlines (indicating wrapping occurred)
	lines := strings.Split(result, "\n")
	if len(lines) <= 2 {
		t.Errorf("Expected text to be wrapped into multiple lines, got %d lines", len(lines))
	}

	// Verify no line is too long (accounting for ANSI codes and prefixes)
	for i, line := range lines {
		stripped := stripANSI(line)
		if len(stripped) > 52 { // Allow small margin for edge cases
			t.Errorf("Line %d is too long (%d chars): %q", i, len(stripped), stripped)
		}
	}
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no ANSI codes",
			input:    "plain text",
			expected: "plain text",
		},
		{
			name:     "with ANSI codes",
			input:    "\033[1;36mcolored text\033[0m",
			expected: "colored text",
		},
		{
			name:     "multiple ANSI codes",
			input:    "\033[32mgreen\033[0m and \033[31mred\033[0m",
			expected: "green and red",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripANSI(tt.input)
			if result != tt.expected {
				t.Errorf("stripANSI() = %q, expected %q", result, tt.expected)
			}
		})
	}
}
