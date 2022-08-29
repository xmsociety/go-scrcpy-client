package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/systray"
	"go-scrcpy-client/ui"
	"os"
)

func main() {
	App := app.NewWithID(ui.AppName)
	mainWindow := App.NewWindow(ui.AppName)
	App.SetIcon(ui.Logo)
	ui.MainWindow(mainWindow)
	if desk, ok := App.(desktop.App); ok {
		ui.SetupSystray(desk, mainWindow)
	}
	mainWindow.Resize(fyne.Size{Width: ui.Width, Height: ui.Height})
	mainWindow.SetFixedSize(true)
	mainWindow.CenterOnScreen()
	mainWindow.SetIcon(ui.Logo)

	// setting intercept not to close app, but hide window,
	// and close only via tray
	mainWindow.SetCloseIntercept(func() {
		ui.Notification(ui.AppName, fmt.Sprintf("%s  minimized!", ui.AppName))
		mainWindow.Hide()
	})

	mainWindow.Show()
	// start
	App.Lifecycle().SetOnStarted(func() {
		systray.SetTooltip(ui.AppName)
		systray.SetTitle(ui.AppName)
	})
	App.Run()
	err := os.Unsetenv("FYNE_FONT")
	if err != nil {
		return
	}
}
