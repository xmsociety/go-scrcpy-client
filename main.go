package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"go-scrcpy-client/ui"
	"os"
)

func main() {
	a := app.NewWithID("go-scrcpy-client")
	w := a.NewWindow("go-scrcpy-client")
	ui.MainWindow(w)
	w.Resize(fyne.Size{Width: 640, Height: 400})
	w.CenterOnScreen()
	w.ShowAndRun()
	err := os.Unsetenv("FYNE_FONT")
	if err != nil {
		return
	}
}
