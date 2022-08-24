package scrcpy

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/xmsociety/adbutils"
	"io"
	"log"
	"net"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func getCurrentFile() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("Can not get current file info"))
	}
	return getParentDirectory(file)
}

type resolution struct {
	W uint16
	H uint16
}

type Client struct {
	Device                adbutils.AdbDevice
	MaxWith               int
	Bitrate               int //8000000
	MaxFps                int
	Alive                 bool
	Flip                  bool
	BlockFrame            bool
	StayAwake             bool
	LockScreenOrientation int
	ConnectionTimeout     int    //3000
	EncoderName           string // "-"
	serverStream          adbutils.AdbConnection
	videoSocket           net.Conn
	controlSocket         net.Conn
	VideoSender           chan<- interface{}
	resolution
}

func readFully(conn net.Conn, n int) []byte {
	t := 0
	buffer := make([]byte, n)
	result := bytes.NewBuffer(nil)
	for t < n {
		length, err := conn.Read(buffer[0:n])
		if length == 0 || err != nil {
			log.Println(err.Error())
			break
		}
		result.Write(buffer[0:length])
		t += length
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Println(err.Error())
			}
		}
	}
	return result.Bytes()
}

func (client *Client) deployServer() {
	jarName := "scrcpy-server.jar"
	currentPath := getCurrentFile()
	src, _ := filepath.Abs(path.Join(currentPath, jarName))
	client.Device.Push(src, fmt.Sprintf("/data/local/tmp/%v", jarName))
	stayAwake := "false"
	if client.StayAwake {
		stayAwake = "true"
	}
	// CLASSPATH=/data/local/tmp/scrcpy-server.jar   app_process / com.genymobile.scrcpy.Server 1.20  info 0 100000 0 -1 true - false ture 0 false true - - false
	cmd := []string{
		fmt.Sprintf("CLASSPATH=/data/local/tmp/%v", jarName),
		"app_process",
		"/",
		"com.genymobile.scrcpy.Server",
		"1.20",                       // Scrcpy server version
		"debug",                      // Log level: info, verbose...
		strconv.Itoa(client.MaxWith), // Max screen width (long side)
		strconv.Itoa(client.Bitrate), // Bitrate of video
		strconv.Itoa(client.MaxFps),  // Max frame per second
		strconv.Itoa(LockScreenOrientationUnlocked), // Lock screen orientation: LOCK_SCREEN_ORIENTATION
		"true",    // Tunnel forward
		"-",       // Crop screen
		"false",   // Send frame rate to client
		"true",    // Control enabled
		"0",       // Display id
		"false",   // Show touches
		stayAwake, // Stay awake
		"-",       // Codec (video encoding) options
		"-",       // Encoder name
		"false",   // Power off screen after server closed
	}
	serverStream := client.Device.Shell(strings.Join(cmd, " "), true)
	client.serverStream = *serverStream.(*adbutils.AdbConnection)
	res := client.serverStream.ReadString(100)
	log.Println("deploy server res: ", res)
}

func (client *Client) initServerConnection() {
	if client.ConnectionTimeout == 0 {
		client.ConnectionTimeout = 3000
	}
	for i := 0; i < client.ConnectionTimeout; i += 100 {
		client.videoSocket = client.Device.CreateConnection(adbutils.LOCALABSTRACT, "scrcpy")
		if client.videoSocket != nil {
			break
		}
	}
	if client.videoSocket == nil {
		log.Fatal("Failed to connect scrcpy-server after 3 seconds")
	}
	buf := readFully(client.videoSocket, 1)
	if buf == nil || len(buf) == 0 || buf[0] != []byte("\x00")[0] {
		log.Fatal("Did not receive Dummy Byte!")
	}
	client.controlSocket = client.Device.CreateConnection(adbutils.LOCALABSTRACT, "scrcpy")
	nameBuf := readFully(client.videoSocket, 64)
	if nameBuf == nil || len(nameBuf) == 0 || strings.TrimSuffix(string(nameBuf), "\x00") == "" {
		log.Fatal("Did not receive Device Name! err: ", nameBuf)
	}
	resBuf := readFully(client.videoSocket, 4)
	r := bytes.NewReader(resBuf)

	if err := binary.Read(r, binary.BigEndian, &client.resolution); err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	log.Println(client.resolution)
}

func (client *Client) Start() {
	client.deployServer()
	client.initServerConnection()
	client.Alive = true
	client.streamLoop()
}

func (client *Client) Stop() {
	client.Alive = false
	client.serverStream.Close()
	client.controlSocket.Close()
	client.videoSocket.Close()
}

func (client *Client) streamLoop() {
	// TODO decode h264
	//dec, err := NewH264Decoder(pps)
	for client.Alive {
		buf := readFully(client.videoSocket, 0x10000)
		//for i, n := range nal[1:] {
		//	img, err := dec.Decode(n)
		//	if err != nil {
		//		continue
		//	}
		fmt.Println(buf)
		if len(buf) == 0 {
			time.Sleep(time.Second * 1)
		}
		time.Sleep(time.Microsecond * 100)
		//client.VideoSender <- struct {
		//}{}
	}
}
