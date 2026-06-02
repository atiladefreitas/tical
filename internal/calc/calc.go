// Package calc implements an incremental calculator engine that behaves like a
// classic pocket calculator: digits build the current operand, an operator
// stores it, and the next operator (or "=") flushes the pending computation.
package calc

import (
	"math"
	"strconv"
	"strings"
)

// Calculator holds the state of an in-progress calculation.
type Calculator struct {
	display   string  // the operand currently being entered
	stored    float64 // the left-hand operand of a pending operation
	op        byte    // the pending operator, or 0 when none is set
	expr      string  // human-readable trail shown above the display
	overwrite bool    // when true the next digit replaces the display
	errored   bool    // set when the last operation was invalid (e.g. /0)
}

// New returns a Calculator reset to zero.
func New() *Calculator {
	return &Calculator{display: "0", overwrite: true}
}

// Display returns the operand currently shown to the user.
func (c *Calculator) Display() string {
	if c.errored {
		return "Error"
	}
	return c.display
}

// Expr returns the running expression trail (e.g. "12 ×").
func (c *Calculator) Expr() string { return c.expr }

// InputDigit appends a single 0-9 digit to the current operand.
func (c *Calculator) InputDigit(d byte) {
	if d < '0' || d > '9' {
		return
	}
	c.recoverFromError()
	if c.overwrite || c.display == "0" {
		c.display = string(d)
		c.overwrite = false
		return
	}
	if len(stripSign(c.display)) >= 15 { // keep the display readable
		return
	}
	c.display += string(d)
}

// InputDecimal adds a decimal point if the operand does not already have one.
func (c *Calculator) InputDecimal() {
	c.recoverFromError()
	if c.overwrite {
		c.display = "0."
		c.overwrite = false
		return
	}
	if !strings.Contains(c.display, ".") {
		c.display += "."
	}
}

// InputOperator flushes any pending operation and stores op as the next one.
// op must be one of '+', '-', '*', '/', '%'.
func (c *Calculator) InputOperator(op byte) {
	if c.errored {
		return
	}
	cur := c.value()
	if c.op != 0 && !c.overwrite {
		// Chain operations: 2 + 3 + → first resolve 2 + 3.
		c.stored = c.apply(c.stored, cur, c.op)
		c.display = format(c.stored)
	} else {
		c.stored = cur
	}
	if c.errored {
		return
	}
	c.op = op
	c.overwrite = true
	c.expr = format(c.stored) + " " + symbol(op)
}

// Equals resolves the pending operation and shows the result.
func (c *Calculator) Equals() {
	if c.errored || c.op == 0 {
		return
	}
	cur := c.value()
	result := c.apply(c.stored, cur, c.op)
	if c.errored {
		return
	}
	c.expr = format(c.stored) + " " + symbol(c.op) + " " + format(cur) + " ="
	c.display = format(result)
	c.stored = result
	c.op = 0
	c.overwrite = true
}

// Clear resets the calculator to its initial zero state.
func (c *Calculator) Clear() {
	c.display = "0"
	c.stored = 0
	c.op = 0
	c.expr = ""
	c.overwrite = true
	c.errored = false
}

// Backspace removes the last character of the current operand.
func (c *Calculator) Backspace() {
	if c.errored {
		c.Clear()
		return
	}
	if c.overwrite {
		return
	}
	if len(c.display) <= 1 || (len(c.display) == 2 && c.display[0] == '-') {
		c.display = "0"
		c.overwrite = true
		return
	}
	c.display = c.display[:len(c.display)-1]
}

// ToggleSign flips the sign of the current operand.
func (c *Calculator) ToggleSign() {
	if c.errored || c.display == "0" {
		return
	}
	if strings.HasPrefix(c.display, "-") {
		c.display = c.display[1:]
	} else {
		c.display = "-" + c.display
	}
}

// value parses the current display into a float64.
func (c *Calculator) value() float64 {
	v, err := strconv.ParseFloat(c.display, 64)
	if err != nil {
		return 0
	}
	return v
}

// apply performs a binary operation, flagging division/modulo by zero.
func (c *Calculator) apply(a, b float64, op byte) float64 {
	switch op {
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	case '/':
		if b == 0 {
			c.errored = true
			return 0
		}
		return a / b
	case '%':
		if b == 0 {
			c.errored = true
			return 0
		}
		return math.Mod(a, b)
	}
	return b
}

func (c *Calculator) recoverFromError() {
	if c.errored {
		c.Clear()
	}
}

// symbol maps an operator byte to its display glyph.
func symbol(op byte) string {
	switch op {
	case '+':
		return "+"
	case '-':
		return "−"
	case '*':
		return "×"
	case '/':
		return "÷"
	case '%':
		return "%"
	}
	return ""
}

// format renders a float64 without trailing zeros and with sane precision.
func format(f float64) string {
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return "Error"
	}
	s := strconv.FormatFloat(f, 'f', -1, 64)
	if strings.Contains(s, ".") && len(stripSign(s)) > 16 {
		s = strconv.FormatFloat(f, 'g', 12, 64)
	}
	return s
}

func stripSign(s string) string { return strings.TrimPrefix(s, "-") }
