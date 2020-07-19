package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/skx/subcommands"

	expect "github.com/google/goexpect"
)

// Structure for our options and state.
type expectCommand struct {

	// Timeout for running our command/waiting.
	//
	// This is set via the script-file, rather than a command-line argument.
	timeout time.Duration

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags
}

// Info returns the name of this subcommand.
func (ec *expectCommand) Info() (string, string) {
	return "expect", `A simple utility for scripting interactive commands.

Details:

This command allows you to execute an arbitrary process, sending input for
matching output which is received.  It is a simple alternative to the 'expect'
utility, famously provided with TCL.

The command requires a configuration file to be specified which contains details
of the process to be executed, and the output/input to receive/send.

Here is a simple example, note that the output the command produces is matched via regular expressions, rather than literally.  That's why you'll see "\." used to match a literal period:

    # Comments are prefixed with '#'
    # Timeout is expressed in seconds
    TIMEOUT 10

    # The command to run
    SPAWN telnet telehack.com

    # Now the dialog
    EXPECT \n\.
    SEND   date\r\n
    EXPECT \n\.
    SEND   quit\r\n

You'll see we use '\r\n' because we're using telnet, for driving bash and other normal commands you'd use '\n' instead as you would expect:

    TIMEOUT 10
    SPAWN   /bin/bash --login
    EXPECT  $
    SEND    touch /tmp/meow\n
    EXPECT  $
    SEND    exit\n

If you wish to execute a command, or arguments, containing spaces that is supported via quoting:

    SPAWN /path/to/foo arg1 "argument two" arg3 ..
`
}

// Run a command, and return something suitable for matching against with
// the expect library we're using..
func (ec *expectCommand) expectExec(cmd []string) (*expect.GExpect, func() error, error) {

	c := exec.CommandContext(
		context.Background(),
		cmd[0], cmd[1:]...)

	// write error out to my stdout
	c.Stderr = os.Stderr

	stdIn, err := c.StdinPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("error creating pipe: %s", err)
	}

	stdOut, err := c.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("error creating pipe: %s", err)
	}

	if err = c.Start(); err != nil {
		return nil, nil, fmt.Errorf("unexpected error starting command: %+v", err)
	}

	waitCh := make(chan error, 1)

	e, _, err := expect.SpawnGeneric(
		&expect.GenOptions{
			In:  stdIn,
			Out: stdOut,
			Wait: func() error {
				er := c.Wait()
				waitCh <- er
				return err
			},
			Close: c.Process.Kill,
			Check: func() bool {
				if c.Process == nil {
					return false
				}
				return c.Process.Signal(syscall.Signal(0)) == nil
			},
		},
		ec.timeout,
		expect.Verbose(true),
		expect.VerboseWriter(os.Stdout),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating expect: %s", err)
	}

	wait := func() error {
		err := <-waitCh
		return err
	}

	return e, wait, nil
}

// Execute is invoked if the user specifies `expect` as the subcommand.
func (ec *expectCommand) Execute(args []string) int {

	// Ensure we have a config-file
	if len(args) <= 0 {
		fmt.Printf("Usage: expect /path/to/config.script\n")
		return 1
	}

	// We'll now open the configuration file
	handle, err := os.Open(args[0])
	if err != nil {
		fmt.Printf("error opening %s : %s\n", args[0], err.Error())
		return 1
	}

	// Timeout Value
	ec.timeout = 60 * time.Second

	// Command
	cmd := "/bin/sh"

	// Read/Send stuff
	interaction := []expect.Batcher{}

	// Allow reading line by line
	reader := bufio.NewReader(handle)

	line, err := reader.ReadString(byte('\n'))
	for err == nil {

		// Lose the space
		line = strings.TrimSpace(line)

		// Timeout?
		if strings.HasPrefix(line, "TIMEOUT ") {

			line = strings.TrimPrefix(line, "TIMEOUT ")
			line = strings.TrimSpace(line)

			val, er := strconv.Atoi(line)
			if er != nil {
				fmt.Printf("error converting timeout value %s to number: %s\n", line, er)
				return 1
			}

			ec.timeout = time.Duration(val) * time.Second
		}

		// Command
		if strings.HasPrefix(line, "SPAWN ") {
			cmd = strings.TrimPrefix(line, "SPAWN ")
			cmd = strings.TrimSpace(cmd)
		}

		// Expect
		if strings.HasPrefix(line, "EXPECT ") {

			line = strings.TrimPrefix(line, "EXPECT ")
			line = strings.TrimSpace(line)
			line = strings.ReplaceAll(line, "\\n", "\n")
			line = strings.ReplaceAll(line, "\\r", "\r")
			line = strings.ReplaceAll(line, "\\t", "\t")
			interaction = append(interaction, &expect.BExp{R: line})
		}

		// Send
		if strings.HasPrefix(line, "SEND ") {
			line = strings.TrimPrefix(line, "SEND ")
			line = strings.TrimSpace(line)
			line = strings.ReplaceAll(line, "\\n", "\n")
			line = strings.ReplaceAll(line, "\\r", "\r")
			line = strings.ReplaceAll(line, "\\t", "\t")
			interaction = append(interaction, &expect.BSnd{S: line})
		}

		// Loop again
		line, err = reader.ReadString(byte('\n'))
	}

	// Launch the command
	fmt.Printf("Running: '%s'\n", cmd)

	// Split the command into fields, taking into account quoted strings.
	//
	// So the user can run things like this:
	//   echo "foo bar" 3
	//
	// https://stackoverflow.com/questions/47489745/
	//
	r := csv.NewReader(strings.NewReader(cmd))
	r.Comma = ' '
	record, err := r.Read()
	if err != nil {
		fmt.Printf("failed to split %s : %s\n", cmd, err)
		return 1
	}

	// Launch the command using the record array we've just parsed.
	e, wait, err := ec.expectExec(record)
	if err != nil {
		fmt.Printf("error launching %s: %s\n", cmd, err)
		return 1
	}
	defer e.Close()

	// Wire up the expect-magic.
	_, err = e.ExpectBatch(interaction, ec.timeout)
	if err != nil {
		fmt.Printf("error running recipe:%s\n", err)
		return 1
	}

	// Now await completion of the command/process.
	if err := wait(); err != nil {
		fmt.Printf("error waiting for process: %s\n", err)
		return 1
	}

	return 0
}
