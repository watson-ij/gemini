package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/watson-ij/gemini/internal/ui"
)

func main() {
	// Default start URL
	startURL := "gemini://geminiprotocol.net"

	// Check for URL argument
	if len(os.Args) > 1 {
		startURL = os.Args[1]
	}

	// Create the model
	m := ui.NewModel(startURL)

	// Create the program
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
