package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var parent fyne.Window

func InitParent(w fyne.Window) {
	parent = w
}

func Notification(title, content string) {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   title,
		Content: content,
	})
}

func MessageError(err error) {
	d := dialog.NewError(err, parent)
	d.Resize(fyne.NewSize(float32(DialogWidth), float32(DialogHeight)))
	d.Show()
}

func MessageInformation(title, content string) {
	dialog.NewInformation(title, content, parent).Show()
}

func MessageConfirm(title, message string, callback func(bool)) {
	dialog.NewConfirm(title, message, callback, parent).Show()
}

func PopError(err error) {
	content := container.NewVBox(widget.NewLabel(err.Error()), widget.NewButton(err.Error(), func() {
		parent.Close()
	}))

	widget.NewModalPopUp(content, parent.Canvas()).Show()
}
