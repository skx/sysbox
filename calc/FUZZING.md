# Fuzz-Testing

If you don't have the appropriate tools installed you can fetch them via:

    go get github.com/dvyukov/go-fuzz/go-fuzz
    go get github.com/dvyukov/go-fuzz/go-fuzz-build

Now you can build the `calc` package with fuzzing enabled:

    $ go-fuzz-build github.com/skx/sysbox/calc

Create a location to hold the work, and give it copies of some sample-programs:

    $ mkdir -p workdir/corpus
    $ echo "1 + 2"     > workdir/corpus/1.txt
    $ echo "1 - 2"     > workdir/corpus/2.txt
    $ echo "1 / 2"     > workdir/corpus/3.txt
    $ echo "1 * 2"     > workdir/corpus/4.txt
    $ echo "1 + 2 * 3" > workdir/corpus/5.txt


Now you can actually launch the fuzzer - here I use `-procs 1` so that
my desktop system isn't complete overloaded:

    $ go-fuzz -procs 1 -bin calc-fuzz.zip  -workdir workdir/



## Output

Once the fuzzer starts running it will generate test-cases, and
run them.  The output will look something like this:

      2020/06/11 08:25:55 workers: 1, corpus: 240 (2m39s ago), crashers: 0, restarts: 1/9967, execs: 3508543 (4348/sec), cover: 846, uptime: 13m27s
      2020/06/11 08:25:58 workers: 1, corpus: 240 (2m42s ago), crashers: 0, restarts: 1/9954, execs: 3514055 (4338/sec), cover: 846, uptime: 13m30s

Here you see the fuzzer been running for 13 minutes, and 30 seconds, 3514055 test-cases have been executed and "crashers: 0" means there have been found no crashes.

If any crashes are detected they'll be saved in `workdir/crashers`.
