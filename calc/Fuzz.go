// +build gofuzz

package calc

// Fuzz is the function that our fuzzer-application uses.
// See `FUZZING.md` in our distribution for how to invoke it.
func Fuzz(data []byte) int {
	cal := New()
	cal.Load(string(data))
	cal.Run()
	return 1
}
