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
	label := widget.NewLabelWithStyle("Frame Loading", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	content := container.NewCenter(imageLabel, label)
	go handleMessage(client)
	go sendImage(label, imageLabel, w, client)
	w.Resize(fyne.NewSize(float32(ScreenBaseWidth), float32(ScreenBaseHeight)))
	w.CenterOnScreen()
	return content
}

func sendImage(l fyne.Widget, imageLabel *OverRideImageWidget, w fyne.Window, client *scrcpy.Client) {
	for {
		img := <-client.VideoSender
		if !client.Alive {
			return
		}
		l.Hide()
		// h264.NewDecoder(h264.PixelFormatBGR) 拿出来BGR
		imageLabel.Image.Image = img
		imageLabel.Image.SetMinSize(fyne.NewSize(float32(client.Resolution.W), float32(client.Resolution.H)))
		w.Resize(fyne.NewSize(float32(client.Resolution.W), float32(client.Resolution.H)))
		imageLabel.Refresh()
	}
}

func handleMessage(client *scrcpy.Client) {
	for {
		err := <-client.ErrReceiver
		if !client.Alive {
			return
		}
		MessageError(err)
		LiveMap[client.Device.Serial].Close()
		return
	}
}
