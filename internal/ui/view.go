package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model.
func (m Model) View() string {
	header := renderHeader(m.st, m.c.Expr(), m.c.Display())
	body := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		m.renderGrid(),
	)
	helpView := m.st.help.Render(m.help.View(m.keys))
	return m.st.app.Render(lipgloss.JoinVertical(lipgloss.Left, body, helpView)) + "\n"
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
	return strings.Join(lines, "\n")
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
