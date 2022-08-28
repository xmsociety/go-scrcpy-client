package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// ==============================================> SYSTRAY 任务栏设置
func SetupSystray(desk desktop.App, w fyne.Window) {
	// Set up menu
	menu := fyne.NewMenu(AppName,
		fyne.NewMenuItem(Open, w.Show),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Tomato", func() {
			// TODO FUNC
		}),
	)
	desk.SetSystemTrayMenu(menu)
}
