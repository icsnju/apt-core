package main

import (
	"apsaras/comm"
	"bufio"
	"bytes"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

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

const MAX_CHUNK = 1024
const WS_PORT = ":9002"
const PORT_START = 1711
const PORT_END = 1720
const PORT_FREE = "free"

type MiniPortManager struct {
	portMap map[int]string
	lock    *sync.Mutex
}

var portManager *MiniPortManager = GetMiniPortManager()

func GetMiniPortManager() *MiniPortManager {
	portMap := make(map[int]string)
	for i := PORT_START; i <= PORT_END; i++ {
		portMap[i] = PORT_FREE
	}
	return &MiniPortManager{portMap, new(sync.Mutex)}
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

//start http server when the slave start
func startWebSocket() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	if err := http.ListenAndServe(WS_PORT, nil); err != nil {
		log.Println("ListenAndServe:", err)
		return
	}
}

//when a device is miss
func stopMinicap(id string) {
	portManager.freePort(id)
}

//when a new device is connected
func startMinicap(id string) {
	port := portManager.allocatePort(id)
	if port == -1 {
		log.Println("Port is not enough for", id)
		return
	}
	ps := strconv.Itoa(port)
	out := comm.ExeCmd(getADBPath() + " -s " + id + " forward tcp:" + ps + " localabstract:minicap")
	if len(out) > 0 {
		log.Println(out)
		portManager.freePort(id)
		return
	}
	log.Println("Start minicap on", port, " for ", id)
	//regist this device in websocket server
	registDeviceInWS(id)
	//run minicap in the device
	runMCinDevice(id)
}

//run minicap in device
func runMCinDevice(id string) {
	out := comm.ExeCmd("./minicap.sh autosize " + id)
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
	defer ws.Close()
	id := ws.Request().URL.Path
	id = strings.TrimPrefix(id, "/")
	log.Println("new client connect to", id)

	port := portManager.getPort(id)
	if port < 0 {
		log.Println("Device is not connected", id)
		return
	}
	ps := strconv.Itoa(port)
	//connect to device
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:"+ps)
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()
	reader := bufio.NewReader(conn)

	//start send bytes to websocket
	sendImage(ws, reader)

}

func sendImage(ws *websocket.Conn, reader *bufio.Reader) {
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

					err = websocket.Message.Send(ws, frameBody.Bytes())
					if err != nil {
						log.Println("Send frame error", le, err)
						return
					}

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
