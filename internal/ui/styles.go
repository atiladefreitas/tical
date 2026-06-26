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
	dimBlue  = lipgloss.Color("#3d59a1") // utility keys (C, ⌫, ±): blue family, dialed back
	green    = lipgloss.Color("#9ece6a")
	magenta  = lipgloss.Color("#bb9af7")
	keyBg    = lipgloss.Color("#2a2f45") // raised surface so digit keys read as keys
	red      = lipgloss.Color("#f7768e") // copy-status errors only, never keys
	darkText = lipgloss.Color("#1a1b26")
)

// Key geometry: btnW/btnH are a single key's inner size; gapX/gapY are the
// breathing room between keys. The mouse hit-testing in ui.go derives its
// pitch from these same constants, so spacing stays click-accurate.
const (
	btnW = 9
	btnH = 1
	gapX = 1
	gapY = 0
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
		MarginRight(gapX)

	return styles{
		app: lipgloss.NewStyle().
			Padding(1, 2).
			Background(bg).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(blue),

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

		// Digits sit on a raised surface so they read as keys against the
		// panel without competing with the coloured action keys.
		btn: base.
			Foreground(fg).Background(keyBg),

		// Action keys stay in one colour family — utility dialed back, the
		// operators bright — so the eye groups them but still ranks them.
		btnFn: base.
			Foreground(fg).Background(dimBlue),

		btnOp: base.
			Foreground(darkText).Background(blue).Bold(true),

		btnEq: base.
			Foreground(darkText).Background(green).Bold(true),

		// Focus is the single brightest key on the grid.
		btnFocus: base.
			Foreground(darkText).Background(magenta).Bold(true),

		help: lipgloss.NewStyle().
			Foreground(comment).
			MarginTop(1),
	}
}
