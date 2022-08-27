package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"go-scrcpy-client/scrcpy"
)

func ScreenWindow(w fyne.Window, client *scrcpy.Client) fyne.CanvasObject {
	imageLabel := NewOverRideImageWidget(nil, client)
	go client.Start()
	label := widget.NewLabel("Frame Loading~~~")
	content := container.NewBorder(nil, nil, nil, nil, label, imageLabel)
	imageLabel.Image.SetMinSize(fyne.NewSize(float32(480), float32(640)))
	go sendImage(label, imageLabel, w, client)
	w.CenterOnScreen()
	return content
}

func sendImage(l fyne.Widget, imageLabel *OverRideImageWidget, w fyne.Window, client *scrcpy.Client) {
	for {
		img := <-client.VideoSender
		l.Hide()
		// h264.NewDecoder(h264.PixelFormatBGR) 拿出来BGR
		imageLabel.Image.Image = img
		imageLabel.Image.SetMinSize(fyne.NewSize(float32(client.Resolution.W), float32(client.Resolution.H)))
		w.Resize(fyne.NewSize(float32(client.Resolution.W), float32(client.Resolution.H)))
		imageLabel.Refresh()
	}
}
