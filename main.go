package main

import (
	"katana/ui"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// Create a new Fyne application
	a := app.New()

	// Create the main window
	w := a.NewWindow("Katana Time Tracker")
	w.Resize(fyne.NewSize(400, 320)) // Initial size only
	// Do not call SetFixedSize or SetMinSize, allow full dynamic resizing

	// Create and set the main UI
	mainUI := ui.NewMainUI()
	w.SetContent(mainUI.Container)

	// Show and run
	w.ShowAndRun()
}
