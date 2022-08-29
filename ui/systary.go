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
		fyne.NewMenuItem("All Check", func() {
			allCheck.SetChecked(!allCheck.Checked)
		}),
		fyne.NewMenuItem(AllStart, func() {
			allStartBtn.Tapped(new(fyne.PointEvent))
		}),
		fyne.NewMenuItem(AllStop, func() {
			allStopBtn.Tapped(new(fyne.PointEvent))
		}),
	)
	desk.SetSystemTrayMenu(menu)
}
