package scrcpy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
	"sync"
	"time"
)

type ControlSender struct {
	ControlConn net.Conn
	W           int
	H           int
	Lock        sync.Mutex
}

func (control *ControlSender) Keycode(keyCode, action, repeat int) {
	control.Lock.Lock()
	defer control.Lock.Unlock()

	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint8(TypeInjectKeycode), // base
		uint8(action),            // B 1 byte
		uint32(keyCode),          // i 4 byte
		uint32(repeat),           // i 4 byte
		uint32(0),                // i 4 byte
	}
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("Keycode binary.Write failed:", err)
			return
		}
	}
	_, err := control.ControlConn.Write(buf.Bytes())
	if err != nil {
		log.Fatal("send Keycode error! ", err.Error())
		return
	}
}

func (control *ControlSender) Text(text string) {
	control.Lock.Lock()
	defer control.Lock.Unlock()

	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint8(TypeInjectTEXT),     // base
		uint32(len([]byte(text))), // i 4 byte
	}
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("Text binary.Write failed:", err)
		}
	}
	msg := append(buf.Bytes(), []byte(text)...)
	_, err := control.ControlConn.Write(msg)
	if err != nil {
		log.Fatal("send Text error! ", err.Error())
		return
	}

}

func (control *ControlSender) Touch(x, y, action int) {
	control.Lock.Lock()
	defer control.Lock.Unlock()

	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint8(TypeInjectTOUCHEvent), // base
		uint8(action),               // B 1 byte
		int64(-1),                   // q 8 byte
		uint32(x),                   // i 4 byte
		uint32(y),                   // i 4 byte
		uint16(control.W),           // H 2 byte
		uint16(control.H),           // H 2 byte
		uint16(0xffff),              // H 2 byte
		uint32(1),                   // i 4 byte
	}
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("Touch binary.Write failed:", err)
		}
	}
	_, err := control.ControlConn.Write(buf.Bytes())
	if err != nil {
		log.Println("send Touch error! ", err.Error())
	}
}

func (control *ControlSender) Scroll(x, y, h, v int) {
	control.Lock.Lock()
	defer control.Lock.Unlock()
	X := math.Max(float64(x), 0)
	Y := math.Max(float64(y), 0)

	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint8(TypeInjectTOUCHEvent), // base
		uint32(X),                   // i 4 byte
		uint32(Y),                   // i 4 byte
		uint16(control.W),           // H 2 byte
		uint16(control.H),           // H 2 byte
		uint32(h),                   // i 4 byte
		uint32(v),                   // i 4 byte
	}
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("TypeInjectTOUCHEvent binary.Write failed:", err)
		}
	}
	_, err := control.ControlConn.Write(buf.Bytes())
	if err != nil {
		log.Fatal("send TypeInjectTOUCHEvent error! ", err.Error())
		return
	}
}

func (control *ControlSender) BackOrTurnScreenOn(action int) {
	control.Lock.Lock()
	defer control.Lock.Unlock()

	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint8(TypeBACKORScreenON), // base
		uint8(action),             // B 1 byte
	}
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("TypeBACKORScreenON binary.Write failed:", err)
		}
	}
	_, err := control.ControlConn.Write(buf.Bytes())
	if err != nil {
		log.Fatal("send TypeBACKORScreenON error! ", err.Error())
		return
	}
}

func (control *ControlSender) ExpandNotificationPanel() {
	control.Lock.Lock()
	defer control.Lock.Unlock()

	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint8(TypeEXPANDNOTIFICATIONPANEL), // base
	}
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("ExpandNotificationPanel binary.Write failed:", err)
		}
	}
	_, err := control.ControlConn.Write(buf.Bytes())
	if err != nil {
		log.Fatal("send ExpandNotificationPanel error! ", err.Error())
		return
	}
}

func (control *ControlSender) ExpandSettingsPanel() {
	control.Lock.Lock()
	defer control.Lock.Unlock()

	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint8(TypeEXPANDSETTINGSPANEL), // base
	}
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("TypeEXPANDSETTINGSPANEL binary.Write failed:", err)
		}
	}
	_, err := control.ControlConn.Write(buf.Bytes())
	if err != nil {
		log.Fatal("send TypeEXPANDSETTINGSPANEL error! ", err.Error())
		return
	}
}

func (control *ControlSender) CollapsePanels() {
	control.Lock.Lock()
	defer control.Lock.Unlock()

	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint8(TypeCOLLAPSEPANELS), // base
	}
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("TypeCOLLAPSEPANELS binary.Write failed:", err)
		}
	}
	_, err := control.ControlConn.Write(buf.Bytes())
	if err != nil {
		log.Fatal("send TypeCOLLAPSEPANELS error! ", err.Error())
		return
	}
}

func (control *ControlSender) GetClipboard() string {
	control.Lock.Lock()
	defer control.Lock.Unlock()
	// TODO 清理之前的数据
	//var buf bytes.Buffer
	//_, err := io.Copy(&buf, control.ControlConn)
	//if err != nil {
	//	// Error handler
	//}
	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint8(TypeGETCLIPBOARD), // base
	}
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			log.Fatal("GetClipboard binary.Write failed:", err)
		}
	}
	_, err := control.ControlConn.Write(buf.Bytes())
	if err != nil {
		log.Fatal("send GetClipboard error! ", err.Error())
	}
	var code *uint8
	recvCodeBuf := make([]byte, 1)
	_, err = control.ControlConn.Read(recvCodeBuf)
	if err != nil {
		log.Println(err.Error())
	}
	readerCode := bytes.NewReader(recvCodeBuf)
	err = binary.Read(readerCode, binary.BigEndian, &code)
	if err != nil {
		log.Fatal("GetClipboard binary.Read failed:", err.Error())
	}
	if *code != uint8(0) {
		log.Fatal("GetClipboard binary.Read failed: code != 0", err.Error())
	}
	var length *uint32
	recvlengthBuf := make([]byte, 4)
	_, err = control.ControlConn.Read(recvlengthBuf)
	if err != nil {
		log.Println(err.Error())
	}
	readerLength := bytes.NewReader(recvlengthBuf)
	err = binary.Read(readerLength, binary.BigEndian, &length)
	if err != nil {
		log.Fatal("GetClipboard binary.Read failed:", err.Error())
	}
	lengthBuf := make([]byte, int(*length))
	_, err = control.ControlConn.Read(lengthBuf)
	if err != nil {
		log.Println(err.Error())
	}
	return string(lengthBuf)
}

func (control *ControlSender) SetClipBoard(text string, pasted bool) {
	control.Lock.Lock()
	defer control.Lock.Unlock()

	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint8(TypeSETCLIPBOARD),   // base
		pasted,                    // ? 1 byte
		uint32(len([]byte(text))), // i 4 byte
	}
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("SetClipBoard binary.Write failed:", err)
		}
	}
	msg := append(buf.Bytes(), []byte(text)...)
	_, err := control.ControlConn.Write(msg)
	if err != nil {
		log.Fatal("send SetClipBoard error! ", err.Error())
		return
	}

}

func (control *ControlSender) SetScreenPowerMode(mode int) {
	control.Lock.Lock()
	defer control.Lock.Unlock()

	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint8(TypeSETScreenPowerMode), // base
		uint8(mode),                   // b 4 byte
	}
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("SetScreenPowerMode binary.Write failed:", err)
		}
	}
	_, err := control.ControlConn.Write(buf.Bytes())
	if err != nil {
		log.Fatal("send SetScreenPowerMode error! ", err.Error())
		return
	}

}

func (control *ControlSender) RotateDevice() {
	control.Lock.Lock()
	defer control.Lock.Unlock()

	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint8(TypeROTATEDEVICE), // base
	}
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("RotateDevice binary.Write failed:", err)
		}
	}
	_, err := control.ControlConn.Write(buf.Bytes())
	if err != nil {
		log.Fatal("send RotateDevice error! ", err.Error())
		return
	}

}

func (control *ControlSender) Swipe(startX, startY, endX, endY, stepLength int, delay float32) {
	control.Lock.Lock()
	defer control.Lock.Unlock()
	control.Touch(startX, startY, ActionDown)
	nextX := startX
	nextY := startY

	if endX > control.W {
		endX = control.W
	}
	if endY > control.H {
		endY = control.H
	}
	decreaseX := false
	decreaseY := false

	if startX > endX {
		decreaseX = true
	}
	if startY > endY {
		decreaseY = true
	}

	for {
		if decreaseX {
			nextX -= stepLength
			if nextX < endX {
				nextX = endX
			}
		} else {
			nextX += stepLength
			if nextX > endX {
				nextX = endX
			}
		}
		if decreaseY {
			nextY -= stepLength
			if nextY < endY {
				nextY = endY
			}
		} else {
			nextY += stepLength
			if nextY > endY {
				nextY = endY
			}
		}
		control.Touch(nextX, nextY, ActionDown)
		if nextX == endX && nextY == endY {
			control.Touch(nextX, nextY, ActionUp)
			break
		}
		time.Sleep(time.Duration(delay) * time.Second)
	}
}
