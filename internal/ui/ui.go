// Package ui wires the calculator engine to a Bubble Tea program: a clickable,
// keyboard-navigable button grid rendered with Lip Gloss in Tokyo Night colors.
package ui

import (
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/atiladefreitas/tical/internal/calc"
)

const (
	cols     = 4
	rows     = 5
	cellW    = btnW + gapX // horizontal pitch: key width + gap
	rowPitch = btnH + gapY // vertical pitch: key height + gap
	gridLeft = 3           // outer border (1) + left padding (2)
	gridTop  = 2           // outer border (1) + top padding (1)
)

type btnKind int

const (
	kDigit btnKind = iota
	kOp
	kFn
	kEq
)

type action int

const (
	aDigit action = iota
	aOp
	aClear
	aBack
	aSign
	aDot
	aEq
)

// button describes a single key in the calculator grid.
type button struct {
	label string
	kind  btnKind
	act   action
	data  byte // the digit or operator this key feeds to the engine
}

// grid is the fixed 5x4 calculator layout.
var grid = [rows][cols]button{
	{
		{"C", kFn, aClear, 0},
		{"⌫", kFn, aBack, 0},
		{"%", kOp, aOp, '%'},
		{"÷", kOp, aOp, '/'},
	},
	{
		{"7", kDigit, aDigit, '7'},
		{"8", kDigit, aDigit, '8'},
		{"9", kDigit, aDigit, '9'},
		{"×", kOp, aOp, '*'},
	},
	{
		{"4", kDigit, aDigit, '4'},
		{"5", kDigit, aDigit, '5'},
		{"6", kDigit, aDigit, '6'},
		{"−", kOp, aOp, '-'},
	},
	{
		{"1", kDigit, aDigit, '1'},
		{"2", kDigit, aDigit, '2'},
		{"3", kDigit, aDigit, '3'},
		{"+", kOp, aOp, '+'},
	},
	{
		{"±", kFn, aSign, 0},
		{"0", kDigit, aDigit, '0'},
		{".", kDigit, aDot, 0},
		{"=", kEq, aEq, 0},
	},
}

type keymap struct {
	Up, Down, Left, Right key.Binding
	Equals, Press, Copy   key.Binding
	Quit, Help            key.Binding
}

func defaultKeys() keymap {
	return keymap{
		Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/↓/←/→", "move")),
		Down:   key.NewBinding(key.WithKeys("down", "j")),
		Left:   key.NewBinding(key.WithKeys("left", "h")),
		Right:  key.NewBinding(key.WithKeys("right", "l")),
		Equals: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "equals")),
		Press:  key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "press")),
		Copy:   key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "copy result")),
		Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c", "esc"), key.WithHelp("q", "quit")),
		Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	}
}

// ShortHelp implements help.KeyMap: the bottom bar shows only help and quit.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp implements help.KeyMap.
func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Equals, k.Press, k.Copy},
		{k.Help, k.Quit},
	}
}

// statusClearMsg wipes the transient status line (e.g. "copied ✓") after a beat.
type statusClearMsg struct{}

// Model is the Bubble Tea model for Tical.
type Model struct {
	c              *calc.Calculator
	st             styles
	keys           keymap
	help           help.Model
	row, col       int    // focused cell
	winW, winH     int    // terminal size, so the panel can be centred
	panelW, panelH int    // rendered panel size, for mouse offset maths
	status         string // transient feedback (copy result), "" when idle
	statusErr      bool   // colours the status line red on failure
}

// New returns an initialised Tical model.
func New() Model {
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(blue)
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(comment)
	h.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(comment)
	m := Model{
		c:    calc.New(),
		st:   newStyles(),
		keys: defaultKeys(),
		help: h,
		row:  4, col: 3, // start focused on "="
	}
	// The panel is a fixed size, so measure it once for mouse hit-testing.
	panel := m.renderPanel()
	m.panelW = lipgloss.Width(panel)
	m.panelH = lipgloss.Height(panel)
	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.MouseMsg:
		return m.handleMouse(msg)
	case tea.WindowSizeMsg:
		m.winW, m.winH = msg.Width, msg.Height
		return m, nil
	case statusClearMsg:
		m.status, m.statusErr = "", false
		return m, nil
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Any keypress clears stale status; the copy handler sets a fresh one.
	m.status, m.statusErr = "", false
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Copy):
		return m.copyResult()
	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll
		return m, nil
	case key.Matches(msg, m.keys.Up):
		m.row = (m.row - 1 + rows) % rows
		return m, nil
	case key.Matches(msg, m.keys.Down):
		m.row = (m.row + 1) % rows
		return m, nil
	case key.Matches(msg, m.keys.Left):
		m.col = (m.col - 1 + cols) % cols
		return m, nil
	case key.Matches(msg, m.keys.Right):
		m.col = (m.col + 1) % cols
		return m, nil
	case key.Matches(msg, m.keys.Equals):
		m.c.Equals()
		return m, nil
	case key.Matches(msg, m.keys.Press):
		m.press(grid[m.row][m.col])
		return m, nil
	}

	// Direct typing: route raw characters straight to the engine and move the
	// focus highlight to the matching key for visual feedback.
	if s := msg.String(); len(s) == 1 {
		switch ch := s[0]; ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			m.c.InputDigit(ch)
			m.focusData(ch)
		case '.', ',':
			m.c.InputDecimal()
		case '+', '-', '*', '/', '%':
			m.c.InputOperator(ch)
			m.focusData(ch)
		case '=':
			m.c.Equals()
		case 'c', 'C':
			m.c.Clear()
		}
	}
	switch msg.Type {
	case tea.KeyBackspace, tea.KeyDelete:
		m.c.Backspace()
	}
	return m, nil
}

func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	r, c, ok := m.cellAt(msg.X, msg.Y)
	if !ok {
		return m, nil
	}
	// Hover moves the focus highlight; a left press activates the key.
	m.row, m.col = r, c
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
		m.press(grid[r][c])
	}
	return m, nil
}

// press dispatches a button's action to the calculator engine.
func (m *Model) press(b button) {
	switch b.act {
	case aDigit:
		m.c.InputDigit(b.data)
	case aOp:
		m.c.InputOperator(b.data)
	case aClear:
		m.c.Clear()
	case aBack:
		m.c.Backspace()
	case aSign:
		m.c.ToggleSign()
	case aDot:
		m.c.InputDecimal()
	case aEq:
		m.c.Equals()
	}
}

// focusData points the focus highlight at the key carrying the given byte.
func (m *Model) focusData(data byte) {
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c].data == data {
				m.row, m.col = r, c
				return
			}
		}
	}
}

// cellAt converts terminal coordinates to a grid cell, if one is under them.
// It accounts for the panel being centred in the terminal and for the gaps
// between keys (clicks landing in a gap return ok == false).
func (m Model) cellAt(x, y int) (row, col int, ok bool) {
	offX, offY := m.offsets()
	gx := offX + gridLeft
	gy := offY + gridTop + headerHeight()

	relY := y - gy
	if relY < 0 || relY%rowPitch >= btnH { // above the grid or on a row gap
		return 0, 0, false
	}
	row = relY / rowPitch
	if row >= rows {
		return 0, 0, false
	}

	relX := x - gx
	if relX < 0 || relX%cellW >= btnW { // left of the grid or on a column gap
		return 0, 0, false
	}
	col = relX / cellW
	if col >= cols {
		return 0, 0, false
	}
	return row, col, true
}

// offsets returns the top-left corner where the centred panel is drawn.
func (m Model) offsets() (x, y int) {
	if m.winW > m.panelW {
		x = (m.winW - m.panelW) / 2
	}
	if m.winH > m.panelH {
		y = (m.winH - m.panelH) / 2
	}
	return x, y
}

// copyResult copies the current display value to the system clipboard.
func (m Model) copyResult() (tea.Model, tea.Cmd) {
	val := m.c.Display()
	if val == "Error" {
		m.status, m.statusErr = "nothing to copy", true
		return m, clearStatusCmd()
	}
	if err := clipboard.WriteAll(val); err != nil {
		m.status, m.statusErr = "clipboard unavailable", true
	} else {
		m.status, m.statusErr = "copied "+val+" ✓", false
	}
	return m, clearStatusCmd()
}

// clearStatusCmd hides the status line after a short delay.
func clearStatusCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg { return statusClearMsg{} })
}

// Init/headerHeight: the header (title + screen) sits above the grid; its line
// count tells cellAt where the grid begins. It is constant for a given layout.
func headerHeight() int {
	return lipgloss.Height(renderHeader(newStyles(), "", "0"))
}
