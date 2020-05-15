package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Structure for our options and state.
type chooseFileCommand struct {

	// Filenames we'll let the user choose between
	files []string

	// The users' choice
	chosen string
}

// Arguments adds per-command args to the object.
func (cf *chooseFileCommand) Arguments(f *flag.FlagSet) {
}

// Info returns the name of this subcommand.
func (cf *chooseFileCommand) Info() (string, string) {
	return "choose-file", `Choose a file, interactively.

Details:

This command presents a directory view, showing you all the files beneath
the named directory.  You can navigate with the keyboard, and press RETURN
to select a file.

Optionally you can press TAB to filter the list via an input field.

Uses:

This is ideal for choosing videos, roms, etc.  For example launch the
given video file:

   $ xine "$(sysbox choose-file ~/Videos)"`
}

// Execute is invoked if the user specifies `choose-file` as the subcommand.
func (cf *chooseFileCommand) Execute(args []string) int {

	//
	// Get our starting directory
	//
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	//
	// Find files
	//
	filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				if !strings.Contains(path, "/.") && !strings.HasPrefix(path, ".") {
					cf.files = append(cf.files, path)
				}
			}
			return nil
		})

	//
	// Create the console-GUI application.
	//
	app := tview.NewApplication()

	//
	// Create a list to hold our files.
	//
	list := tview.NewList()
	list.ShowSecondaryText(false)
	for _, entry := range cf.files {
		list.AddItem(entry, "", ' ', nil)
	}

	//
	// If the user presses return in the list then choose that item.
	//
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			selected := list.GetCurrentItem()
			cf.chosen, _ = list.GetItemText(selected)
			app.Stop()
		}

		return event
	})

	//
	// Create a filter input-view
	//
	inputField := tview.NewInputField().
		SetLabel("Filter: ").
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {

				// get the selected index
				selected := list.GetCurrentItem()

				// less than the entry count?
				if list.GetItemCount() > 0 {
					cf.chosen, _ = list.GetItemText(selected)
				}
				app.Stop()
			}
		})

	//
	// Setup the filter-function, to filter the list to
	// only matches present in the input-field
	//
	inputField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		// Get text
		input := strings.TrimSpace(currentText)

		// empty? All items should be visible
		if input == "" {
			list.Clear()
			for _, entry := range cf.files {
				list.AddItem(entry, "", ' ', nil)
			}
			return
		}

		// Otherwise filter by input
		input = strings.ToLower(input)
		list.Clear()
		for _, entry := range cf.files {
			if strings.Contains(strings.ToLower(entry), input) {
				list.AddItem(entry, "", ' ', nil)
			}
		}

		return
	})

	//
	// Create a layout grid, add the filter-box and the list.
	//
	grid := tview.NewFlex().SetFullScreen(true).SetDirection(tview.FlexRow)
	grid.AddItem(inputField, 1, 0, false)
	grid.AddItem(list, 0, 1, true)

	//
	// Global keyboard handler, use "TAB" to switch focus.
	//
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case '\t':
			if list.HasFocus() {
				app.SetFocus(inputField)
			} else {
				app.SetFocus(list)
			}
			return nil
		}
		if event.Key() == tcell.KeyEscape {
			app.Stop()
		}
		return event
	})

	//
	// Launch the application.
	//
	if err := app.SetRoot(grid, true).SetFocus(grid).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	if cf.chosen != "" {
		fmt.Printf("%s\n", cf.chosen)
		return 0
	}

	return 1

}
