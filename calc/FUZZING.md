# Fuzz-Testing

The 1.18 release of the golang compiler/toolset has integrated support for
fuzz-testing.

Fuzz-testing is basically magical and involves generating new inputs "randomly"
and running test-cases with those inputs.


## Running

If you're running 1.18beta1 or higher you can run the fuzz-testing against
the calculator package like so:

    $ go test -fuzztime=300s -parallel=1 -fuzz=FuzzCalculator -v
    === RUN   TestBasic
    --- PASS: TestBasic (0.00s)
    ..
    ..
    fuzz: elapsed: 0s, gathering baseline coverage: 0/111 completed
    fuzz: elapsed: 0s, gathering baseline coverage: 111/111 completed, now fuzzing with 1 workers
    fuzz: elapsed: 3s, execs: 63894 (21292/sec), new interesting: 9 (total: 120)
    fuzz: elapsed: 6s, execs: 76044 (4051/sec), new interesting: 12 (total: 123)
    fuzz: elapsed: 9s, execs: 76044 (0/sec), new interesting: 12 (total: 123)
    fuzz: elapsed: 12s, execs: 76044 (0/sec), new interesting: 12 (total: 123)
    ..
    fuzz: elapsed: 5m0s, execs: 5209274 (12462/sec), new interesting: 162 (total: 273)
    fuzz: elapsed: 5m1s, execs: 5209274 (0/sec), new interesting: 162 (total: 273)
    --- PASS: FuzzCalculator (301.01s)
    PASS
    ok  	github.com/skx/sysbox/calc	301.010s

You'll note that I've added `-parallel=1` to the test, because otherwise my desktop system becomes unresponsive while the testing is going on.
