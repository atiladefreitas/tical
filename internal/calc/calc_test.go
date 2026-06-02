package calc

import "testing"

// feed runs a compact script against a fresh calculator: digits and '.' are
// typed, operators are applied, and '=' resolves. Returns the final display.
func feed(script string) string {
	c := New()
	for i := 0; i < len(script); i++ {
		ch := script[i]
		switch ch {
		case '+', '-', '*', '/', '%':
			c.InputOperator(ch)
		case '=':
			c.Equals()
		case '.':
			c.InputDecimal()
		default:
			c.InputDigit(ch)
		}
	}
	return c.Display()
}

func TestBasicOperations(t *testing.T) {
	cases := map[string]string{
		"2+3=":      "5",
		"10-4=":     "6",
		"6*7=":      "42",
		"20/5=":     "4",
		"17%5=":     "2",
		"2+3*4=":    "20", // left-to-right, not precedence
		"100/8=":    "12.5",
		"0.1+0.2=":  "0.3",
		"9-12=":     "-3",
		"5*0=":      "0",
	}
	for in, want := range cases {
		if got := feed(in); got != want {
			t.Errorf("feed(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDivideByZero(t *testing.T) {
	if got := feed("5/0="); got != "Error" {
		t.Errorf("5/0 = %q, want Error", got)
	}
	if got := feed("5%0="); got != "Error" {
		t.Errorf("5%%0 = %q, want Error", got)
	}
}

func TestChaining(t *testing.T) {
	// 2 + 3 + 4 = should fold to 9 (intermediate 5 then +4).
	if got := feed("2+3+4="); got != "9" {
		t.Errorf("2+3+4 = %q, want 9", got)
	}
}

func TestClearAndBackspace(t *testing.T) {
	c := New()
	c.InputDigit('1')
	c.InputDigit('2')
	c.InputDigit('3')
	c.Backspace()
	if got := c.Display(); got != "12" {
		t.Errorf("after backspace = %q, want 12", got)
	}
	c.Clear()
	if got := c.Display(); got != "0" {
		t.Errorf("after clear = %q, want 0", got)
	}
}

func TestToggleSign(t *testing.T) {
	c := New()
	c.InputDigit('4')
	c.ToggleSign()
	if got := c.Display(); got != "-4" {
		t.Errorf("toggle sign = %q, want -4", got)
	}
	c.ToggleSign()
	if got := c.Display(); got != "4" {
		t.Errorf("toggle sign back = %q, want 4", got)
	}
}

func TestRecoverAfterError(t *testing.T) {
	c := New()
	c.InputDigit('5')
	c.InputOperator('/')
	c.InputDigit('0')
	c.Equals()
	if c.Display() != "Error" {
		t.Fatalf("expected Error state")
	}
	c.InputDigit('7') // typing should reset and start fresh
	if got := c.Display(); got != "7" {
		t.Errorf("after error recovery = %q, want 7", got)
	}
}
