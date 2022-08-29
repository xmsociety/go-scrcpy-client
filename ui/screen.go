package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"go-scrcpy-client/scrcpy"
	"log"
)

func ScreenWindow(w fyne.Window, client *scrcpy.Client) fyne.CanvasObject {
	w.Resize(fyne.NewSize(float32(ScreenBaseWidth), float32(ScreenBaseHeight)))
	imageLabel := NewOverRideImageWidget(nil, client)
	label := widget.NewLabelWithStyle("Frame Loading", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	content := container.NewCenter(imageLabel, label)
	go client.Start()
	go handler(label, imageLabel, w, client)
	w.CenterOnScreen()
	return content
}

func handler(l fyne.Widget, imageLabel *OverRideImageWidget, w fyne.Window, client *scrcpy.Client) {
	defer func() {
		log.Printf("[show] %v: goroutine show handler quit！", client.Device.Serial)
	}()
	for {
		select {
		case <-client.Ctx.Done():
			return
		case img := <-client.VideoSender:
			l.Hide()
			// h264.NewDecoder(h264.PixelFormatBGR) 拿出来BGR
			imageLabel.Image.Image = img
			imageLabel.Image.SetMinSize(fyne.NewSize(float32(client.Resolution.W), float32(client.Resolution.H)))
			if !client.Alive {
				return
			}
			// window close cause panic
			w.Resize(fyne.NewSize(float32(client.Resolution.W), float32(client.Resolution.H)))
			imageLabel.Refresh()
		case err := <-client.ErrReceiver:
			MessageError(err)
			FakeFunc(w, client.Device.Serial)
			return
		}
	}
}

// FakeFunc Before close Live window, do something
func FakeFunc(w fyne.Window, sn string) {
	defer func() {
		w.Close()
	}()
	textMap[sn][Show] = Show
	clientCancelMap[sn]()
	w.Hide()
}
