package calc

import (
	"fmt"
	"math"
)

// Evaluator holds the parser-state
type Evaluator struct {

	// tokens holds the series of tokens which our lexer produced from our input.
	tokens []*Token

	// Current position within the array of tokens.
	token int

	// holder for variables
	variables map[string]float64
}

// New creates a new evaluation object.
func New() *Evaluator {

	// Create the new object.
	e := &Evaluator{}

	// Populate the variable storage-store.
	e.variables = make(map[string]float64)

	// Load default constants.
	e.variables["pi"] = math.Pi
	e.variables["e"] = math.E

	return e
}

// Load is used to load a program into the evaluator.
//
// Note that the existing variables will maintain their state
// if it isn't reset.
func (e *Evaluator) Load(input string) {

	// Create a lexer for splitting the program
	lexer := NewLexer(input)

	// Remove any existing tokens
	e.tokens = nil

	// Parse the input into tokens, and
	// save them away.
	for {
		tok := lexer.Next()
		if tok.Type == EOF {
			break
		}
		e.tokens = append(e.tokens, tok)
	}

	// Add an extra pair of EOF tokens so that NextToken
	// can always be called
	e.tokens = append(e.tokens, &Token{Value: "EOF", Type: EOF})
	e.tokens = append(e.tokens, &Token{Value: "EOF", Type: EOF})
	e.token = 0
}

// NextToken returns the next token from our input
func (e *Evaluator) NextToken() *Token {
	tok := e.tokens[e.token]
	e.token++
	return tok
}

// PeekToken returns the next token in our stream
//
// Note it is always possible to peek at the next token,
// because we deliberately add an extra/spare EOF token
// in our constructor.
func (e *Evaluator) PeekToken() *Token {
	tok := e.tokens[e.token]
	return tok
}

// term() - return a term
func (e *Evaluator) term() *Token {

	f1 := e.factor()

	// error handling
	if f1.Type == ERROR {
		return f1
	}

	op := e.PeekToken()
	for op.Type == MULTIPLY || op.Type == DIVIDE {

		op = e.NextToken()

		f2 := e.factor()

		// error handling
		if f2.Type == ERROR {
			return f2
		}

		if f1.Type != NUMBER {
			return &Token{Type: ERROR, Value: fmt.Sprintf("%v is not a number", f1)}
		}
		if f2.Type != NUMBER {
			return &Token{Type: ERROR, Value: fmt.Sprintf("%v is not a number", f2)}
		}

		if op.Type == MULTIPLY {
			f1.Value = f1.Value.(float64) * f2.Value.(float64)
		}
		if op.Type == DIVIDE {
			if f2.Value.(float64) == 0 {
				f1 = &Token{Type: ERROR, Value: fmt.Sprintf("Attempted division by zero: %v/%v", f1, f2)}
			} else {
				f1.Value = f1.Value.(float64) / f2.Value.(float64)
			}
		}

		op = e.PeekToken()
	}

	return f1
}

// expr() parse an expression
func (e *Evaluator) expr() *Token {

	t1 := e.term()

	// Sleazy.
	if t1.Type == LET {

		// Get the identifier.
		ident := e.NextToken()
		if ident.Type != IDENT {
			return &Token{Type: ERROR, Value: fmt.Sprintf("%v is not an identifier", ident)}
		}

		// Skip the assignment statement
		assign := e.NextToken()
		if assign.Type != ASSIGN {
			return &Token{Type: ERROR, Value: fmt.Sprintf("%v is not an assignment statement", ident)}
		}

		// Calculate the result
		result := e.expr()

		// Save it, and also return the value.
		if result.Type == NUMBER {
			e.variables[ident.Value.(string)] = result.Value.(float64)
		}
		return result
	}

	tok := e.PeekToken()
	for tok.Type == PLUS || tok.Type == MINUS {

		tok = e.NextToken()
		t2 := e.term()

		if t1.Type != NUMBER {
			return &Token{Type: ERROR, Value: fmt.Sprintf("%v is not a number", t1)}
		}
		if t2.Type != NUMBER {
			return &Token{Type: ERROR, Value: fmt.Sprintf("%v is not a number", t2)}
		}

		if tok.Type == PLUS {
			t1.Value = t1.Value.(float64) + t2.Value.(float64)
		}
		if tok.Type == MINUS {
			t1.Value = t1.Value.(float64) - t2.Value.(float64)
		}

		tok = e.PeekToken()
	}

	return t1
}

// factor() - return a token
func (e *Evaluator) factor() *Token {
	tok := e.NextToken()

	switch tok.Type {
	case EOF:
		return &Token{Type: ERROR, Value: "unexpected EOF in factor()"}
	case NUMBER:
		return tok
	case IDENT:
		val, ok := e.variables[tok.Value.(string)]
		if ok {
			return &Token{Value: val, Type: NUMBER}
		}
		return &Token{Type: ERROR, Value: fmt.Sprintf("undefined variable: %s", tok.Value.(string))}
	case LET:
		return tok

		// TODO: Handle LPAREN here
		// Until RPAREN
	}

	return &Token{Type: ERROR, Value: fmt.Sprintf("Unexpected token inside factor() - %v\n", tok)}
}

// Run launches the program we've loaded.
//
// If multiple statements are available each are executed in turn,
// and the result of the last one returned
func (e *Evaluator) Run() *Token {

	var result *Token

	// Process each statement
	for e.PeekToken().Type != EOF {

		// Get the result
		result = e.expr()

		// Error? Then abort
		if result.Type == ERROR {
			return result
		}
	}

	return result
}
