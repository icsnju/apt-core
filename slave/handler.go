package main

import (
	"apsaras/comm"
	"apsaras/comm/comp"
	"encoding/gob"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

//Commucate with master
func diaMaster(serverIP string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverIP)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	mIP := conn.LocalAddr().String()
	mIP = strings.Split(mIP, ":")[0]

	//say hi
	_, err = conn.Write([]byte(comm.HIMASTER))
	if err != nil {
		log.Fatalln(err)
	}

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)
	//send heart beat
	for {
		var beat comp.SlaveInfo
		beat.IP = mIP

		//copy deivces information
		beat.DeviceStates = deviceManager.getDeviceInfo()

		//record all task state
		beat.TaskStates = taskManager.getTaskInfo()

		log.Println(beat.IP, " send beat to master!")
		//send beat
		err = encoder.Encode(&beat)
		if err != nil {
			log.Fatalln(err)
		}

		//get response
		conn.SetReadDeadline(time.Now().Add(comm.WAITFORDIA))
		var newTasks comp.RunTaskList
		err = decoder.Decode(&newTasks)
		if err != nil {
			log.Fatalln(err)
		}

		//allocate devices to these tasks
		allocateTasks(newTasks.Tasks)

		time.Sleep(comm.HEARTTIME)
	}
}

func allocateTasks(tasks []comp.RunTask) {
	if len(tasks) <= 0 {
		return
	}

	for _, ts := range tasks {
		//get the device
		ok := deviceManager.giveDevice(ts.TaskInfo.DeviceId)
		if ok {
			err := createLocalFile(ts)
			if err != nil {
				ok = false
			} else {
				//run this comp
				log.Println("Start task: ", ts.TaskInfo.JobId, ":", ts.TaskInfo.DeviceId)
				go runTask(ts)
				ts.TaskInfo.State = comp.TASK_RUN
			}
		}

		if !ok {
			ts.TaskInfo.State = comp.TASK_FAIL
			deviceManager.reclaim(ts.TaskInfo.DeviceId)
		}
		taskManager.addTask(ts)
	}
}

func createLocalFile(ts comp.RunTask) error {
	//create local file
	ex, err := comm.FileExists(ts.TaskInfo.JobId)
	if err != nil {
		return err
	}

	if !ex {
		//create local file
		err = os.Mkdir(ts.TaskInfo.JobId, os.ModePerm)
		if err != nil {
			return err
		}
		//make test locally
		ts.Frame = ts.Frame.MoveTestFile(ts.TaskInfo.JobId)
	}

	return nil
}

////Allot devices to tasks
//func allotDevice() {
//	taskLock.Lock()
//	for key, ts := range taskMap {
//		if ts.TaskInfo.State == task.TASK_QUEUE {
//			devLock.Lock()
//			dev, exist := deviceMap[ts.TaskInfo.DeviceId]
//			if exist && dev.State == comp.DEVICE_FREE {
//				//run this task
//				fmt.Println("Start task: " + ts.TaskInfo.JobId + "--" + ts.TaskInfo.DeviceId)
//				dev.State = comp.DEVICE_RUN
//				deviceMap[ts.TaskInfo.DeviceId] = dev
//				go runTask(ts)
//				ts.TaskInfo.State = task.TASK_RUN
//				taskMap[key] = ts
//			} else if !exist {
//				ts.TaskInfo.State = task.TASK_FAIL
//				taskMap[key] = ts
//			}
//			devLock.Unlock()
//		}
//	}
//	taskLock.Unlock()
//}

//Run a test task
func runTask(ts comp.RunTask) {
	jobId := ts.TaskInfo.JobId
	devId := ts.TaskInfo.DeviceId

	taskOutPath := path.Join(jobId, devId)
	err := os.Mkdir(taskOutPath, os.ModePerm)
	if err != nil {
		log.Println("Cannot create out file!", err)
		ts.TaskInfo.FinishTime = time.Now()
		ts.TaskInfo.State = comp.TASK_FAIL
	} else {
		//start task
		ts.Frame.TaskExecutor(jobId, devId)
		ts.TaskInfo.FinishTime = time.Now()
		ts.TaskInfo.State = comp.TASK_COMPLETE
	}
	taskManager.updateTaskStates(ts)
	//Update device information
	deviceManager.reclaim(devId)

}
