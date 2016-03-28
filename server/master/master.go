package master

import (
	"apsaras/comm"
	"apsaras/comm/comp"
	"log"
	"net"
	"strings"
	"time"
)

func StartMaster(port string) {
	//create TCP listener
	service := ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		log.Fatalln(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	//start find finished job
	go jobManager.updateJobInDB()
	//go shift() //TODO

	//start to wait for slave
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// set maxium request length to 128KB to prevent flood attack
		message := make([]byte, 128)
		//wait hi
		conn.SetReadDeadline(time.Now().Add(comm.WAITFORDIA))
		mLen, err := conn.Read(message)
		if err != nil {
			log.Println("Connet miss: ", err)
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
				log.Println("A new slave connet to me!")
				go handleSlave(conn)
			}
		}
	}
	//never end
}
