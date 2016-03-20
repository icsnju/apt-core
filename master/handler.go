package main

import (
	"apsaras/device"
	"apsaras/node"
	"apsaras/task"
	"apsaras/tools"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

//handle job query
func handleJobQuery(conn net.Conn, jobId string) {
	var job task.Job
	jobLock.Lock()
	job, ex := jobMap[jobId]
	jobLock.Unlock()
	encoder := gob.NewEncoder(conn)
	if !ex {
		job.JobId = "unknown"
	}
	err := encoder.Encode(&job)
	if err != nil {
		fmt.Println("handleJobQuery: send job error!")
		fmt.Println(err)
	}
}

//handle sub job
func handleSubJob(conn net.Conn, kind string) {
	defer conn.Close()
	idLock.Lock()
	mid := strconv.Itoa(jobid)
	respo := mid
	jobid++
	idLock.Unlock()
	//tell client jobid
	_, err := conn.Write([]byte(respo))
	if err != nil {
		fmt.Println(err)
	}

	var subjob task.SubJob
	conn.SetReadDeadline(time.Now().Add(tools.WAITFORDIA)) // set 2 minutes timeout
	if kind == tools.SUBJOB {
		//read job information
		decoder := gob.NewDecoder(conn)

		err = decoder.Decode(&subjob)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		message := make([]byte, 1024) // set maxium request length to 128KB to prevent flood attack
		mLen, err := conn.Read(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		me := message[:mLen]
		fmt.Println(string(me))
		subjob, err = task.ParserSubJobFromJson(me)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	var job task.Job
	job.JobId = mid
	job.JobInfo = subjob
	job.StartTime = time.Now()
	job.LatestTime = time.Now()

	//get current device list
	var devList []string
	slavesLock.Lock()
	var devMap map[string]device.Device = make(map[string]device.Device)
	for _, si := range slavesMap {
		for k, v := range si.DeviceStates {
			devMap[k] = v
		}
	}
	slavesLock.Unlock()
	//get this job
	devList = job.JobInfo.Filter.GetDeviceSet(devMap)
	var taskMap map[string]task.Task = make(map[string]task.Task)
	for _, devId := range devList {
		var t task.Task
		t.JobId = job.JobId
		t.DeviceId = devId
		t.TargetId = devId //勿忘初心
		t.State = task.TASK_WAIT
		taskMap[devId] = t
	}
	job.TaskMap = taskMap

	//add this job in map
	jobLock.Lock()
	_, ex := jobMap[job.JobId]
	if ex {
		fmt.Println("Error! Job id repetitive: " + job.JobId)
	} else {
		jobMap[job.JobId] = job
	}
	jobLock.Unlock()

}

//respose client request
func handleClient(conn net.Conn, kind string) {
	defer conn.Close()
	encoder := gob.NewEncoder(conn)
	if kind == tools.CHECKJOBS {
		var jobs []task.JobBrief = make([]task.JobBrief, 0)

		jobLock.Lock()
		for _, jbif := range jobMap {
			temp := task.BriefThisJob(jbif)
			jobs = append(jobs, temp)
		}
		jobLock.Unlock()

		jobjson, err := json.Marshal(jobs)
		if err != nil {
			fmt.Println("get jobs err!")
			fmt.Println(err)
			return
		}

		_, err = conn.Write([]byte(jobjson))
		if err != nil {
			fmt.Println(err)
		}

	} else if kind == tools.CHECKSLAVES {
		slavesLock.Lock()
		var sm node.SlaveMap
		sm.Map = slavesMap
		err := encoder.Encode(sm)
		if err != nil {
			fmt.Println("Check slaves error!")
			fmt.Println(err)
		}
		slavesLock.Unlock()
	} else if kind == tools.CHECKDEVICES {
		var devices []device.DeviceInfo = make([]device.DeviceInfo, 0)

		slavesLock.Lock()
		for _, slave := range slavesMap {
			for _, device := range slave.DeviceStates {
				devices = append(devices, device.Info)
			}
		}
		slavesLock.Unlock()

		devjson, err := json.Marshal(devices)
		if err != nil {
			fmt.Println("get devices err!")
			fmt.Println(err)
			return
		}

		_, err = conn.Write([]byte(devjson))
		if err != nil {
			fmt.Println(err)
		}
	}
}

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
		conn.SetReadDeadline(time.Now().Add(tools.WAITFORHEARTTIME)) // set 5 minutes timeout

		//get beat content from slave
		var getBeat node.SlaveInfo
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
		slavesLock.Lock()
		slavesMap[mIP] = getBeat
		slavesLock.Unlock()

		//update job queue information
		jobLock.Lock()
		taskinfo := getBeat.TaskStates
		for _, t := range taskinfo {
			jid := t.JobId
			did := t.DeviceId
			job, ex := jobMap[jid]
			if !ex {
				fmt.Println("Job not exist! ", jid)
			}
			_, ex = job.TaskMap[did]
			if !ex {
				fmt.Println("Device not exist! ", did)
			}
			job.TaskMap[did] = t
			jobMap[jid] = job
		}
		jobLock.Unlock()

		//heart beat response
		var taskList task.RunTaskList
		taskList.Tasks = make([]task.RunTask, 0)

		//find task to send
		slave, _ := slavesMap[mIP]
		for id, ds := range slave.DeviceStates {
			if ds.State == device.DEVICE_FREE {
				//find a good job
				jobLock.Lock()

				bestJobId := findBestJob(id)
				//bestJobId := findOldJob(id)
				//TODO

				if bestJobId != "-1" {
					//udpate slave information
					ds.State = device.DEVICE_RUN
					slave.DeviceStates[id] = ds
					slavesMap[mIP] = slave
					//update job information
					j := jobMap[bestJobId]
					ts := j.TaskMap[id]
					ts.State = task.TASK_RUN
					ts.StartTime = time.Now()
					j.TaskMap[id] = ts
					j.LatestTime = time.Now()
					jobMap[bestJobId] = j
					//create run task
					var rt task.RunTask
					rt.Frame = j.JobInfo.Frame
					rt.TaskInfo = j.TaskMap[id]

					taskList.Tasks = append(taskList.Tasks, rt)
					fmt.Println("Send task " + j.JobId + "--" + id + " to: " + mIP)
				}
				jobLock.Unlock()
			}
		}
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
	slavesLock.Lock()
	delete(slavesMap, ip)
	slavesLock.Unlock()
}

//Start find finished job
func updateJobState() {
	for {
		//update job status
		var finishedJobs []string = make([]string, 0)
		jobLock.Lock()
		for jid, job := range jobMap {

			isFinished := true
			for _, ts := range job.TaskMap {
				if ts.State == task.TASK_COMPLETE || ts.State == task.TASK_FAIL {

				} else {
					isFinished = false
				}
			}
			if isFinished {
				finishedJobs = append(finishedJobs, jid)
				job.FinishTime = time.Now()
				job.Finished = true
				jobMap[jid] = job

				//save this job in DB
				temp := task.BriefThisJob(job)
				saveJobInDB(temp)
				delete(jobMap, jid)
			}
		}
		jobLock.Unlock()

		//delete finished job
		jobLock.Lock()
		for _, id := range finishedJobs {
			delete(jobMap, id)
		}
		jobLock.Unlock()

		time.Sleep(tools.HEARTTIME)
	}
}
