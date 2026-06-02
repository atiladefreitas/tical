// Tical is a minimalist terminal calculator built with Bubble Tea, Bubbles and
// Lip Gloss, themed with Tokyo Night colors and driven by mouse or keyboard.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/atiladefreitas/tical/internal/ui"
)

func main() {
	p := tea.NewProgram(
		ui.New(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "tical:", err)
		os.Exit(1)
	}
}
