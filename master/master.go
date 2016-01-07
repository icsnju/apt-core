package main

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand"
	"nata/andevice"
	"nata/framework"
	"nata/node"
	"nata/task"
	"nata/tools"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

const JOBPATH = "job"
const JSON = ".txt"

// slave map  ip: info
var slavesMap map[string]node.SlaveInfo
var jobMap map[string]task.Job

var shareDirPath string
var jobid int

var slavesLock *sync.Mutex
var jobLock *sync.Mutex
var idLock *sync.Mutex

func main() {

	//init files
	var err error
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

	cf.Close()

	service := ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	tools.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	defer listener.Close()
	tools.CheckError(err)

	//create job file
	err = os.Mkdir(JOBPATH, os.ModePerm)
	if err != nil {
		fmt.Println("Cannot create job file!")
		fmt.Println(err)
		os.Exit(1)
	}

	//register in gob
	framework.RigisterGob()
	//start find finished job
	go updateJobState()
	go shift()
	//TODO

	//start to wait for connect
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
			fmt.Println("Connet miss: ")
			fmt.Println(err)
			continue
		}

		mIP := conn.RemoteAddr().String()
		mIP = strings.Split(mIP, ":")[0]
		me := string(message[:mLen])
		if me == tools.HIMASTER {
			fmt.Println("Slave is ready!")
			//write slaves info
			slavesLock.Lock()
			_, ok := slavesMap[mIP]
			if ok {
				//this slave is in the list
				fmt.Println("This slave has connet to me : " + mIP)
				slavesLock.Unlock()
				return
			} else {
				slavesMap[mIP] = node.SlaveInfo{mIP, make(map[string]andevice.Device), make(map[string]task.Task)}
			}
			slavesLock.Unlock()
			//it is a slave, handle it
			go handleSlave(conn)
		} else if me == tools.CHECKJOBS || me == tools.CHECKSLAVES {
			go handleClient(conn, me)
		} else if me == tools.SUBJOB {
			go handleSubJob(conn)
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
}

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
func handleSubJob(conn net.Conn) {
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

	//read job information
	conn.SetReadDeadline(time.Now().Add(tools.WAITFORDIA)) // set 2 minutes timeout
	decoder := gob.NewDecoder(conn)
	var subjob task.SubJob
	err = decoder.Decode(&subjob)
	if err != nil {
		fmt.Println(err)
		return
	}

	var job task.Job
	job.JobId = mid
	job.JobInfo = subjob
	job.StartTime = time.Now()
	job.LatestTime = time.Now()

	//get current device list
	var devList []string
	slavesLock.Lock()
	var devMap map[string]andevice.Device = make(map[string]andevice.Device)
	for _, si := range slavesMap {
		for k, v := range si.DeviceStates {
			devMap[k] = v
		}
	}

	//get this job
	devList = job.JobInfo.Filter.GetDeviceSet(devMap)
	slavesLock.Unlock()

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
		jobLock.Lock()
		var jm task.JobMap
		jm.Map = jobMap
		err := encoder.Encode(jm)
		if err != nil {
			fmt.Println("Check jobs error!")
			fmt.Println(err)
		}
		jobLock.Unlock()
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
			if ds.State == andevice.DEVICE_FREE {
				//find a good job
				jobLock.Lock()

				bestJobId := findBestJob(id)
				//bestJobId := findOldJob(id)
				//TODO

				if bestJobId != "-1" {
					//udpate slave information
					ds.State = andevice.DEVICE_RUN
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

//find the oldest job
func findOldJob(id string) string {

	var maxJobId string = "-1"
	oldTime := time.Now()

	//find a old job
	for ke, jo := range jobMap {
		if jo.Finished {
			continue
		}
		ts, ex := jo.TaskMap[id]
		if ex && ts.State == task.TASK_WAIT {
			if jo.StartTime.Before(oldTime) {
				maxJobId = ke
				oldTime = jo.StartTime
			}
		}
	}
	return maxJobId
}

//find the best job
func findBestJob(id string) string {
	var maxPriority float64 = -1
	var count1 int = 0
	var maxJob1 string = "-1"

	var maxWaitTime float64 = -1
	var count2 int = 0
	var maxJob2 string = "-1"

	//find a good job
	for ke, jo := range jobMap {
		if jo.Finished {
			continue
		}
		ts, ex := jo.TaskMap[id]
		if ex && ts.State == task.TASK_WAIT {
			p := jo.GetPriority()
			if p == -1 {
				if jo.GetWaitTime() > maxWaitTime {
					maxWaitTime = jo.GetWaitTime()
					maxJob2 = ke
				}
				count2++
			} else {
				if p > maxPriority {
					maxPriority = p
					maxJob1 = ke
				}
				count1++
			}
		}
	}

	sum := count1 + count2
	if sum <= 0 {
		return "-1"
	}

	rand.Seed(time.Now().UnixNano())
	lottery := rand.Intn(count1 + count2)
	if lottery < count1 {
		return maxJob1
	} else {
		return maxJob2
	}
}

//Do something before close slave
func closeSlave(ip string) {
	//delete slave
	slavesLock.Lock()
	delete(slavesMap, ip)
	slavesLock.Unlock()
}

//Shift it
func shift() {
	for {
		jobLock.Lock()
		for jid, job := range jobMap {
			if !job.Finished && job.GetWaitTime() > task.MAXWAIT.Seconds() {
				slavesLock.Lock()
				var devMap map[string]andevice.Device = make(map[string]andevice.Device)
				for _, si := range slavesMap {
					for k, v := range si.DeviceStates {
						devMap[k] = v
					}
				}

				var nTaskMap map[string]task.Task = make(map[string]task.Task)
				for tid, ts := range job.TaskMap {
					if ts.State == task.TASK_WAIT {
						//get similar device
						target := ts.TargetId
						idMap := andevice.FindSimilarDevs(target, devMap, 0)
						//TODO

						//remove repeated id
						var idList []string = make([]string, 0)
						for cid, _ := range idMap {
							_, ex := job.TaskMap[cid]
							if !ex || cid == target {
								idList = append(idList, cid)
							}
						}

						//get a redom id
						if len(idList) > 0 {
							//give me a replacement for this device
							rand.Seed(time.Now().UnixNano())
							r := rand.Intn(len(idList))
							ts.DeviceId = idList[r]
							tid = ts.DeviceId
						}
					}
					nTaskMap[tid] = ts
				}
				slavesLock.Unlock()
				job.TaskMap = nTaskMap
				jobMap[jid] = job
			}
		}
		jobLock.Unlock()
		time.Sleep(tools.SHIFTTIME)
	}
}

//Start find finished job
func updateJobState() {
	for {
		var fiJobs []string = make([]string, 0)
		jobLock.Lock()
		for jid, job := range jobMap {
			if job.Finished {
				continue
			}
			isFinished := true
			for _, ts := range job.TaskMap {
				if ts.State == task.TASK_COMPLETE || ts.State == task.TASK_FAIL {

				} else {
					isFinished = false
				}
			}
			if isFinished {
				job.FinishTime = time.Now()
				job.Finished = true
				jobMap[jid] = job

				fiJobs = append(fiJobs, jid)
			}
		}

		//remove finished jobs
		for _, id := range fiJobs {
			djob, ex := jobMap[id]
			if ex {
				jp := path.Join(JOBPATH, id+JSON)
				jobFile, err1 := os.Create(jp)
				if err1 != nil {
					fmt.Println("updateJobState err!")
					fmt.Println(err1)
					continue
				}

				var content string = ""
				var sumRun int64 = 0
				var sumAll int64 = 0
				start := djob.StartTime
				var mdevlist []string = make([]string, 0)
				for _, ts := range djob.TaskMap {
					sumRun += ts.FinishTime.Sub(ts.StartTime).Nanoseconds() / 1000000
					sumAll += ts.FinishTime.Sub(start).Nanoseconds() / 1000000
					mdevlist = append(mdevlist, ts.DeviceId)
				}
				countDev := strconv.Itoa(len(djob.TaskMap))
				jobl := int64(len(djob.TaskMap))
				sumRun = sumRun / jobl
				sumAll = sumAll / jobl

				sum1 := strconv.FormatInt(sumRun, 10)
				sum2 := strconv.FormatInt(sumAll, 10)

				content = sum1 + " " + sum2 + " " + countDev + "\n"
				_, err := jobFile.Write([]byte(content))

				if err != nil {
					fmt.Println("updateJobState err!")
					fmt.Println(err)
					continue
				}

				devcon, err := json.Marshal(mdevlist)
				if err != nil {
					fmt.Println("updateJobState err!")
					fmt.Println(err)
					continue
				}

				_, err = jobFile.Write(devcon)
				if err != nil {
					fmt.Println("updateJobState err!")
					fmt.Println(err)
					continue
				}

				fmt.Println("Job " + id + " finishd!")
				jobFile.Sync()
				jobFile.Close()

				delete(jobMap, id)
			}
		}
		jobLock.Unlock()
		time.Sleep(tools.HEARTTIME)
	}
}
