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

// 鼠标键盘事件 绑定到某个widget 或者 container  很垃圾
type Image struct {
}

//左键点击
func (i *Image) Tapped(e *fyne.PointEvent) {
	fmt.Println("Tapped")
}

//左键双击
func (i *Image) DoubleTapped(e *fyne.PointEvent) {
	fmt.Println("DoubleTapped")
}

// 鼠标键盘事件 end

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
		func() (int, int) { return 1, len(headersMap) },
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")
		},
		func(id widget.TableCellID, c fyne.CanvasObject) {
			//text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*widget.Label).SetText(headersMap[id.Col])
		})
	table := widget.NewTable(
		func() (int, int) { return 1, len(headersMap) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*widget.Label).SetText(text)
		})
	for i := 0; i < len(headersMap); i++ {
		table.SetColumnWidth(i, 100)
	}

	selectRadio := widget.NewRadioGroup([]string{"Select All"}, func(s string) {})
	allStartBtn := widget.NewButton("All Start", func() {})
	allStopBtn := widget.NewButton("All Stop", func() {})

	imageLabel := canvas.NewImageFromImage(nil)

	//container.NewBorder(nil, nil, nil, nil, imageLabel)
	bottom := container.NewHBox(selectRadio, allStartBtn, allStopBtn)
	//w.SetContent(container.NewBorder(container.NewBorder(head, nil, nil, nil, headers), bottom, nil, nil, table))
	w.SetContent(container.NewBorder(container.NewBorder(head, nil, nil, nil, headers), bottom, nil, nil, imageLabel))
	w.SetMaster()
	go ClientStart(imageLabel, w)
}

func setCurrentTime(head *widget.Label) {
	for {
		head.SetText(fmt.Sprintf("Current Time: %v", time.Now().Format("2006-01-02 15:04:05")))
	}
}

func ClientStart(imageLabel *canvas.Image, w fyne.Window) {
	adb := adbutils.AdbClient{Host: "localhost", Port: 5037, SocketTime: 10}
	snNtid := adbutils.SerialNTransportID{
		Serial: "127.0.0.1:5555",
	}
	fmt.Println(adb.Device(snNtid).SayHello())
	client := scrcpy.Client{Device: adb.Device(snNtid), MaxWith: 800, Bitrate: 5000000, VideoSender: VideoTransfer}
	go sendImage(imageLabel, w, &client)
	go client.Start()

}

func sendImage(imageLabel *canvas.Image, w fyne.Window, client *scrcpy.Client) {
	for {
		img := <-VideoTransfer
		// h264.NewDecoder(h264.PixelFormatBGR) 拿出来BGR
		imageLabel.Image = img
		imageLabel.SetMinSize(fyne.NewSize(float32(client.Resolution.W), float32(client.Resolution.H)))
		w.Resize(fyne.NewSize(float32(client.Resolution.W), float32(client.Resolution.H)))
		imageLabel.Refresh()
	}
}
