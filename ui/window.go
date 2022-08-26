package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/xmsociety/adbutils"
	"go-scrcpy-client/scrcpy"
	"image"
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

var VideoTransfer = make(chan image.Image)

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

	imageLable := canvas.NewImageFromImage(nil)

	//container.NewBorder(nil, nil, nil, nil, imageLable)
	bottom := container.NewHBox(selectRadio, allStartBtn, allStopBtn)
	//w.SetContent(container.NewBorder(container.NewBorder(head, nil, nil, nil, headers), bottom, nil, nil, table))
	w.SetContent(container.NewBorder(container.NewBorder(head, nil, nil, nil, headers), bottom, nil, nil, imageLable))
	w.SetMaster()
	go ClientStart(imageLable, w)
}

func setCurrentTime(head *widget.Label) {
	for {
		head.SetText(fmt.Sprintf("Current Time: %v", time.Now().Format("2006-01-02 15:04:05")))
	}
}

func ClientStart(imageLable *canvas.Image, w fyne.Window) {
	adb := adbutils.AdbClient{Host: "localhost", Port: 5037, SocketTime: 10}
	snNtid := adbutils.SerialNTransportID{
		Serial: "127.0.0.1:5555",
	}
	fmt.Println(adb.Device(snNtid).SayHello())
	client := scrcpy.Client{Device: adb.Device(snNtid), MaxWith: 800, Bitrate: 5000000, VideoSender: VideoTransfer}
	go sendImage(imageLable, w, &client)
	go client.Start()

}

func sendImage(imageLable *canvas.Image, w fyne.Window, client *scrcpy.Client) {
	for {
		img := <-VideoTransfer
		// h264.NewDecoder(h264.PixelFormatBGR) 拿出来BGR
		imageLable.Image = img
		imageLable.SetMinSize(fyne.NewSize(float32(client.Resolution.W), float32(client.Resolution.H)))
		w.Resize(fyne.NewSize(float32(client.Resolution.W), float32(client.Resolution.H)))
		imageLable.Refresh()
	}
}
