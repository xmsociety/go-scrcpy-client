package scrcpy

import "C"
import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/xmsociety/adbutils"
	"go-scrcpy-client/scrcpy/h264"
	"gocv.io/x/gocv"
	"image"
	"io"
	"log"
	"net"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
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
	VideoSender           chan image.Image
	ErrReceiver           chan error
	Resolution            resolution
	Control               ControlSender
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

func (client *Client) deployServer() bool {
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
	serverStream := client.Device.Shell(strings.Join(cmd, " "), true, 3)
	client.serverStream = *serverStream.(*adbutils.AdbConnection)
	//res := client.serverStream.ReadString(100) wait serverStream ready
	time.Sleep(time.Millisecond * 400)
	return true
}

func (client *Client) initServerConnection() bool {
	if client.ConnectionTimeout == 0 {
		client.ConnectionTimeout = 3000
	}
	for i := 0; i < client.ConnectionTimeout; i += 100 {
		client.videoSocket = client.Device.CreateConnection(adbutils.LOCALABSTRACT, "scrcpy")
		if client.videoSocket != nil {
			break
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}
	if client.videoSocket == nil {
		client.ErrReceiver <- errors.New("Failed to connect scrcpy-server after 3 seconds")
		log.Println("Failed to connect scrcpy-server after 3 seconds")
		return false
	}
	buf := readFully(client.videoSocket, 1)
	if buf == nil || len(buf) == 0 || buf[0] != []byte("\x00")[0] {
		client.ErrReceiver <- errors.New("Did not receive Dummy Byte!")
		log.Println("Did not receive Dummy Byte!")
		return false
	}
	client.controlSocket = client.Device.CreateConnection(adbutils.LOCALABSTRACT, "scrcpy")
	nameBuf := readFully(client.videoSocket, 64)
	if nameBuf == nil || len(nameBuf) == 0 || strings.TrimSuffix(string(nameBuf), "\x00") == "" {
		client.ErrReceiver <- errors.New(fmt.Sprintf("Did not receive Device Name! err: %v", nameBuf))
		log.Println("Did not receive Device Name! err: ", nameBuf)
		return false
	}
	resBuf := readFully(client.videoSocket, 4)
	r := bytes.NewReader(resBuf)

	resolutionTmp := resolution{}
	if err := binary.Read(r, binary.BigEndian, &resolutionTmp); err != nil {
		client.ErrReceiver <- errors.New(fmt.Sprintf("binary.Read failed:: %v", err))
		fmt.Println("binary.Read failed:", err)
		return false
	}
	client.Resolution = resolutionTmp
	client.Control = ControlSender{
		ControlConn: client.controlSocket,
		W:           int(client.Resolution.W),
		H:           int(client.Resolution.H),
		Lock:        sync.Mutex{},
	}
	log.Println(client.Resolution)
	return true
}

func (client *Client) Start() {
	client.Alive = true
	if !client.deployServer() {
		return
	}
	if !client.initServerConnection() {
		return
	}
	client.streamLoop()

}

func (client *Client) Stop() {
	client.Alive = false
	if &client.serverStream != nil {
		client.serverStream.Close()
	}
	if client.controlSocket != nil {
		client.controlSocket.Close()
	}
	if client.videoSocket != nil {
		client.videoSocket.Close()
	}
}

func (client *Client) streamLoop() {
	// thanks https://github.com/mike1808/h264decoder totaly what I want
	codec, err := h264.NewDecoder(h264.PixelFormatBGR)
	if err != nil {
		client.ErrReceiver <- err
		log.Println(err.Error())
		return
	}
	// noqa 看起来不错
	// TODO need fix -> could not determine kind of name for C.sws_addVec
	for client.Alive {
		buf := readFully(client.videoSocket, 0x10000)
		frames, err := codec.Decode(buf)
		if err != nil {
			client.ErrReceiver <- err
			log.Println(err.Error())
			return

		}
		if len(frames) == 0 {
			log.Println("no frames")
			break
		} else {
			for _, frame := range frames {
				client.Resolution = resolution{W: uint16(frame.Width), H: uint16(frame.Height)}
				client.Control.W = frame.Width
				client.Control.H = frame.Height
				imgCv, _ := gocv.NewMatFromBytes(frame.Height, frame.Width, gocv.MatTypeCV8UC3, frame.Data)
				if imgCv.Empty() {
					log.Println("empty")
					continue
				}
				imageRGB, err := imgCv.ToImage()
				if err != nil {
					log.Println(err.Error())
					continue
				}
				client.VideoSender <- imageRGB
			}
			log.Printf(fmt.Sprintf("found %d frames", len(frames)))
			time.Sleep(time.Microsecond * 100)
		}
	}
}
