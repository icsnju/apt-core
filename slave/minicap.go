package main

import (
	"apsaras/comm"
	"bufio"
	"bytes"
	"errors"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type Banner struct {
	Version       int
	Length        int
	Pid           int
	RealWidth     int
	RealHeight    int
	VirtualWidth  int
	VirtualHeight int
	Orientation   int
	Quirks        int
}

var isRunning bool = true

const MAX_CHUNK = 2048
const WS_PORT = ":9002"
const PORT_START = 1711
const PORT_END = 1720
const PORT_FREE = "free"
const SCREEN_SIZE = 500

type DevWS struct {
	id       string
	isNeeded bool
	ws       *websocket.Conn
	lock     *sync.Mutex
}

func CreateDevWS(id string) *DevWS {
	return &DevWS{id, false, new(websocket.Conn), new(sync.Mutex)}
}

func (this *DevWS) takeWS(ws *websocket.Conn) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.isNeeded && ws != this.ws {
		this.ws.Close()
	}
	this.isNeeded = true
	this.ws = ws
}
func (this *DevWS) freeWS(ws *websocket.Conn) {
	this.lock.Lock()
	defer this.lock.Unlock()
	ws.Close()
	if this.ws == ws {
		this.isNeeded = false
	}
}

//send screen image frame
func (this *DevWS) sendFrame(frame []byte) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if len(frame) > 0 && this.isNeeded {
		err := websocket.Message.Send(this.ws, frame)
		if err != nil {
			log.Println("Send frame error", err)
		}
	}
}

//send log
func (this *DevWS) sendLog(content []byte) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if len(content) > 0 && this.isNeeded {
		err := websocket.Message.Send(this.ws, string(content))
		if err != nil {
			log.Println("Send log error", err)
		}
	}
}

type MiniPortManager struct {
	portMap map[int]string
	wsMap   map[string]*DevWS
	lock    *sync.Mutex
}

var portManager *MiniPortManager = GetMiniPortManager()

func GetMiniPortManager() *MiniPortManager {
	portMap := make(map[int]string)
	for i := PORT_START; i <= PORT_END; i++ {
		portMap[i] = PORT_FREE
	}
	wsMap := make(map[string]*DevWS)
	return &MiniPortManager{portMap, wsMap, new(sync.Mutex)}
}

func (this *MiniPortManager) getPort(id string) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	for mport, mid := range this.portMap {
		if mid == id {
			return mport
		}
	}
	return -1
}

func (this *MiniPortManager) allocatePort(id string) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	for mport, mid := range this.portMap {
		if mid == PORT_FREE {
			this.portMap[mport] = id
			return mport
		}
	}
	return -1
}

func (this *MiniPortManager) freePort(id string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for mport, mid := range this.portMap {
		if mid == id {
			this.portMap[mport] = PORT_FREE
		}
	}
}

func (this *MiniPortManager) addDevWS(id string, devws *DevWS) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.wsMap[id] = devws
}

func (this *MiniPortManager) getDevWS(id string) (*DevWS, bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	devWS, ex := this.wsMap[id]
	return devWS, ex
}

func (this *MiniPortManager) deleteDevWS(id string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	_, ex := this.wsMap[id]
	if ex {
		delete(this.wsMap, id)
	}
}

//start http server when the slave start
func startWebSocket() {
	//http.Handle("/", http.FileServer(http.Dir(".")))
	if err := http.ListenAndServe(WS_PORT, nil); err != nil {
		log.Println("ListenAndServe:", err)
		return
	}
	//http.Handle("/*", websocket.Handler(clientHandler))
}

//when a device is miss
func stopMinicap(id string) {
	portManager.freePort(id)
	portManager.deleteDevWS(id)
}

//when a new device is connected, start minicap server in the device and read the images data
func startMinicap(id, resolution string) {
	port := portManager.allocatePort(id)
	if port == -1 {
		log.Println("Port is not enough for", id)
		return
	}
	portManager.addDevWS(id, CreateDevWS(id))
	defer stopMinicap(id)

	//forward port
	ps := strconv.Itoa(port)
	out := comm.ExeCmd(getADBPath() + " -s " + id + " forward tcp:" + ps + " localabstract:minicap")
	if len(out) > 0 {
		log.Println(out)
		return
	}

	//regist this device in websocket server
	registDeviceInWS(id)

	//run minicap in the device
	cmd, err := runMCinDeviceCmd(id, resolution)
	if err != nil {
		log.Println("Minicap cmd create err", err)
		return
	}
	cmd.Start()
	defer func() {
		err = cmd.Process.Kill()
		if err != nil {
			log.Println("process kill err", err)
		}
	}()
	time.Sleep(2 * time.Second)

	//start get log and sent the log
	startLogcat(id)

	//start to send the image buffer
	log.Println("Start minicap on", port, " for ", id)
	parserImage(id, strconv.Itoa(port))
}

//get command for running minicap in device
func runMCinDeviceCmd(id, resolution string) (*exec.Cmd, error) {
	//defer portManager.freePort(id)
	var command *exec.Cmd
	var err error
	wh := strings.Split(resolution, "x")
	if len(wh) != 2 {
		return command, errors.New("resolution err")
	}
	w, err := strconv.Atoi(wh[0])
	if err != nil {
		return command, errors.New("resolution err")
	}
	h, err := strconv.Atoi(wh[1])
	if err != nil {
		return command, errors.New("resolution err")
	}

	sw := w * SCREEN_SIZE / h
	args := resolution + "@" + strconv.Itoa(sw) + "x" + strconv.Itoa(SCREEN_SIZE) + "/0"
	cmd := comm.CreateCmd("./minicap.sh " + args + " " + id)
	return cmd, nil
}

//run minicap in device
func runMCinDevice(id, resolution string) {
	//defer portManager.freePort(id)
	wh := strings.Split(resolution, "x")
	if len(wh) != 2 {
		return
	}
	w, err := strconv.Atoi(wh[0])
	if err != nil {
		return
	}
	h, err := strconv.Atoi(wh[1])
	if err != nil {
		return
	}

	sw := w * SCREEN_SIZE / h
	args := resolution + "@" + strconv.Itoa(sw) + "x" + strconv.Itoa(SCREEN_SIZE) + "/0"
	out := comm.ExeCmd("./minicap.sh " + args + " " + id)
	log.Println("minicap: ", out)
	portManager.freePort(id)
}

//regist the device in server
func registDeviceInWS(id string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err) //TODO panic: http: multiple registrations for
		}
	}()
	http.Handle("/"+id, websocket.Handler(clientHandler))
	//TODO cannot delete this handler now
}

//when a new client connet to
func clientHandler(ws *websocket.Conn) {
	id := ws.Request().URL.Path
	id = strings.TrimPrefix(id, "/")
	log.Println("new client connect to", id)

	dev, err := deviceManager.getDevice(id)
	if err != nil {
		log.Println("This device is not connected", id)
		return
	}

	devWS, ex := portManager.getDevWS(id)
	if !ex {
		log.Println("Device websocket information dosenot exist", id)
		return
	}
	//take up this device
	devWS.takeWS(ws)
	defer devWS.freeWS(ws)

	getEvent(ws, id, dev.Info.Resolution)
}

//get UI event from client
func getEvent(ws *websocket.Conn, id, resolution string) {
	wh := strings.Split(resolution, "x")
	if len(wh) != 2 {
		log.Println("resolution err: ", resolution)
		return
	}

	h, err := strconv.Atoi(wh[1])
	if err != nil {
		log.Println("resolution h: ", wh[1])
		return
	}

	zoom := float64(h) / float64(SCREEN_SIZE)
	for {
		var evtString string
		if err := websocket.Message.Receive(ws, &evtString); err != nil {
			log.Println("Get event err", err)
			break
		}
		xy := strings.Split(evtString, ",")
		if len(xy) == 2 {
			x, err := strconv.Atoi(xy[0])
			if err != nil {
				continue
			}
			y, err := strconv.Atoi(xy[1])
			if err != nil {
				continue
			}
			x = int(zoom * float64(x))
			y = int(zoom * float64(y))
			cmd := getADBPath() + " -s " + id + " shell input tap " + strconv.Itoa(x) + " " + strconv.Itoa(y)
			comm.ExeCmd(cmd)
		} else if len(xy) == 4 {
			x1, err := strconv.Atoi(xy[0])
			if err != nil {
				continue
			}
			y1, err := strconv.Atoi(xy[1])
			if err != nil {
				continue
			}
			x2, err := strconv.Atoi(xy[2])
			if err != nil {
				continue
			}
			y2, err := strconv.Atoi(xy[3])
			if err != nil {
				continue
			}
			x1 = int(zoom * float64(x1))
			y1 = int(zoom * float64(y1))
			x2 = int(zoom * float64(x2))
			y2 = int(zoom * float64(y2))
			cmd := getADBPath() + " -s " + id + " shell input swipe " + strconv.Itoa(x1) + " " + strconv.Itoa(y1) + " " + strconv.Itoa(x2) + " " + strconv.Itoa(y2)
			comm.ExeCmd(cmd)
		} else if len(xy) == 1 {
			cmd := getADBPath() + " -s " + id + " shell input keyevent "
			isKey := true
			switch xy[0] {
			case "back":
				cmd += "4"
			case "home":
				cmd += "3"
			case "menu":
				cmd += "1"
			case "power":
				cmd += "26"
			case "reboot":
				cmd = getADBPath() + " -s " + id + " reboot "
			default:
				isKey = false
				log.Println("Get event", xy[0])
			}
			if isKey {
				comm.ExeCmd(cmd)
			}
		} else {
			log.Println("Get event", xy)
		}

	}
}

//start logcat and send the content to websocket
func startLogcat(id string) {
	devWS, ex := portManager.getDevWS(id)
	if !ex {
		log.Println("device websocket not exist")
		return
	}

	cmd := comm.CreateCmd(getADBPath() + " -s " + id + " logcat -v time")

	// Create stdout, stderr streams of type io.Reader
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}
	// Start command
	err = cmd.Start()
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		read := bufio.NewReader(stdout)
		for {
			content, _, err := read.ReadLine()
			if err != nil {
				log.Println(err)
				break
			}
			if len(content) > 0 {
				devWS.sendLog(content)
			}
		}
	}()
}

//send images to client
func parserImage(id, port string) {

	//connect to device minicap
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:"+port)
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Println("conn close err", err)
		}
	}()
	reader := bufio.NewReader(conn)
	//get device websocket information
	devWS, ex := portManager.getDevWS(id)
	if !ex {
		log.Println("device websocket not exist")
		return
	}

	readBannerBytes := 0
	bannerLength := 2
	readFrameBytes := 0
	frameBodyLength := 0
	frameBody := bytes.NewBuffer([]byte{})
	var banner Banner = Banner{0, 0, 0, 0, 0, 0, 0, 0, 0}

	for isRunning {
		chunkLen := 0
		var chunk []byte = make([]byte, MAX_CHUNK)
		var err error
		for chunkLen <= 0 {
			chunkLen, err = reader.Read(chunk)
			if err != nil {
				log.Println("data read: ", err)
				return
			}
		}

		for cursor := 0; cursor < chunkLen; {
			data := int(chunk[cursor])
			if readBannerBytes < bannerLength {
				//read banner from data
				switch readBannerBytes {
				case 0:
					banner.Version = data
				case 1:
					banner.Length = data
					bannerLength = data
				case 2:
					fallthrough
				case 3:
					fallthrough
				case 4:
					fallthrough
				case 5:
					banner.Pid += data << uint((readBannerBytes-2)*8)
				case 6:
					fallthrough
				case 7:
					fallthrough
				case 8:
					fallthrough
				case 9:
					banner.RealWidth += data << uint((readBannerBytes-6)*8)
				case 10:
					fallthrough
				case 11:
					fallthrough
				case 12:
					fallthrough
				case 13:
					banner.RealHeight += data << uint((readBannerBytes-10)*8)
				case 14:
					fallthrough
				case 15:
					fallthrough
				case 16:
					fallthrough
				case 17:
					banner.VirtualWidth += data << uint((readBannerBytes-14)*8)
				case 18:
					fallthrough
				case 19:
					fallthrough
				case 20:
					fallthrough
				case 21:
					banner.VirtualHeight += data << uint((readBannerBytes-18)*8)
				case 22:
					banner.Orientation += data * 90
				case 23:
					banner.Quirks = data
				}
				cursor += 1
				readBannerBytes += 1
				if readBannerBytes == bannerLength {
					log.Println(banner)
				}
			} else if readFrameBytes < 4 {
				//read 4 bytes (the length of frame) in the head
				frameBodyLength += data << uint(readFrameBytes*8)
				cursor += 1
				readFrameBytes += 1
				//log.Println("headerbyte:", readFrameBytes, frameBodyLength, data)
			} else {
				if chunkLen-cursor >= frameBodyLength {
					le, err := frameBody.Write(chunk[cursor : cursor+frameBodyLength])

					if le != frameBodyLength || err != nil {
						log.Println("Frame body does not start with JPG header", frameBody, err)
						return
					}

					//log.Println("Get a frame len=", frameBody.Len())
					//send frame bytes to websocket
					devWS.sendFrame(frameBody.Bytes())

					cursor += frameBodyLength
					frameBodyLength = 0
					readFrameBytes = 0
					frameBody.Reset()
				} else {
					le, err := frameBody.Write(chunk[cursor:chunkLen])
					if le != chunkLen-cursor || err != nil {
						log.Println("Append frame err ", err)
						return
					}
					frameBodyLength -= chunkLen - cursor
					readFrameBytes += chunkLen - cursor
					cursor = chunkLen
				}
			}
		}

	}

}
