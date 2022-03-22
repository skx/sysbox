package calc

import (
	"fmt"
	"math"
)

// Evaluator holds the state of the evaluation-object.
type Evaluator struct {

	// tokens holds the series of tokens which our
	// lexer produced from our input.
	tokens []*Token

	// Current position within the array of tokens.
	token int

	// holder for any variables the user has defined.
	variables map[string]float64
}

// New creates a new evaluation object.
//
// The evaluation object starts out as being empty,
// but you can call Load to load an expression and
// then Run to execute it.
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

// Variable allows you to return the value of the given variable
func (e *Evaluator) Variable(name string) (float64, bool) {
	res, ok := e.variables[name]
	return res, ok
}

// Load is used to load a program into the evaluator.
//
// Note that the existing variables will maintain their state
// if not reset explicitly.
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

	// Add an extra pair of EOF tokens so that nextToken
	// can always be called
	e.tokens = append(e.tokens, &Token{Value: "EOF", Type: EOF})
	e.tokens = append(e.tokens, &Token{Value: "EOF", Type: EOF})
	e.token = 0
}

// nextToken returns the next token from our input.
//
// When Load is called the input-expression is broken
// down into a series of Tokens, and this function
// advances to the next token, returning it.
func (e *Evaluator) nextToken() *Token {
	tok := e.tokens[e.token]
	e.token++
	return tok
}

// peekToken returns the next pending token in our stream.
//
// NOTE it is always possible to peek at the next token,
// because we deliberately add an extra/spare EOF token
// in our constructor.
func (e *Evaluator) peekToken() *Token {
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

	op := e.peekToken()
	for op.Type == MULTIPLY || op.Type == DIVIDE {

		op = e.nextToken()

		f2 := e.factor()

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

		op = e.peekToken()
	}

	if op.Type == ERROR {
		return &Token{Type: ERROR, Value: fmt.Sprintf("Unexpected token inside term() - %v\n", op)}
	}

	return f1
}

// expr() parse an expression
func (e *Evaluator) expr() *Token {

	t1 := e.term()

	//
	// If we have an assignment we'll save the result
	// of the expression here.
	//
	// We do this to avoid repetition for "let x = ..." and
	// "y = ..."
	//
	variable := &Token{Type: ERROR, Value: "cant happen"}

	//
	// Assignment without LET ?
	//
	if t1.Type == IDENT {
		nxt := e.peekToken()
		if nxt.Type == ASSIGN {

			// Skip the assignment
			e.nextToken()

			// And we've found a variable to assign to
			variable = t1
		}
	}

	//
	// Assignment with LET
	//
	if t1.Type == LET {

		// Get the identifier.
		ident := e.nextToken()
		if ident.Type != IDENT {
			return &Token{Type: ERROR, Value: fmt.Sprintf("%v is not an identifier", ident)}
		}

		// Skip the assignment statement
		assign := e.nextToken()
		if assign.Type != ASSIGN {
			return &Token{Type: ERROR, Value: fmt.Sprintf("%v is not an assignment statement", ident)}
		}

		variable = ident
	}

	//
	// OK if we have an assignment, of either form, then
	// process it here.
	//
	if variable.Type == IDENT {

		// Calculate the result
		result := e.expr()

		// Save it, and also return the value.
		if result.Type == NUMBER {
			e.variables[variable.Value.(string)] = result.Value.(float64)
		}
		return result
	}

	//
	// If we reach here we're now done with assignments.
	//
	tok := e.peekToken()
	for tok.Type == PLUS || tok.Type == MINUS {

		tok = e.nextToken()
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

		tok = e.peekToken()
	}

	if tok.Type == ERROR {
		return &Token{Type: ERROR, Value: fmt.Sprintf("Unexpected token inside expr() - %v\n", tok)}
	}

	return t1
}

// factor() - return a token
func (e *Evaluator) factor() *Token {
	tok := e.nextToken()

	switch tok.Type {
	case EOF:
		return &Token{Type: ERROR, Value: "unexpected EOF in factor()"}
	case NUMBER:
		return tok
	case IDENT:

		// sleazy hack here.
		//
		// We're getting a factor, but if we have a variable
		// AND the next token is an assignment then we return
		// the ident (i.e. current token) to allow Run() to
		// process "var = expr"
		//
		// Without this we'd interpret "foo = 1 + 2" as
		// a reference to the preexisting variable "foo"
		// which would not exist.
		//
		nxt := e.peekToken()
		if nxt.Type == ASSIGN {
			return tok
		}

		//
		// OK lookup the content of an existing variable.
		//
		val, ok := e.variables[tok.Value.(string)]
		if ok {
			return &Token{Value: val, Type: NUMBER}
		}
		return &Token{Type: ERROR, Value: fmt.Sprintf("undefined variable: %s", tok.Value.(string))}
	case LET:
		return tok
	case LPAREN:
		//
		// We don't need to skip past the `(` here
		// because the `expr` call will do that
		// to find its arguments
		//

		// evaluate the expression
		res := e.expr()

		// next token should be ")"
		if e.peekToken().Type != RPAREN {
			return &Token{Type: ERROR, Value: fmt.Sprintf("expected ')' after expression found %v", e.peekToken())}
		}

		// skip that ")"
		e.nextToken()

		return res
	case MINUS:
		// If the next token is a number then we're good
		if e.peekToken().Type == NUMBER {
			val := e.nextToken()
			cur := val.Value.(float64)
			val.Value = cur * -1
			return val
		}
	}

	return &Token{Type: ERROR, Value: fmt.Sprintf("Unexpected token inside factor() - %v\n", tok)}
}

// Run launches the program we've loaded.
//
// If multiple statements are available each are executed in turn,
// and the result of the last one returned.  However errors will
// cause early-termination.
func (e *Evaluator) Run() *Token {

	var result *Token

	// Process each statement
	for e.peekToken().Type != EOF && e.peekToken().Type != ERROR {

		// Get the result
		result = e.expr()

		// Error? Then abort
		if result.Type == ERROR {
			return result
		}

		// Otherwise loop again.
	}

	// Did we terminate on an error?
	if e.peekToken().Type == ERROR {
		return e.peekToken()
	}

	// If we evaluated something we'll have a result which
    // we'll save in the `result` variable.
    //
    // (We might receive input such as "", which will result
    // in nothing being evaluated)
	if result != nil {
		e.variables["result"] = result.Value.(float64)
	}

	// All done.
	return result
}
