// Package chooseui presents a simple console-based UI.
//
// The user-interface is constructed with an array of strings,
// and will allow the user to choose one of them.  The list may
// be filtered, and the user can cancel if they wish.
package chooseui

import (
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ChooseUI is the structure for holding our state.
type ChooseUI struct {

	// The items the user will choose from.
	Choices []string

	// The users' choice.
	chosen string
}

// New creates a new UI, allowing the user to select from the available options.
func New(choices []string) *ChooseUI {
	sort.Strings(choices)
	return &ChooseUI{Choices: choices}
}

// Choose launches our user interface.
func (ui *ChooseUI) Choose() string {

	//
	// Create the console-GUI application.
	//
	app := tview.NewApplication()

	//
	// Create a list to hold our files.
	//
	list := tview.NewList()
	list.ShowSecondaryText(false)
	list.SetWrapAround(false)

	//
	// Add all the choices to it.
	//
	for _, entry := range ui.Choices {
		list.AddItem(entry, "", ' ', nil)
	}

	//
	// If the user presses return in the list then choose that item.
	//
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			selected := list.GetCurrentItem()
			ui.chosen, _ = list.GetItemText(selected)
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
					ui.chosen, _ = list.GetItemText(selected)
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
			for _, entry := range ui.Choices {
				list.AddItem(entry, "", ' ', nil)
			}
			return
		}

		// Otherwise filter by input
		input = strings.ToLower(input)
		list.Clear()
		for _, entry := range ui.Choices {
			if strings.Contains(strings.ToLower(entry), input) {
				list.AddItem(entry, "", ' ', nil)
			}
		}

		return
	})

	//
	// Help text
	//
	help := tview.NewBox().SetBorder(true).SetTitle("TAB to switch focus, ENTER to select, ESC to cancel, arrows/etc to move")

	//
	// Create a layout grid, add the filter-box and the list.
	//
	grid := tview.NewFlex().SetFullScreen(true).SetDirection(tview.FlexRow)
	grid.AddItem(inputField, 1, 0, true)
	grid.AddItem(list, 0, 1, false)
	grid.AddItem(help, 2, 1, false)

	//
	// Global keyboard handler, use "TAB" to switch focus.
	//
	// Arrows and HOME/END work as expected regardless of focus-state
	//
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {

		// Home
		case tcell.KeyHome:
			list.SetCurrentItem(0)

		// End
		case tcell.KeyEnd:
			list.SetCurrentItem(list.GetItemCount())

		// Up arrow
		case tcell.KeyUp:
			selected := list.GetCurrentItem()
			if selected > 0 {
				selected--
			} else {
				selected = list.GetItemCount()
			}
			list.SetCurrentItem(selected)
			return nil

		// Down arrow
		case tcell.KeyDown:
			selected := list.GetCurrentItem()
			selected++
			list.SetCurrentItem(selected)
			return nil

		// TAB
		case tcell.KeyTab, tcell.KeyBacktab:
			if list.HasFocus() {
				app.SetFocus(inputField)
			} else {
				app.SetFocus(list)
			}
			return nil

		// Escape
		case tcell.KeyEscape:
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

	//
	// Return the choice
	//
	return ui.chosen
}
