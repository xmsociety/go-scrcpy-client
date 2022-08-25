package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

var headersMap = map[int]string{
	0: "id",
	1: "Device",
	2: "SerialNum",
	3: "RunMode",
	4: "Operate",
	5: "Other",
}

func MainWindow(w fyne.Window) {
	w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("File",
		fyne.NewMenuItem("New", func() { fmt.Println("Menu New") }),
		// a quit item will be appended to our first menu
	), fyne.NewMenu("Edit",
		fyne.NewMenuItem("Cut", func() { fmt.Println("Menu Cut") }),
		fyne.NewMenuItem("Copy", func() { fmt.Println("Menu Copy") }),
		fyne.NewMenuItem("Paste", func() { fmt.Println("Menu Paste") }),
	)))

	head := widget.NewLabel(fmt.Sprintf("Current Time: %v ", time.Now().Format("2006-01-02 15:04:05")))
	go setCurrentTime(head)
	headers := widget.NewTable(
		func() (int, int) { return 1, 6 },
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")
		},
		func(id widget.TableCellID, c fyne.CanvasObject) {
			//text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*widget.Label).SetText(headersMap[id.Col])
		})
	table := widget.NewTable(
		func() (int, int) { return 1, 6 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*widget.Label).SetText(text)
		})
	for i := 0; i < 7; i++ {
		table.SetColumnWidth(i, 100)
	}

	selectRadio := widget.NewRadioGroup([]string{"Select All"}, func(s string) {})
	allStartBtn := widget.NewButton("All Start", func() {})
	allStopBtn := widget.NewButton("All Stop", func() {})

	bottom := container.NewHBox(selectRadio, allStartBtn, allStopBtn)
	w.SetContent(container.NewBorder(container.NewBorder(head, nil, nil, nil, headers), bottom, nil, nil, table))
	w.SetMaster()
}

func setCurrentTime(head *widget.Label) {
	for {
		head.SetText(fmt.Sprintf("Current Time: %v", time.Now().Format("2006-01-02 15:04:05")))
	}
}
