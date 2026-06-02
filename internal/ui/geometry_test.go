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

	// With no WindowSizeMsg yet, View renders the panel at the origin, so the
	// centring offset is zero and grid coordinates use the base offsets.
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			label := grid[r][c].label
			// Center of the button: column/row pitch offset + half the size.
			x := gridLeft + c*cellW + btnW/2
			y := gridTop + headerHeight() + r*rowPitch

			if y >= len(lines) {
				t.Fatalf("cell (%d,%d): y=%d beyond %d rendered lines", r, c, y, len(lines))
			}
			if !strings.Contains(lines[y], label) {
				t.Errorf("cell (%d,%d): label %q not on line %d: %q", r, c, label, y, lines[y])
			}
			gr, gc, ok := m.cellAt(x, y)
			if !ok || gr != r || gc != c {
				t.Errorf("cellAt(%d,%d) = (%d,%d,%v), want (%d,%d,true)", x, y, gr, gc, ok, r, c)
			}
		}
	}

	// A click in the gap between two rows must not register as any key.
	if _, _, ok := m.cellAt(gridLeft+btnW/2, gridTop+headerHeight()+btnH); ok && gapY > 0 {
		t.Errorf("click on the row gap should miss, but hit a cell")
	}
}
