package templatedcmd

import (
	"testing"
)

// Basic test
func TestBasic(t *testing.T) {

	type TestCase struct {
		template string
		input    string
		split    string
		expected []string
	}

	tests := []TestCase{
		{"xine {}", "This is a file", "", []string{"xine", "This is a file"}},
		{"xine {1}", "This is a file", "", []string{"xine", "This"}},
		{"xine {3}", "This is a file", "", []string{"xine", "a"}},
		{"xine {10}", "This is a file", "", []string{"xine", ""}},
		{"foo bar", "", "", []string{"foo", "bar"}},
		{"id {1}", "root:0:0...", ":", []string{"id", "root"}},
	}

	for _, test := range tests {

		out := Expand(test.template, test.input, test.split)

		if len(out) != len(test.expected) {
			t.Fatalf("Expected to have %d pieces, found %d", len(test.expected), len(out))
		}

		for i, x := range test.expected {

			if out[i] != x {
				t.Errorf("expected '%s' for piece %d, got '%s'", x, i, out[i])
			}
		}
	}
}
