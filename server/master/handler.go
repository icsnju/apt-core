package master

import (
	"apsaras/comm"
	"apsaras/comm/comp"
	"encoding/gob"
	"fmt"
	"net"
	"strings"
	"time"
)

//handle slave communication
func handleSlave(conn net.Conn) {

	defer conn.Close() // close connection before exit
	mIP := conn.RemoteAddr().String()
	mIP = strings.Split(mIP, ":")[0]
	fmt.Println("New slave dial me: " + mIP)

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
			fmt.Println(err)
			fmt.Println("Slave is dead! ")
			return
		}
		if getBeat.IP != mIP {
			fmt.Println("IP Error!")
			return
		}
		fmt.Println("Get heart beat from: " + getBeat.IP)

		//update slave information
		slaveManager.updateSlave(getBeat)

		//update job queue information
		jobManager.updateJobs(getBeat)

		//heart beat response
		taskList := getWaitTasks(mIP)
		err = encoder.Encode(&taskList)
		if err != nil {
			fmt.Println(err)
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

	slave := slaveManager.getSlave(ip)
	//find task to send
	for id, ds := range slave.DeviceStates {
		if ds.State == comp.DEVICE_FREE {

			bestJobId := jobManager.findBestJob(id) //TODO
			if bestJobId != "-1" {
				//udpate slave information
				slaveManager.updateSlaveDeviceState(ip, id, comp.DEVICE_RUN)

				//update job information
				jobManager.updateJobTaskState(bestJobId, id, comp.TASK_RUN)

				//create run task
				rt := jobManager.createRuntask(bestJobId, id)

				taskList.Tasks = append(taskList.Tasks, rt)
				fmt.Println("Send task " + bestJobId + "--" + id + " to: " + ip)
			}
		}
	}
	return taskList
}
