package ui

import "github.com/charmbracelet/lipgloss"

// Tokyo Night palette. https://github.com/folke/tokyonight.nvim
var (
	bg       = lipgloss.Color("#1a1b26")
	bgHL     = lipgloss.Color("#24283b")
	bgDark   = lipgloss.Color("#16161e")
	fg       = lipgloss.Color("#c0caf5")
	comment  = lipgloss.Color("#565f89")
	blue     = lipgloss.Color("#7aa2f7")
	cyan     = lipgloss.Color("#7dcfff")
	green    = lipgloss.Color("#9ece6a")
	magenta  = lipgloss.Color("#bb9af7")
	red      = lipgloss.Color("#f7768e")
	orange   = lipgloss.Color("#ff9e64")
	darkText = lipgloss.Color("#1a1b26")
)

// btnW and btnH are the inner dimensions of a single calculator key.
const (
	btnW = 7
	btnH = 1
)

type styles struct {
	app      lipgloss.Style
	title    lipgloss.Style
	expr     lipgloss.Style
	display  lipgloss.Style
	screen   lipgloss.Style
	btn      lipgloss.Style
	btnOp    lipgloss.Style
	btnFn    lipgloss.Style
	btnEq    lipgloss.Style
	btnFocus lipgloss.Style
	help     lipgloss.Style
}

// newStyles builds the full style set for the calculator chrome.
func newStyles() styles {
	base := lipgloss.NewStyle().
		Width(btnW).Height(btnH).
		Align(lipgloss.Center, lipgloss.Center).
		MarginRight(1)

	return styles{
		app: lipgloss.NewStyle().
			Padding(1, 2).
			Background(bg).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(magenta),

		title: lipgloss.NewStyle().
			Foreground(magenta).Bold(true),

		// screenW keeps the display box edges flush with the button grid:
		// grid = 4*cellW cols; box = content + 2 padding + 2 border.
		expr: lipgloss.NewStyle().
			Foreground(comment).
			Align(lipgloss.Right).
			Width(cols*cellW - 6),

		display: lipgloss.NewStyle().
			Foreground(fg).Bold(true).
			Align(lipgloss.Right).
			Width(cols*cellW - 6),

		screen: lipgloss.NewStyle().
			Background(bgDark).
			Padding(1, 2).
			MarginBottom(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(bgHL),

		btn: base.
			Foreground(fg).Background(bgHL),

		btnOp: base.
			Foreground(darkText).Background(blue).Bold(true),

		btnFn: base.
			Foreground(darkText).Background(red).Bold(true),

		btnEq: base.
			Foreground(darkText).Background(green).Bold(true),

		btnFocus: base.
			Foreground(darkText).Background(orange).Bold(true),

		help: lipgloss.NewStyle().
			Foreground(comment).
			MarginTop(1),
	}
}
