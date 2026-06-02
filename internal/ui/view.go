package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model: it draws the panel and centres it so the
// calculator fills the full width and height of the terminal.
func (m Model) View() string {
	panel := m.renderPanel()
	if m.winW == 0 || m.winH == 0 { // size not known yet (first frame)
		return panel
	}
	return lipgloss.Place(
		m.winW, m.winH,
		lipgloss.Center, lipgloss.Center,
		panel,
		lipgloss.WithWhitespaceBackground(bg),
	)
}

// renderPanel builds the bordered calculator box (constant size).
func (m Model) renderPanel() string {
	body := lipgloss.JoinVertical(
		lipgloss.Left,
		renderHeader(m.st, m.c.Expr(), m.c.Display()),
		m.renderGrid(),
		m.renderStatus(),
	)
	helpView := m.st.help.Render(m.help.View(m.keys))
	return m.st.app.Render(lipgloss.JoinVertical(lipgloss.Left, body, helpView))
}

// renderStatus draws the transient feedback line (always one row tall, so the
// panel height — and therefore mouse hit-testing — stays constant).
func (m Model) renderStatus() string {
	base := lipgloss.NewStyle().
		Width(cols * cellW).Height(1).
		Align(lipgloss.Center).MarginTop(1)
	if m.status == "" {
		return base.Render(" ")
	}
	color := green
	if m.statusErr {
		color = red
	}
	return base.Foreground(color).Render(m.status)
}

// renderHeader draws the title and the calculator screen (expression + result).
func renderHeader(st styles, expr, display string) string {
	title := st.title.Render("  Tical") +
		lipgloss.NewStyle().Foreground(comment).Render("  · terminal calculator")

	if expr == "" {
		expr = " "
	}
	screen := st.screen.Render(lipgloss.JoinVertical(
		lipgloss.Right,
		st.expr.Render(expr),
		st.display.Render(display),
	))
	return lipgloss.JoinVertical(lipgloss.Left, title, "", screen)
}

// renderGrid draws the 5x4 button grid, styling the focused key distinctly.
func (m Model) renderGrid() string {
	lines := make([]string, 0, rows)
	for r := 0; r < rows; r++ {
		cells := make([]string, 0, cols)
		for c := 0; c < cols; c++ {
			cells = append(cells, m.renderButton(r, c))
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}
	// Separate rows by gapY blank lines for vertical breathing room.
	return strings.Join(lines, "\n"+strings.Repeat("\n", gapY))
}

// renderButton picks the right style for a key based on kind and focus.
func (m Model) renderButton(r, c int) string {
	b := grid[r][c]
	style := m.st.btn
	switch b.kind {
	case kOp:
		style = m.st.btnOp
	case kFn:
		style = m.st.btnFn
	case kEq:
		style = m.st.btnEq
	}
	if r == m.row && c == m.col {
		style = m.st.btnFocus
	}
	return style.Render(b.label)
}
