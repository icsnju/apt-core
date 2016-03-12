package main

import (
	"apsaras/device"
	"apsaras/framework"
	"apsaras/node"
	"apsaras/task"
	"apsaras/tools"
	"bufio"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

//
const DBURL = "localhost"

// slave map ip: info
var slavesMap map[string]node.SlaveInfo

// job map
var jobMap map[string]task.Job

var shareDirPath string
var jobid int

var slavesLock *sync.Mutex
var jobLock *sync.Mutex
var idLock *sync.Mutex

func main() {

	//init
	jobid = 0
	slavesMap = make(map[string]node.SlaveInfo)
	jobMap = make(map[string]task.Job)
	slavesLock = new(sync.Mutex)
	jobLock = new(sync.Mutex)
	idLock = new(sync.Mutex)

	//read config file
	cf, err := os.Open("master.conf")
	tools.CheckError(err)

	reader := bufio.NewReaderSize(cf, 1024)
	//share path
	line, _, err := reader.ReadLine()
	tools.CheckError(err)
	sublines := strings.Split(string(line), "=")
	if len(sublines) == 2 && sublines[0] == "share" {
		shareDirPath = sublines[1]
		os.RemoveAll(shareDirPath) //clean old files
		subPath := path.Join(shareDirPath, tools.MASTER)
		os.MkdirAll(subPath, os.ModePerm)
		subPath = path.Join(shareDirPath, tools.SLAVE)
		os.MkdirAll(subPath, os.ModePerm)
		fmt.Println("share file path: " + shareDirPath)
	} else {
		fmt.Println("share path error: " + string(line))
		os.Exit(1)
	}

	//port
	line, _, err = reader.ReadLine()
	tools.CheckError(err)
	sublines = strings.Split(string(line), "=")
	var port string
	if len(sublines) == 2 && sublines[0] == "port" {
		port = sublines[1]
		fmt.Println("port: " + port)
	} else {
		fmt.Println("port error: " + string(line))
		os.Exit(1)
	}

	//init DB
	err = initDB(DBURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer closeDB()

	//close config file
	cf.Close()

	//create TCP listener
	service := ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	tools.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	defer listener.Close()
	tools.CheckError(err)

	//register in gob
	framework.RigisterGob()
	//start find finished job
	go updateJobState()
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
		conn.SetReadDeadline(time.Now().Add(tools.WAITFORDIA)) // set 2 minutes timeout
		mLen, err := conn.Read(message)
		if err != nil {
			fmt.Print("Connet miss: ")
			fmt.Println(err)
			continue
		}

		mIP := conn.RemoteAddr().String()
		mIP = strings.Split(mIP, ":")[0]
		me := string(message[:mLen])
		if me == tools.HIMASTER {
			fmt.Println("A new slave connet to me!")
			//write slaves info
			slavesLock.Lock()
			_, ok := slavesMap[mIP]
			if ok {
				//this slave is in the list
				fmt.Println("This slave has connet to me : " + mIP)
			} else {
				slavesMap[mIP] = node.SlaveInfo{mIP, make(map[string]device.Device), make(map[string]task.Task)}
			}
			slavesLock.Unlock()
			//it is a slave, handle it
			go handleSlave(conn)
		} else if me == tools.CHECKJOBS || me == tools.CHECKSLAVES || me == tools.CHECKDEVICES {
			go handleClient(conn, me)
		} else if me == tools.SUBJOB || me == tools.WEBSUBJOB {
			go handleSubJob(conn, me)
		} else {
			if strings.HasPrefix(me, "job:") {
				terms := strings.Split(me, ":")
				if len(terms) == 2 {
					go handleJobQuery(conn, terms[1])
				}
			} else {
				fmt.Println("Message is wrong: " + me)
			}
		}
	}
	//never end
}
