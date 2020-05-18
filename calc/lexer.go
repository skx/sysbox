package calc

import (
	"fmt"
	"strconv"
	"strings"
)

// These constants are used to return the type of
// token which has been lexed.
const (
	// Basic token-types
	EOF    = "EOF"
	IDENT  = "IDENT"
	NUMBER = "NUMBER"
	ERROR  = "ERROR"

	// Assignment-magic
	LET    = "LET"
	ASSIGN = "="

	// Paren
	LPAREN = "("
	RPAREN = ")"

	// Operations
	PLUS     = "+"
	MINUS    = "-"
	MULTIPLY = "*"
	DIVIDE   = "/"
)

// Token holds a lexed token from our input.
type Token struct {

	// The type of the token.
	Type string

	// The value of the token.
	//
	// If the type of the token is NUMBER then this
	// will be stored as a float64.  Otherwise the
	// value will be a string.
	Value interface{}
}

// Lexer holds our lexer state.
type Lexer struct {

	// input is the string we're lexing
	input string

	// position is the current position within the input-string
	position int
}

// NewLexer creates a new lexer, for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// Next returns the next token from our input stream.
//
// This is pretty naive lexer, however it is sufficient to
// recognize numbers, identifiers, and our small set of
// operators.
func (l *Lexer) Next() *Token {

	// Known-token-types
	known := make(map[string]string)
	known["*"] = MULTIPLY
	known["+"] = PLUS
	known["-"] = MINUS
	known["/"] = DIVIDE
	known["="] = ASSIGN
	known["("] = LPAREN
	known[")"] = RPAREN

	// Loop until we've exhausted our input.
	for l.position < len(l.input) {

		// Get the next character
		char := string(l.input[l.position])

		// Is this a known character/token?
		t, ok := known[char]
		if ok {
			// skip the character, and return the token
			l.position++
			return &Token{Value: char, Type: t}
		}

		// If we reach here it is something more complex.
		switch char {

		// Skip whitespace
		case " ", "\t", "\n", ";":
			l.position++
			continue

			// Is it a digit?
		case "-", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ".":
			//
			// Loop for more digits
			//

			// Starting offset of our number
			start := l.position

			// ending offset of our number.
			end := l.position

			// keep walking forward, minding we don't wander
			// out of our input.
			for end < len(l.input) {
				c := string(l.input[end])

				if c != "0" &&
					c != "1" &&
					c != "2" &&
					c != "3" &&
					c != "4" &&
					c != "5" &&
					c != "6" &&
					c != "7" &&
					c != "8" &&
					c != "9" &&
					c != "." &&
					c != "-" {
					break
				}
				end++
			}

			l.position = end

			// Here we have the number
			token := l.input[start:end]

			// too many periods?
			bits := strings.Split(token, ".")
			if len(bits) > 2 {
				return &Token{Type: ERROR, Value: fmt.Sprintf("too many periods in '%s'", token)}
			}

			// We can only have a "-" at the start of the number
			for idx, chr := range token {
				if chr == rune('-') && idx != 0 {
					return &Token{Type: ERROR, Value: fmt.Sprintf("- can only appear at the start of a number, fond: %s", token)}
				}
			}

			// Convert to float64
			number, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return &Token{Value: err.Error(), Type: ERROR}
			}

			return &Token{Value: number, Type: NUMBER}

		default:
			//
			// We'll assume we have an identifier
			//

			// Starting offset of our ident
			start := l.position

			// ending offset of our ident.
			end := l.position

			// keep walking forward, minding we don't wander
			// out of our input.
			for end < len(l.input) {

				c := string(l.input[end])

				if c == " " ||
					c == "\t" ||
					c == "\n " ||
					c == "1" ||
					c == "2" ||
					c == "3" ||
					c == "4" ||
					c == "5" ||
					c == "6" ||
					c == "7" ||
					c == "8" ||
					c == "9" ||
					c == "." ||
					c == "+" ||
					c == "-" ||
					c == "/" ||
					c == "*" ||
					c == ";" ||
					c == "=" {
					break
				}
				end++
			}

			l.position = end
			token := l.input[start:end]

			// We only have a single keyword, LET, handle it here.
			if token == "let" {
				return &Token{Value: "let", Type: LET}
			}

			// If it wasn't `let` it was an identifier.
			return &Token{Value: token, Type: IDENT}
		}

	}

	return &Token{Value: "", Type: EOF}
}
