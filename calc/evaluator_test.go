package calc

import (
	"math"
	"strings"
	"testing"
)

const float64EqualityThreshold = 1e-5

// Floating points are hard.
func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

// Test some basic operations
func TestBasic(t *testing.T) {

	tests := []struct {
		input  string
		output float64
	}{
		{"1", 1},
		{"1 + 2", 3},
		{"1 + 2 * 3", 7},
		{"1 / 3", 1.0 / 3},
		{"1 / 3 * 9", 3},
		{"( 1 / 3 ) * 9", 3},
		{"1 - 3", -2},
		{"-1 + 3", 2},
		{"( 1 + 2 ) * 4", 12},
		{"( ( 1 + 2 ) * 4 )", 12},
	}

	for _, test := range tests {

		p := New()
		p.Load(test.input)

		out := p.Run()

		if out.Type != NUMBER {
			t.Fatalf("Output was not a number: %v\n", out)
		}
		if !almostEqual(out.Value.(float64), test.output) {
			t.Fatalf("Got wrong result for '%s', expected '%f' found '%f'", test.input, test.output, out.Value.(float64))
		}
	}
}

// Test for errors
func TestDivideZero(t *testing.T) {

	tests := []struct {
		input string
	}{
		{"1 / 0"},
		{"let a = 1 ; let b = 0 ; a / b ;"},
	}

	for _, test := range tests {

		p := New()
		p.Load(test.input)

		out := p.Run()

		if out.Type != ERROR {
			t.Fatalf("expected error, found none")
		}
		if !strings.Contains(out.Value.(string), "division by zero") {
			t.Fatalf("division by zero error expected, but found %s", out.Value.(string))
		}
	}
}

// Test for errors
func TestMissingVariable(t *testing.T) {

	tests := []struct {
		input string
	}{
		{"let a = 1 + b"},
		{"let a = 1 - b"},
		{"let a = 1 / b"},
		{"let a = 1 * b"},

		{"let a =  b + 1"},
		{"let a =  b - 1"},
		{"let a =  b / 1"},
		{"let a =  b * 2"},
	}

	for _, test := range tests {

		p := New()
		p.Load(test.input)

		out := p.Run()

		if out.Type != ERROR {
			t.Fatalf("expected error, found none")
		}
		if !strings.Contains(out.Value.(string), "undefined variable") {
			t.Fatalf("undefined variable error expected, but found %s", out.Value.(string))
		}
	}
}

// TestErrorCases looks for some basic errors.
func TestErrorCases(t *testing.T) {

	tests := []struct {
		input string
		error string
	}{
		{"let 1 = 1", "is not an identifier"},
		{"let foo = ; ", "EOF"},
		{"let foo foo ; ", "not an assignment statement"},
		{"let foo = ( 1 + 2 * 3 ", "expected ')'"},
		{")", "Unexpected token inside factor"},
	}

	for _, test := range tests {

		p := New()
		p.Load(test.input)

		out := p.Run()

		if out.Type != ERROR {
			t.Fatalf("expected error, found none for input '%s'", test.input)
		}
		if !strings.Contains(out.Value.(string), test.error) {
			t.Fatalf("expected error '%s', but found %s", test.error, out.Value.(string))
		}
	}
}

// TestAssign tests that assignment work.
func TestAssign(t *testing.T) {
	tests := []struct {
		input    string
		variable string
		value    float64
	}{
		// with let
		{"let a = 3", "a", 3},
		{"let a = 1; let b = 2; let c = 3; let d = a+ b * c", "d", 7},

		// without let
		{"a = 6", "a", 6},
		{"a = 1; b = 2; c = 3;  d = a + b * c", "d", 7},
	}

	for _, test := range tests {

		p := New()
		p.Load(test.input)

		out := p.Run()

		if out.Type == ERROR {
			t.Fatalf("unexpected error '%s': %s", test.input, out.Value.(string))
		}

		// get the variable
		result, found := p.Variable(test.variable)
		if !found {
			t.Fatalf("failed to lookup variable %s in %s", test.variable, test.input)
		}

		if result != test.value {
			t.Fatalf("result of '%s' should have been %f, got %f", test.input, test.value, result)
		}
	}
}
