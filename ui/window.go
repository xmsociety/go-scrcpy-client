package ui

import (
	"context"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/xmsociety/adbutils"
	"go-scrcpy-client/scrcpy"
	"go-scrcpy-client/service"
	"image"
	"sort"
	"strconv"
	"sync"
	"time"
)

var (
	headersMap = map[int]string{
		0: "id",
		1: "check",
		2: "SerialNum",
		3: "NickName",
		4: "RunMode",
		5: "Run",
		6: "Operate",
		7: "Other",
	}
	devicesList      = make([]map[int]interface{}, 0)
	adb              = adbutils.AdbClient{Host: "localhost", Port: 5037, SocketTime: 10}
	textMap          = make(map[string]map[string]string)
	LiveMap          = make(map[string]fyne.Window)
	themeSettingOn   = false
	editMap          = make(map[string]fyne.Window)
	clientMap        = make(map[string]*scrcpy.Client)
	clientCancelMap  = make(map[string]context.CancelFunc)
	serviceMap       = make(map[string]*service.Service)
	checkBoxMap      = make(map[string]*widget.Check)
	serviceButtonMap = make(map[string]*widget.Button)
	allCheck         = &widget.Check{}
	allStartBtn      = &widget.Button{}
	allStopBtn       = &widget.Button{}
	opLock           = sync.Mutex{}
	maxWidthCol2     = 0
	maxWidthCol3     = 0
)

func NewClient(sn string, VideoTransfer chan image.Image, ErrReceiver chan error) *scrcpy.Client {
	if sn == "" {
		sn = "127.0.0.1:5555"
	}
	snNtid := adbutils.SerialNTransportID{
		Serial: sn,
	}
	return &scrcpy.Client{Device: adb.Device(snNtid), MaxWith: MaxWidth, Bitrate: 5000000, VideoSender: VideoTransfer, ErrReceiver: ErrReceiver}
}

func MainWindow(w fyne.Window) {
	InitParent(w)
	w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("About",
		fyne.NewMenuItem("about1", func() {}),
		fyne.NewMenuItem("about2", func() {}),
		// a quit item will be appended to our first menu
	), fyne.NewMenu("Settings",
		fyne.NewMenuItem("Theme", func() {
			if themeSettingOn {
				return
			}
			s := settings.NewSettings()
			appearance := s.LoadAppearanceScreen(w)
			tabs := container.NewAppTabs(
				&container.TabItem{Text: "Appearance", Icon: s.AppearanceIcon(), Content: appearance})
			tabs.SetTabLocation(container.TabLocationLeading)
			themeWindow := fyne.CurrentApp().NewWindow("Theme Settings")
			themeWindow.SetContent(tabs)
			themeWindow.Show()
			themeSettingOn = true
			themeWindow.SetOnClosed(func() {
				fmt.Println("close Theme Setting")
				themeSettingOn = false
			})
			fmt.Println("Menu New")
		}),
	), fyne.NewMenu("Help",
		fyne.NewMenuItem("", func() {
		})),
	))

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
		// create table
		func() fyne.CanvasObject {
			return container.NewMax(
				widget.NewLabel("placeholder"),
				widget.NewCheck("", func(b bool) {}),
				widget.NewLabel("placeholder"),
				widget.NewLabel("placeholder"),
				widget.NewLabel("placeholder"),
				NewButton(Run, func() { fmt.Println("Init Run") }),
				NewButton(Edit, func() { fmt.Println("Init Edit") }),
				NewButton(Show, func() { fmt.Println("Init Show") }),
			)
		},
		// update table
		func(id widget.TableCellID, c fyne.CanvasObject) {
			device := devicesList[id.Row]
			objs, ok := c.(*fyne.Container)
			if !ok {
				return
			}
			tmpWidthCol2 := 0
			tmpWidthCol3 := 0
			for i := 0; i < len(objs.Objects); i++ {
				obj := objs.Objects[i]
				sn := devicesList[id.Row][2].(*widget.Label).Text
				nickName := devicesList[id.Row][3].(*widget.Label).Text
				baseWidthCol2 := len(sn) * 10
				baseWidthCol3 := len(nickName) * 10

				if baseWidthCol2 > tmpWidthCol2 {
					tmpWidthCol2 = baseWidthCol2
				}
				if baseWidthCol3 > tmpWidthCol3 {
					tmpWidthCol3 = baseWidthCol3
				}

				if i == id.Col {
					switch obj.(type) {
					case *widget.Label:
						labObj := obj.(*widget.Label)
						labObj.SetText(device[id.Col].(*widget.Label).Text)
					case *widget.Button:
						buttonObj := obj.(*widget.Button)
						if id.Col == 7 {
							buttonObj.SetText(textMap[sn][Show])
						} else if id.Col == 6 {
							buttonObj.SetText(textMap[sn][Edit])
						} else { // 5
							buttonObj.SetText(textMap[sn][Run])
						}
						buttonObj.OnTapped = device[id.Col].(*widget.Button).OnTapped
					case *widget.Check:
						checkObj := obj.(*widget.Check)
						checkObj.OnChanged = device[id.Col].(*widget.Check).OnChanged
						checkObj.SetChecked(device[id.Col].(*widget.Check).Checked)
					}
					obj.Show()
				} else {
					obj.Hide()
				}
			}
			maxWidthCol2 = tmpWidthCol2
			maxWidthCol3 = tmpWidthCol3
		})
	go listenDevice(table)
	go autoAdjustTableWidth(headers, table)
	// select all
	allCheck = widget.NewCheck(CheckAll, func(b bool) {
		for _, check := range checkBoxMap {
			check.SetChecked(b)
			table.Refresh()
		}
	})
	// all check start
	allStartBtn = NewButton(AllStart, func() {
		startCount := 0
		for sn, check := range checkBoxMap {
			if check.Checked {
				if textMap[sn][Run] == Run {
					serviceButtonMap[sn].OnTapped()
				}
				startCount++
			}
		}
		if startCount == 0 {
			MessageError(errors.New("Nothing to Start!"))
			Notification(AppName, "Nothing to Start!")
		}
	})
	// all check stop
	allStopBtn = NewButton(AllStop, func() {
		stopCount := 0
		for sn, check := range checkBoxMap {
			if check.Checked {
				if textMap[sn][Run] == Running {
					serviceButtonMap[sn].OnTapped()
				}
				stopCount++
			}
		}
		if stopCount == 0 {
			MessageError(errors.New("Nothing to Stop!"))
			Notification(AppName, "Nothing to Stop!")
		}
	})
	bottom := container.NewHBox(allCheck, allStartBtn, allStopBtn)
	w.SetContent(container.NewBorder(container.NewBorder(head, nil, nil, nil, headers), bottom, nil, nil, table))
	w.SetMaster()
}

func setCurrentTime(head *widget.Label) {
	for {
		select {
		case now := <-time.Tick(time.Second * 1):
			head.SetText(fmt.Sprintf("Current Time: %v", now.Format("2006-01-02 15:04:05")))
		}
	}
}

func listenDevice(table *widget.Table) {
	for {
		var refreshFlag bool
		nowSn := []string{}
		allSn := []string{}
		devicesMap := map[string]map[int]interface{}{}
		for index, d := range adb.DeviceList() {
			sn := d.Serial
			nowSn = append(nowSn, sn)
			if listIn(sn) {
				continue
			}
			if _, ok := textMap[sn]; !ok {
				refreshFlag = true
				textMap[sn] = make(map[string]string, 0)
				textMap[sn][Show] = Show
				textMap[sn][Edit] = Edit
				textMap[sn][Run] = Run
				textMap[sn][Check] = False
				LiveMap[sn] = nil
				editMap[sn] = nil
				ch := make(chan image.Image)
				errCh := make(chan error)
				errChService := make(chan error)
				clientMap[sn] = NewClient(sn, ch, errCh)
				serviceMap[sn] = service.NewService(NewClient(sn, ch, errCh), errChService, "config")
			}
			checkWidget := widget.NewCheck("", func(b bool) {
				checkCount := 0
				checkBoxMap[sn].Checked = b
				for _, check := range checkBoxMap {
					if check.Checked {
						checkCount++
					}
				}
				allCheck.Checked = checkCount == len(checkBoxMap)
				allCheck.Refresh()
			})
			checkBoxMap[sn] = checkWidget
			runButton := NewButton(Run, func() {
				// stop quick click
				s := serviceMap[sn]
				opLock.Lock()
				defer func() {
					table.Refresh()
					time.Sleep(time.Second * 1)
					opLock.Unlock()
				}()
				if textMap[sn][Run] == Run {
					ctx, cancel := context.WithCancel(context.Background())
					s.Client.Ctx = ctx
					service.CancelMap[sn] = cancel
					s.Start()
					go serviceListener(s, table)
					textMap[sn][Run] = Running
					Notification(AppName, fmt.Sprintf("%v start run!", sn))
				} else {
					s.Stop()
					textMap[sn][Run] = Run
					Notification(AppName, fmt.Sprintf("%v stop!", sn))
				}

			})
			serviceButtonMap[sn] = runButton
			device := map[int]interface{}{
				0: widget.NewLabel(fmt.Sprintf("%v", index)), // index
				1: checkWidget,                               // check
				2: widget.NewLabel(sn),                       // sn
				3: widget.NewLabel("placeholder"),            // nick_name
				4: widget.NewLabel("placeholder"),            // run mode
				5: runButton,                                 // run
				6: NewButton(Edit, func() {
					// stop quick click
					opLock.Lock()
					defer func() {
						time.Sleep(time.Second * 1)
						opLock.Unlock()
					}()
					if textMap[sn][Edit] == Edit {
						textMap[sn][Edit] = EditIng
						w := fyne.CurrentApp().NewWindow(fmt.Sprintf("%s %s", sn, Edit))
						w.SetContent(EditWindow(sn, w))
						w.SetOnClosed(func() {
							textMap[sn][Edit] = Edit
							table.Refresh()
						})
						w.Show()
						editMap[sn] = w
					} else {
						editMap[sn].Close()      // page
						textMap[sn][Edit] = Edit // text
					}
					table.Refresh()
				}),
				7: NewButton(Show, func() {
					// stop quick click
					opLock.Lock()
					defer func() {
						time.Sleep(time.Second * 1)
						opLock.Unlock()
					}()
					if textMap[sn][Show] == Show {
						textMap[sn][Show] = Hide
						client := clientMap[sn]
						ctx, cancel := context.WithCancel(context.Background())
						clientCancelMap[sn] = cancel
						clientMap[sn].Ctx = ctx
						w := fyne.CurrentApp().NewWindow(fmt.Sprintf("%s %s", sn, Live))
						w.SetContent(ScreenWindow(w, client))
						w.Show()
						w.SetOnClosed(func() {
							textMap[sn][Show] = Show
							table.Refresh() // need Refresh
						})
						//stop close before Resize
						w.SetCloseIntercept(func() {
							FakeFunc(w, sn)
						})
						LiveMap[sn] = w
					} else {
						w := LiveMap[sn]
						FakeFunc(w, sn)
						textMap[d.Serial][Show] = Show
					}
					table.Refresh()
				}),
			}
			devicesList = append(devicesList, device)
		}
		// now all sn
		for _, device := range devicesList {
			sn := device[2].(*widget.Label).Text
			devicesMap[sn] = device
			allSn = append(allSn, sn)
		}
		toDeleteSn := Minus(allSn, nowSn)
		updateMap(toDeleteSn)
		if len(toDeleteSn) > 0 {
			tmp := []map[int]interface{}{}
			refreshFlag = true
			for _, sn := range toDeleteSn {
				if listIn(sn) {
					delete(devicesMap, sn)
				}
			}
			for _, device := range devicesMap {
				tmp = append(tmp, device)
			}
			sort.Slice(tmp, func(i, j int) bool {
				b, _ := strconv.Atoi(devicesList[i][0].(*widget.Label).Text)
				f, _ := strconv.Atoi(devicesList[j][0].(*widget.Label).Text)
				return b < f
			})
			devicesList = tmp
		}
		if len(nowSn) == 0 {
			refreshFlag = true
			devicesList = []map[int]interface{}{}
		}
		if refreshFlag {
			sort.Slice(devicesList, func(i, j int) bool {
				b, _ := strconv.Atoi(devicesList[i][0].(*widget.Label).Text)
				f, _ := strconv.Atoi(devicesList[j][0].(*widget.Label).Text)
				return b < f
			})
			for i := 0; i < len(devicesList); i++ {
				devicesList[i][0].(*widget.Label).SetText(fmt.Sprintf("%v", i))
			}
			table.Refresh()
		}
		time.Sleep(time.Second * 1)
	}
}

func autoAdjustTableWidth(table *widget.Table, headers *widget.Table) {
	for {
		select {
		case <-time.Tick(time.Millisecond * 500):
			table.SetColumnWidth(2, float32(maxWidthCol2))
			table.SetColumnWidth(3, float32(maxWidthCol3))
			headers.SetColumnWidth(2, float32(maxWidthCol2))
			headers.SetColumnWidth(3, float32(maxWidthCol3))
		}
	}
}

func listIn(item string) (in bool) {
	for _, dMap := range devicesList {
		if item == dMap[2].(*widget.Label).Text {
			return true
		}
	}
	return
}

func updateMap(toDeleteSn []string) {
	for _, sn := range toDeleteSn {
		delete(textMap, sn)
		delete(LiveMap, sn)
		delete(editMap, sn)
		delete(clientMap, sn)
		delete(clientCancelMap, sn)
		delete(serviceMap, sn)
		delete(service.CancelMap, sn)
		delete(checkBoxMap, sn)
		delete(serviceButtonMap, sn)
	}
}

func serviceListener(s *service.Service, table *widget.Table) {
	for {
		errText := <-s.ErrReceiver
		MessageError(errors.New(fmt.Sprintf("%s %s", s.Client.Device.Serial, errText)))
		s := serviceMap[s.Client.Device.Serial]
		s.Stop()
		textMap[s.Client.Device.Serial][Run] = Run
		Notification(AppName, fmt.Sprintf("%v stop!", s.Client.Device.Serial))
		table.Refresh()
	}
}
