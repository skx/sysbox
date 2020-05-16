package calc

import (
	"fmt"
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
