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

	// app is the global application
	app *tview.Application

	// list contains the global list of text-entries.
	list *tview.List

	// inputField contains the global text-input field.
	inputField *tview.InputField
}

// New creates a new UI, allowing the user to select from the available options.
func New(choices []string) *ChooseUI {
	sort.Strings(choices)
	return &ChooseUI{Choices: choices}
}

// SetupUI configures the UI.
func (ui *ChooseUI) SetupUI() {

	//
	// Create the console-GUI application.
	//
	ui.app = tview.NewApplication()

	//
	// Create a list to hold our files.
	//
	ui.list = tview.NewList()
	ui.list.ShowSecondaryText(false)
	ui.list.SetWrapAround(false)

	//
	// Add all the choices to it.
	//
	for _, entry := range ui.Choices {
		ui.list.AddItem(entry, "", ' ', nil)
	}

	//
	// Create a filter input-view
	//
	ui.inputField = tview.NewInputField().
		SetLabel("Filter: ").
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {

				// get the selected index
				selected := ui.list.GetCurrentItem()

				// less than the entry count?
				if ui.list.GetItemCount() > 0 {
					ui.chosen, _ = ui.list.GetItemText(selected)
				}
				ui.app.Stop()
			}
		})

	//
	// Setup the filter-function, to filter the list to
	// only matches present in the input-field
	//
	ui.inputField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		// Get text
		input := strings.TrimSpace(currentText)

		// empty? All items should be visible
		if input == "" {
			ui.list.Clear()
			for _, entry := range ui.Choices {
				ui.list.AddItem(entry, "", ' ', nil)
			}
			return
		}

		// Otherwise filter by input
		input = strings.ToLower(input)
		ui.list.Clear()
		for _, entry := range ui.Choices {
			if strings.Contains(strings.ToLower(entry), input) {
				ui.list.AddItem(entry, "", ' ', nil)
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
	grid.AddItem(ui.inputField, 1, 0, true)
	grid.AddItem(ui.list, 0, 1, false)
	grid.AddItem(help, 2, 1, false)

	ui.app.SetRoot(grid, true).SetFocus(grid).EnableMouse(true)

}

// SetupKeyBinding installs the global captures, and list-specific keybindings.
func (ui *ChooseUI) SetupKeyBinding() {

	//
	// If the user presses return in the list then choose that item.
	//
	ui.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			selected := ui.list.GetCurrentItem()
			ui.chosen, _ = ui.list.GetItemText(selected)
			ui.app.Stop()
		}
		return event
	})

	//
	// Global keyboard handler, use "TAB" to switch focus.
	//
	// Arrows and HOME/END work as expected regardless of focus-state
	//
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {

		// Home
		case tcell.KeyHome:
			ui.list.SetCurrentItem(0)

		// End
		case tcell.KeyEnd:
			ui.list.SetCurrentItem(ui.list.GetItemCount())

		// Up arrow
		case tcell.KeyUp:
			selected := ui.list.GetCurrentItem()
			if selected > 0 {
				selected--
			} else {
				selected = ui.list.GetItemCount()
			}
			ui.list.SetCurrentItem(selected)
			return nil

		// Down arrow
		case tcell.KeyDown:
			selected := ui.list.GetCurrentItem()
			selected++
			ui.list.SetCurrentItem(selected)
			return nil

		// TAB
		case tcell.KeyTab, tcell.KeyBacktab:
			if ui.list.HasFocus() {
				ui.app.SetFocus(ui.inputField)
			} else {
				ui.app.SetFocus(ui.list)
			}
			return nil

		// Escape
		case tcell.KeyEscape:
			ui.app.Stop()
		}
		return event
	})

}

// Choose launches our user interface.
func (ui *ChooseUI) Choose() string {

	ui.SetupUI()

	ui.SetupKeyBinding()

	//
	// Launch the application.
	//
	err := ui.app.Run()
	if err != nil {
		panic(err)
	}

	//
	// Return the choice
	//
	return ui.chosen
}
