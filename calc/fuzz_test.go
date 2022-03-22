//go:build go1.18
// +build go1.18

package calc

import (
	"testing"
)

// FuzzCalculator will run a series of random/fuzz-tests against our
// parser and evaluator.
func FuzzCalculator(f *testing.F) {

	// Seed some "interesting" inputs
	f.Add([]byte(""))
	f.Add([]byte("\r"))
	f.Add([]byte("\n"))
	f.Add([]byte("\t"))
	f.Add([]byte("\r \n \t"))
	f.Add([]byte("3 / 3\n"))
	f.Add([]byte("3 - -3\r\n"))
	f.Add([]byte("3 / 0"))
	f.Add([]byte(nil))

	// Run the fuzzer
	f.Fuzz(func(t *testing.T, input []byte) {

		// Create
		cal := New()

		// Parser
		cal.Load(string(input))

		// Evaluate
		cal.Run()
	})
}
