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
	"log"
	"time"
)

var (
	headersMap = map[int]string{
		0: "id",
		1: "Device",
		2: "SerialNum",
		3: "RunMode",
		4: "Operate",
		5: "Other",
	}
	VideoTransfer = make(chan image.Image)
)

type ClientWithUi struct {
	Client *scrcpy.Client
}

func (c *ClientWithUi) SetClient(serial string) {
	if serial == "" {
		serial = "127.0.0.1:5555"
	}
	adb := adbutils.AdbClient{Host: "localhost", Port: 5037, SocketTime: 10}
	snNtid := adbutils.SerialNTransportID{
		Serial: serial,
	}
	fmt.Println(adb.Device(snNtid).SayHello())
	client := scrcpy.Client{Device: adb.Device(snNtid), MaxWith: 800, Bitrate: 5000000, VideoSender: VideoTransfer}
	c.Client = &client
}

func (c *ClientWithUi) Start() {
	c.Client.Start()
}

type OverRideImageWidget struct {
	widget.BaseWidget
	Image  *canvas.Image
	Client *scrcpy.Client
}

func (o *OverRideImageWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(o.Image)
}

func NewOverRideImageWidget(img image.Image, client *scrcpy.Client) *OverRideImageWidget {
	w := &OverRideImageWidget{Image: canvas.NewImageFromImage(img), Client: client}
	w.ExtendBaseWidget(w)
	return w
}

// Tapped 左键点击
func (o *OverRideImageWidget) Tapped(e *fyne.PointEvent) {
	fmt.Println("Tapped")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if o.Client.Control.ControlConn == nil {
		log.Println("o.Client.Control.ControlConn is nil")
		return
	}
	o.Client.Control.Touch(int(e.Position.X), int(e.Position.Y), scrcpy.ActionDown)
	o.Client.Control.Touch(int(e.Position.X), int(e.Position.Y), scrcpy.ActionUp)
}

// DoubleTapped 左键双击
func (o *OverRideImageWidget) DoubleTapped(e *fyne.PointEvent) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if o.Client.Control.ControlConn == nil {
		log.Println("o.Client.Control.ControlConn is nil")
		return
	}
	fmt.Println("DoubleTapped")
}

// TappedSecondary 右键点击
func (o *OverRideImageWidget) TappedSecondary(e *fyne.PointEvent) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if o.Client.Control.ControlConn == nil {
		log.Println("o.Client.Control.ControlConn is nil")
		return
	}
	fmt.Println("TappedSecondary")
}

func (o *OverRideImageWidget) Refresh() {
	canvas.Refresh(o)
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

	c := ClientWithUi{}
	c.SetClient("")

	imageLabel := NewOverRideImageWidget(nil, c.Client)
	bottom := container.NewHBox(selectRadio, allStartBtn, allStopBtn)
	w.SetContent(container.NewBorder(container.NewBorder(head, nil, nil, nil, headers), bottom, nil, nil, imageLabel))
	w.SetMaster()
	go c.Start()
	go sendImage(imageLabel, w, c.Client)
}

func setCurrentTime(head *widget.Label) {
	for {
		head.SetText(fmt.Sprintf("Current Time: %v", time.Now().Format("2006-01-02 15:04:05")))
	}
}

func sendImage(imageLabel *OverRideImageWidget, w fyne.Window, client *scrcpy.Client) {
	for {
		img := <-VideoTransfer
		// h264.NewDecoder(h264.PixelFormatBGR) 拿出来BGR
		imageLabel.Image.Image = img
		imageLabel.Image.SetMinSize(fyne.NewSize(float32(client.Resolution.W), float32(client.Resolution.H)))
		w.Resize(fyne.NewSize(float32(client.Resolution.W), float32(client.Resolution.H)))
		imageLabel.Refresh()
	}
}
