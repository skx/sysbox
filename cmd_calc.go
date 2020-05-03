package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"

	"github.com/skx/subcommands"
)

// Structure for our options and state.
type calcCommand struct {

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags
}

// Info returns the name of this subcommand.
func (c *calcCommand) Info() (string, string) {
	return "calc", `A simple (floating-point) calculator.

Details:

This command allows you to evaluate simple mathematical operations,
with support for floating-point operations - something the standard
'expr' command does not support.

Example:

   $ sysbox calc 3 + 3
   $ sysbox calc '1 / 3 * 9'

Note here we can join arguments, or accept a quoted string.  The arguments
must be quoted if you use '*' because otherwise the shell's globbing would
cause surprises.`
}

// eval evaluates the given AST expression.
func (c *calcCommand) eval(exp ast.Expr) float64 {
	switch exp := exp.(type) {

	// ! and -
	case *ast.BinaryExpr:
		return c.evalBinaryExpr(exp)

	// numbers (+ strings, etc)
	case *ast.BasicLit:
		switch exp.Kind {
		case token.INT, token.FLOAT:
			i, _ := strconv.ParseFloat(exp.Value, 64)
			return i
		default:
			fmt.Printf("unknown literal type: %v %T\n", exp, exp)
			os.Exit(1)
		}

	// parenthesis (e.g. "(1 + 2 ) * 3".)
	case *ast.ParenExpr:
		return (c.eval(exp.X))

	default:
		fmt.Printf("unknown ast.Node: %v %T\n", exp, exp)
		os.Exit(1)

	}

	return 0
}

// evalBinaryExpr evaluate a binary operation (which means there are
// two arguments).
func (c *calcCommand) evalBinaryExpr(exp *ast.BinaryExpr) float64 {
	left := c.eval(exp.X)
	right := c.eval(exp.Y)

	switch exp.Op {
	case token.ADD:
		return left + right
	case token.SUB:
		return left - right
	case token.MUL:
		return left * right
	case token.QUO:
		return left / right
	case token.REM:
		// modulus
		return float64(int(left) % int(right))
	}

	fmt.Printf("Unknown operator '%v'\n", exp.Op)
	os.Exit(1)
	return 0
}

// Execute is invoked if the user specifies `calc` as the subcommand.
func (c *calcCommand) Execute(args []string) int {

	//
	// Join all arguments, in case we have been given "3", "+", "4".
	//
	input := ""

	for _, arg := range args {
		input += arg
		input += " "
	}

	//
	// Parse to AST
	//
	exp, err := parser.ParseExpr(input)
	if err != nil {
		fmt.Printf("failed to parse '%s': %s\n", input, err)
		return 1
	}

	//
	// Evaluate
	//
	res := c.eval(exp)

	//
	// If the result is an int show that, to avoid
	// needless ".0000" suffix.
	//
	if res == float64(int(res)) {
		fmt.Printf("%d\n", int(res))
	} else {

		//
		// OK show the floating-point result.
		//
		fmt.Printf("%f\n", res)
	}

	//
	// All done
	//
	return 0
}
