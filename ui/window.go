package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/xmsociety/adbutils"
	"go-scrcpy-client/scrcpy"
	"image"
	"sort"
	"strconv"
	"time"
)

var (
	headersMap = map[int]string{
		0: "id",
		1: "check",
		2: "NickName",
		3: "SerialNum",
		4: "RunMode",
		5: "Operate",
		6: "Other",
	}
	devicesList = make([]map[int]interface{}, 0)
	adb         = adbutils.AdbClient{Host: "localhost", Port: 5037, SocketTime: 10}
	textMap     = make(map[string]map[string]string)
	liveMap     = make(map[string]fyne.Window)
	editMap     = make(map[string]fyne.Window)
	clientMap   = make(map[string]*scrcpy.Client)
)

func NewClient(sn string, VideoTransfer chan image.Image) *scrcpy.Client {
	if sn == "" {
		sn = "127.0.0.1:5555"
	}
	snNtid := adbutils.SerialNTransportID{
		Serial: sn,
	}
	return &scrcpy.Client{Device: adb.Device(snNtid), MaxWith: 800, Bitrate: 5000000, VideoSender: VideoTransfer}
}

func tableSelected(id widget.TableCellID) {
	fmt.Println(id)
}

func tableUnSelected(id widget.TableCellID) {
	fmt.Println(id)
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
			c.(*widget.Label).SetText(headersMap[id.Col])
		})
	table := widget.NewTable(
		func() (int, int) { return len(devicesList), len(headersMap) },
		// create
		func() fyne.CanvasObject {
			return container.NewMax(
				widget.NewLabel("placeholder"),
				widget.NewCheckGroup([]string{""}, func(strings []string) { fmt.Println("check") }),
				widget.NewLabel("placeholder"),
				widget.NewLabel("placeholder"),
				widget.NewLabel("placeholder"),
				widget.NewButton(Edit, func() { fmt.Println("Init Edit") }),
				widget.NewButton(Show, func() { fmt.Println("Init Show") }),
			)
		},
		// update
		func(id widget.TableCellID, c fyne.CanvasObject) {
			device := devicesList[id.Row]
			objs, ok := c.(*fyne.Container)
			if !ok {
				return
			}
			for i := 0; i < len(objs.Objects); i++ {
				obj := objs.Objects[i]
				if i == id.Col {
					switch obj.(type) {
					case *widget.Label:
						labObj := obj.(*widget.Label)
						labObj.SetText(device[id.Col].(*widget.Label).Text)
					case *widget.Button:
						buttonObj := obj.(*widget.Button)
						if id.Col == 6 {
							sn := devicesList[id.Row][2].(*widget.Label).Text
							buttonObj.SetText(textMap[sn][Show])
						}
						buttonObj.OnTapped = device[id.Col].(*widget.Button).OnTapped
					case *widget.CheckGroup:
						checkObj := obj.(*widget.CheckGroup)
						checkObj.OnChanged = device[id.Col].(*widget.CheckGroup).OnChanged
					}
					obj.Show()
				} else {
					obj.Hide()
				}
			}

		})
	// bind event
	table.OnSelected = tableSelected
	table.OnUnselected = tableUnSelected

	for i := 0; i < len(headersMap); i++ {
		headers.SetColumnWidth(i, 120)
		table.SetColumnWidth(i, 120)
	}
	go ListenDevice(table)
	selectRadio := widget.NewRadioGroup([]string{SelectAll}, func(s string) {})
	allStartBtn := widget.NewButton(AllStart, func() {})
	allStopBtn := widget.NewButton(AllStop, func() {})
	//imageLabel := NewOverRideImageWidget(nil, c.Client)
	bottom := container.NewHBox(selectRadio, allStartBtn, allStopBtn)
	w.SetContent(container.NewBorder(container.NewBorder(head, nil, nil, nil, headers), bottom, nil, nil, table))
	w.SetMaster()
}

func setCurrentTime(head *widget.Label) {
	for {
		head.SetText(fmt.Sprintf("Current Time: %v", time.Now().Format("2006-01-02 15:04:05")))
	}
}

func ListenDevice(table *widget.Table) {
	for {
		newSn := []string{}
		oldSn := []string{}
		for _, device := range devicesList {
			oldSn = append(oldSn, device[2].(*widget.Label).Text)
		}
		for index, d := range adb.DeviceList() {
			newSn = append(newSn, d.Serial)
			if ListIn(d.Serial) {
				continue
			}
			if _, ok := textMap[d.Serial]; !ok {
				textMap[d.Serial] = make(map[string]string, 0)
				textMap[d.Serial][Show] = Show
				textMap[d.Serial][Edit] = Edit
				textMap[d.Serial][Check] = False
				liveMap[d.Serial] = nil
				editMap[d.Serial] = nil
				ch := make(chan image.Image)
				clientMap[d.Serial] = NewClient(d.Serial, ch)
			}

			device := map[int]interface{}{
				0: widget.NewLabel(fmt.Sprintf("%v", index)), // index
				1: widget.NewCheckGroup([]string{""}, func(strings []string) {
					func() {
						fmt.Println("this is check  ", d.Serial)
					}()
				}), // button
				2: widget.NewLabel(d.Serial),      // sn
				3: widget.NewLabel("placeholder"), // nick_name
				4: widget.NewLabel("placeholder"), // run mode
				5: widget.NewButton(Edit, func() {
					func() {
						if textMap[d.Serial][Edit] == Edit {
							textMap[d.Serial][Edit] = EditIng
							w := fyne.CurrentApp().NewWindow(fmt.Sprintf("%s %s", d.Serial, Edit))
							w.SetContent(EditWindow(d.Serial, w))
							w.Show()
							w.SetOnClosed(func() {
								textMap[d.Serial][Edit] = Edit
								table.Refresh()
							})
							editMap[d.Serial] = w
						} else {
							editMap[d.Serial].Close()
							textMap[d.Serial][Edit] = Edit
						}
						table.Refresh()
					}()
				}),
				6: widget.NewButton(Show, func() {
					func() {
						if textMap[d.Serial][Show] == Show {
							textMap[d.Serial][Show] = Hide
							client := clientMap[d.Serial]
							w := fyne.CurrentApp().NewWindow(fmt.Sprintf("%s %s", d.Serial, Live))
							w.SetContent(ScreenWindow(w, client))
							w.Show()
							w.SetOnClosed(func() {
								textMap[d.Serial][Show] = Show
								client.Stop()
								table.Refresh()
							})
							liveMap[d.Serial] = w
						} else {
							liveMap[d.Serial].Close()
							textMap[d.Serial][Show] = Show
						}

						table.Refresh()
					}()
				}),
			}
			devicesList = append(devicesList, device)
		}
		sort.Slice(devicesList, func(i, j int) bool {
			b, _ := strconv.Atoi(devicesList[i][0].(*widget.Label).Text)
			f, _ := strconv.Atoi(devicesList[j][0].(*widget.Label).Text)
			return b < f
		})
		if len(newSn) == 0 {
			devicesList = []map[int]interface{}{}
		}
		var refreshFlag bool
		sort.Slice(oldSn, func(i, j int) bool {
			return oldSn[i] < oldSn[j]
		})
		sort.Slice(newSn, func(i, j int) bool {
			return newSn[i] < newSn[j]
		})

		if len(oldSn) != len(newSn) {
			refreshFlag = true
		} else {
			for i := 0; i < len(newSn); i++ {
				if newSn[i] != oldSn[i] {
					refreshFlag = true
				}
			}
		}

		if refreshFlag {
			table.Refresh()
		}
		time.Sleep(time.Second * 1)
	}
}

func ListIn(item string) (in bool) {
	for _, dMap := range devicesList {
		if item == dMap[2].(*widget.Label).Text {
			return true
		}
	}
	return
}
