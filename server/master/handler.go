package master

import (
	"apsaras/comm"
	"apsaras/comm/comp"
	"encoding/gob"
	"log"
	"net"
	"strings"
	"time"
)

//handle slave communication
func handleSlave(conn net.Conn) {

	mIP := conn.RemoteAddr().String()
	mIP = strings.Split(mIP, ":")[0]
	log.Println("New slave dial me: " + mIP)

	defer conn.Close() // close connection before exit
	defer closeSlave(mIP)

	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)
	for {
		//wait for your heat
		conn.SetReadDeadline(time.Now().Add(comm.WAITFORHEARTTIME)) // set 5 minutes timeout

		//get beat content from slave
		var getBeat comp.SlaveInfo
		err := decoder.Decode(&getBeat)
		if err != nil {
			log.Println("Slave is disconnected!", err)
			break
		}
		if getBeat.IP != mIP {
			log.Println("IP Error!")
			break
		}
		log.Println("Get heart beat from: " + getBeat.IP)

		//update slave information
		slaveManager.updateSlave(getBeat)

		//update job queue information
		jobManager.updateJobs(getBeat.TaskStates)

		//heart beat response
		taskList := getWaitTasks(mIP)
		err = encoder.Encode(&taskList)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

//Do something before close slave
func closeSlave(ip string) {
	//delete slave
	slaveManager.removeSlave(ip)
}

//Get priority tasks to run
func getWaitTasks(ip string) comp.RunTaskList {
	var taskList comp.RunTaskList
	taskList.Tasks = make([]comp.RunTask, 0)

	slave, ex := slaveManager.getSlave(ip)
	//find task to send
	if ex {
		for id, ds := range slave.DeviceStates {
			if ds.State == comp.DEVICE_FREE {

				rt, ex := jobManager.findBestJob(id) //TODO
				if ex {
					//udpate slave information
					slaveManager.updateSlaveDeviceState(ip, id, comp.DEVICE_RUN)
					taskList.Tasks = append(taskList.Tasks, rt)
					log.Println("Send task " + rt.TaskInfo.JobId + "--" + id + " to: " + ip)
				}
			}
		}
	}
	return taskList
}
