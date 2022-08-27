package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func EditWindow(sn string, w fyne.Window) fyne.CanvasObject {
	entry1 := widget.NewEntry()
	w.Resize(fyne.Size{Width: 400, Height: 100})
	v1 := container.NewBorder(nil, nil, widget.NewLabel("label 1:"), nil, entry1)
	widget.NewLabel("lable 2")
	entry2 := widget.NewEntry()
	v2 := container.NewBorder(nil, nil, widget.NewLabel("label 2:"), nil, entry2)
	combox1 := widget.NewSelect([]string{"111", "222"}, func(s string) { fmt.Println("selected", s) })
	v3 := container.NewHBox(widget.NewLabel("label 3:"), combox1)
	okBtn := widget.NewButton(OK, func() {
		saveConfig(sn)
		w.Close()
	})
	cancelBtn := widget.NewButton(Cancel, func() {
		w.Close()
	})

	v4 := container.NewHBox(container.NewBorder(nil, nil, nil, nil, okBtn), container.NewBorder(nil, nil, nil, nil, cancelBtn))
	w.CenterOnScreen()
	return container.NewVBox(
		v1,
		v2,
		v3,
		v4,
	)
}

func saveConfig(sn string) {
	fmt.Println(sn)
}
