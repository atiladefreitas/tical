package ui

import (
	"strings"
	"testing"
)

// TestCellGeometry renders the full view and confirms that cellAt maps the
// terminal coordinates of each visible button label back to its grid cell.
func TestCellGeometry(t *testing.T) {
	m := New()
	lines := strings.Split(m.View(), "\n")

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			label := grid[r][c].label
			// Expected center of the button: column offset + half the width.
			x := gridLeft + c*cellW + btnW/2
			y := gridTop + headerHeight() + r

			if y >= len(lines) {
				t.Fatalf("cell (%d,%d): y=%d beyond %d rendered lines", r, c, y, len(lines))
			}
			if !strings.Contains(lines[y], label) {
				t.Errorf("cell (%d,%d): label %q not on line %d: %q", r, c, label, y, lines[y])
			}
			gr, gc, ok := cellAt(x, y)
			if !ok || gr != r || gc != c {
				t.Errorf("cellAt(%d,%d) = (%d,%d,%v), want (%d,%d,true)", x, y, gr, gc, ok, r, c)
			}
		}
	}
}
