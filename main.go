package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"go-scrcpy-client/ui"
	"log"
	"os"
)

func main() {
	a := app.NewWithID("go-scrcpy-client")
	w := a.NewWindow("go-scrcpy-client")
	a.SetIcon(ui.Logo)
	ui.MainWindow(w)
	w.Resize(fyne.Size{Width: 140 * 6.5, Height: 500})
	w.SetFixedSize(true)
	w.CenterOnScreen()
	w.SetIcon(ui.Logo)

	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("go-scrcpy-client",
			fyne.NewMenuItem("Hello", func() {
				log.Println("Hello")
			}))
		desk.SetSystemTrayMenu(m)
	}

	w.ShowAndRun()

	err := os.Unsetenv("FYNE_FONT")
	if err != nil {
		return
	}
}
