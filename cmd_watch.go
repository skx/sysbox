package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Structure for our options and state.
type watchCommand struct {

	// delay contains the number of seconds to sleep before updating our command.
	delay int

	// count increments once every second.
	count int

	// This can be set in the keyboard handler, and will trigger an immediate re-run
	// of the command, without disturbing the regularly scheduled update(s).
	immediately bool
}

// Arguments adds per-command args to the object.
func (w *watchCommand) Arguments(f *flag.FlagSet) {
	f.IntVar(&w.delay, "n", 5, "The number of seconds to sleep before re-running the specified command.")
}

// Info returns the name of this subcommand.
func (w *watchCommand) Info() (string, string) {
	return "watch", `Watch the output of a command.

Details:

This command allows you execute a command every five seconds,
and see the most recent output.

It is included because Mac OS does not include a watch-command
by default.

The display uses the tview text-based user interface package, to
present a somewhat graphical display - complete with an updating
run-timer.

To exit the application you may press 'q', 'Escape', or Ctrl-c.

`
}

// Execute is invoked if the user specifies `watch` as the subcommand.
func (w *watchCommand) Execute(args []string) int {

	if len(args) < 1 {
		fmt.Printf("Usage: watch cmd arg1 arg2 .. argN\n")
		return 1
	}

	// Command we're going to run
	command := strings.Join(args, " ")

	// Start time so that
	startTime := time.Now()

	// Assume Unix
	shell := "/bin/sh -c"

	switch runtime.GOOS {
	case "windows":
		shell = "cmd /c"
	}

	// Build up the thing to run
	sh := strings.Split(shell, " ")
	sh = append(sh, command)

	// Create the screen
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	err = screen.Init()
	if err != nil {
		panic(err)
	}

	// Create the application
	app := tview.NewApplication()
	app.SetScreen(screen)

	// Create the viewing-area
	viewer := tview.NewTextView()
	viewer.SetScrollable(true)
	viewer.SetBackgroundColor(tcell.ColorDefault)

	//
	// If the user presses 'q' or Esc in the viewer then exit
	//
	viewer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			app.Stop()
		}
		if event.Rune() == 'q' {
			app.Stop()
		}
		if event.Rune() == ' ' {
			// A space will trigger a re-run the next second
			w.immediately = true
		}
		return event
	})

	// Create an elapsed time record
	elapsed := tview.NewTextView()
	elapsed.SetTextColor(tcell.ColorBlack)
	elapsed.SetTextAlign(tview.AlignRight)
	elapsed.SetText("0s")
	elapsed.SetBackgroundColor(tcell.ColorGreen)

	// Setup a title
	title := tview.NewTextView()
	title.SetTextColor(tcell.ColorBlack)
	title.SetText(fmt.Sprintf("%s every %ds", command, w.delay))
	title.SetBackgroundColor(tcell.ColorGreen)

	// The status-bar will have the title and elapsed time
	statusBar := tview.NewFlex()
	statusBar.AddItem(title, 0, 1, false)
	statusBar.AddItem(elapsed, 15, 1, false)

	// The layout will have the status-bar
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(viewer, 0, 1, true)
	flex.AddItem(statusBar, 1, 1, false)
	app.SetRoot(flex, true)

	// Ensure we update
	go func() {
		run := true

		for {

			// Run the command if we should, either:
			//
			//  1.  The first time we start.
			//
			//  2. When the timer has exceeded our second-count
			if run || w.immediately {

				// Command output
				var out []byte

				// Run the command and get the output
				cmd := exec.Command(sh[0], sh[1:]...)

				// Get the output.
				out, err = cmd.CombinedOutput()
				if err != nil {
					app.Stop()
					fmt.Printf("Error running command: %v - %s\n", sh, err)
					os.Exit(1)
				}

				// Once we've done that we're all ready to update the screen
				app.QueueUpdateDraw(func() {

					// Clear the screen
					screen.Clear()

					// Update the main-window's output
					viewer.SetText(tview.TranslateANSI(string(out)))

					// And update our run-time log
					elapsed.SetText(fmt.Sprintf("%v", time.Since(startTime).Round(time.Second)))
				})

				run = false
			} else {

				// Otherwise just update the status-bars elapsed timer.
				app.QueueUpdateDraw(func() {
					elapsed.SetText(fmt.Sprintf("%v", time.Since(startTime).Round(time.Second)))
				})
			}

			// We sleep for a second, and want to reset the to-run flag when we've done that
			// enough times.
			w.count++
			if w.count >= w.delay {
				w.count = 0
				run = true
			}

			// The user can press the space-bar to trigger an immediate run,
			// reset the flag that would have been set in that case.
			if w.immediately {
				w.immediately = false
			}

			// delay before the next test.
			time.Sleep(time.Second)
		}
	}()

	// Run the application
	err = app.Run()
	if err != nil {
		fmt.Printf("Error in watch:%s\n", err)
		return 1
	}

	return 0
}
