package master

import (
	"apsaras/comm"
	"apsaras/comm/comp"
	"apsaras/comm/filter"
	"apsaras/comm/framework"
	"bufio"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

const (
	CONFPATH = "conf/master.conf"
)

var shareDirPath string

func StartMaster() {

	//read config file
	cf, err := os.Open(CONFPATH)
	comm.CheckError(err)

	reader := bufio.NewReaderSize(cf, 1024)
	//share path
	line, _, err := reader.ReadLine()
	comm.CheckError(err)
	sublines := strings.Split(string(line), "=")
	if len(sublines) == 2 && sublines[0] == "share" {
		shareDirPath = sublines[1]
		os.RemoveAll(shareDirPath) //clean old files
		subPath := path.Join(shareDirPath, comm.MASTER)
		os.MkdirAll(subPath, os.ModePerm)
		subPath = path.Join(shareDirPath, comm.SLAVE)
		os.MkdirAll(subPath, os.ModePerm)
		fmt.Println("share file path: " + shareDirPath)
	} else {
		fmt.Println("share path error: " + string(line))
		os.Exit(1)
	}

	//port
	line, _, err = reader.ReadLine()
	comm.CheckError(err)
	sublines = strings.Split(string(line), "=")
	var port string
	if len(sublines) == 2 && sublines[0] == "port" {
		port = sublines[1]
		fmt.Println("port: " + port)
	} else {
		fmt.Println("port error: " + string(line))
		os.Exit(1)
	}

	//close config file
	cf.Close()

	//create TCP listener
	service := ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	comm.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	defer listener.Close()
	comm.CheckError(err)

	//register in gob
	framework.RigisterGob()
	filter.RigisterGob()
	//start find finished job
	go jobManager.updateJobInDB()
	//go shift() //TODO

	//start to wait for slave
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		message := make([]byte, 128) // set maxium request length to 128KB to prevent flood attack
		//wait hi
		conn.SetReadDeadline(time.Now().Add(comm.WAITFORDIA)) // set 2 minutes timeout
		mLen, err := conn.Read(message)
		if err != nil {
			fmt.Print("Connet miss: ")
			fmt.Println(err)
			continue
		}

		mIP := conn.RemoteAddr().String()
		mIP = strings.Split(mIP, ":")[0]
		me := string(message[:mLen])
		if me == comm.HIMASTER {

			slave := comp.SlaveInfo{mIP, make(map[string]comp.Device), make(map[string]comp.Task)}
			ok := slaveManager.addSlave(slave)
			//it is a slave, handle it
			if ok {
				fmt.Println("A new slave connet to me!")
				go handleSlave(conn)
			}
		}
	}
	//never end
}
