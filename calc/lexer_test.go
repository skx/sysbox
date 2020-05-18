package calc

import (
	"fmt"
	"strings"
	"testing"
)

// Test basic invocation of our lexer.
func TestLexer(t *testing.T) {

	tests := []struct {
		expectedType    string
		expectedLiteral string
	}{
		{LET, "let"},
		{MULTIPLY, "*"},
		{ASSIGN, "="},
		{NUMBER, "3"},
		{PLUS, "+"},
		{NUMBER, "4"},
		{MULTIPLY, "*"},
		{NUMBER, "5"},
		{MINUS, "-"},
		{NUMBER, "1"},
		{DIVIDE, "/"},
		{NUMBER, "2"},
		{EOF, ""},
	}

	l := NewLexer("let * = 3 + 4 * 5 - 1 / 2")

	for i, tt := range tests {
		tok := l.Next()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if fmt.Sprintf("%v", tok.Value) != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Value)
		}
	}

}

// Test we can parse numbers correctly
func TestNumbers(t *testing.T) {

	tests := []struct {
		input  string
		error  bool
		errMsg string
	}{
		{"-3", false, ""},
		{".1", false, ""},
		{".1.1", true, "too many"},
		{"12-3", true, "only appear at the start"},
	}

	for n, test := range tests {

		l := NewLexer(test.input)

		// Loop over all tokens and see if we found an error
		err := ""

		tok := l.Next()
		for tok.Type != EOF {
			if tok.Type == ERROR {
				err = tok.Value.(string)
			}
			tok = l.Next()

		}

		if test.error {
			if err == "" {
				t.Fatalf("tests[%d] %s - expected error, got none", n, test.input)
			}
			if !strings.Contains(err, test.errMsg) {
				t.Fatalf("expected error to match '%s', but got '%s'", test.errMsg, err)
			}
		} else {
			if err != "" {
				t.Fatalf("tests[%d] %s - didn't expect error, got %s", n, test.input, err)
			}
		}
	}

}
