package test

import (
	"fmt"
	"github.com/xmsociety/adbutils"
	"go-scrcpy-client/scrcpy"
	"testing"
)

//var adb = adbutils.AdbClient{Host: "localhost", Port: 5037, SocketTime: 10}

func TestEquaBuf(t *testing.T) {
	fmt.Println([]byte("\x00")[0] == []byte("\x00")[0])
}

func TestConnect(t *testing.T) {
	adb := adbutils.AdbClient{Host: "localhost", Port: 5037, SocketTime: 10}
	snNtid := adbutils.SerialNTransportID{
		Serial: "192.168.0.107:5555",
	}
	//adb.Connect("emulator-5554")
	fmt.Println(adb.Device(snNtid).SayHello())
	client := scrcpy.Client{Device: adb.Device(snNtid), MaxWith: 800, Bitrate: 5000000}
	client.Start()
}
