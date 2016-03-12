package main

import (
	"apsaras/device"
	"apsaras/framework"
	"apsaras/node"
	"apsaras/task"
	"apsaras/tools"
	"bufio"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

var serviceIP string
var sharePath string

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Check devices infomation : " + tools.CHECKJOBS + " or " + tools.CHECKSLAVES)
		fmt.Println("Submit job : " + tools.SUBJOB + "path/to/job.json")
		return
	}
	//read config file
	cf, err := os.Open("client.conf")
	tools.CheckError(err)

	reader := bufio.NewReaderSize(cf, 1024)

	//service addr
	line, _, err := reader.ReadLine()
	tools.CheckError(err)
	sublines := strings.Split(string(line), "=")
	if len(sublines) == 2 && sublines[0] == "master" {
		serviceIP = sublines[1]
		fmt.Println("service addr is: " + serviceIP)
	} else {
		fmt.Println("service error: " + string(line))
	}
	//share path
	line, _, err = reader.ReadLine()
	tools.CheckError(err)
	sublines = strings.Split(string(line), "=")
	if len(sublines) == 2 && sublines[0] == "share" {
		sharePath = sublines[1]
		fmt.Println("share path is: " + sharePath)
	} else {
		fmt.Println("share path wrong: " + string(line))
	}
	cf.Close()

	//connet master
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serviceIP)
	tools.CheckError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	tools.CheckError(err)
	defer conn.Close()

	//register in gob
	framework.RigisterGob()

	kind := os.Args[1]
	if kind == tools.CHECKJOBS || kind == tools.CHECKSLAVES {
		checkInfo(conn, kind)
	} else if kind == tools.JOBSTATE && len(os.Args) == 3 {
		checkJob(conn, os.Args[2])

	} else if kind == tools.SUBJOB && len(os.Args) == 3 {
		jspath := os.Args[2]
		subJob(conn, jspath)
	} else {
		fmt.Println("Check devices infomation : " + tools.CHECKJOBS + " or " + tools.CHECKSLAVES)
		fmt.Println("Sub job : " + tools.SUBJOB + " [App.apk]" + " [Test.apk]" + " deviceId[,deviceid]")
	}

}

//Check job
func checkJob(conn *net.TCPConn, id string) {
	_, err := conn.Write([]byte(tools.JOBSTATE + ":" + id))
	tools.CheckError(err)
	conn.SetReadDeadline(time.Now().Add(tools.WAITFORDIA)) // set 2 minutes timeout
	decoder := gob.NewDecoder(conn)
	var job task.Job
	decoder.Decode(&job)
	fmt.Println("job id:" + job.JobId)
	if job.JobId == "unknown" {
		return
	}

	for se, st := range job.TaskMap {
		var state string
		switch st.State {
		case task.TASK_COMPLETE:
			state = "complete"
		case task.TASK_RUN:
			state = "run"
		case task.TASK_WAIT:
			state = "wait"
		case task.TASK_FAIL:
			state = "fail"
		default:
			state = "unknown"
		}
		fmt.Println("\t device id:" + se + "\t" + state)
	}
}

//Sub job
func subJob(conn *net.TCPConn, jsPath string) {

	ex, err := tools.FileExists(jsPath)
	tools.CheckError(err)
	if !ex {
		fmt.Println("Error! Job json file dose not exist! ", jsPath)
		return
	}

	content, err := ioutil.ReadFile(jsPath)
	tools.CheckError(err)

	sj, err := task.ParserSubJobFromJson(content)
	if err != nil {
		fmt.Println(err)
		return
	}

	//say what do you want
	_, err = conn.Write([]byte(tools.SUBJOB))
	tools.CheckError(err)
	conn.SetReadDeadline(time.Now().Add(tools.WAITFORDIA)) // set 2 minutes timeout

	//get job id
	message := make([]byte, 128)
	mLen, err := conn.Read(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	jobid := string(message[:mLen])
	fmt.Println("Job id: " + jobid)

	//copy test file to share dir
	jobPath := path.Join(sharePath, jobid)
	os.RemoveAll(jobPath)
	os.Mkdir(jobPath, os.ModePerm)
	sj.Frame = sj.Frame.MoveTestFile(jobPath)

	//send this job
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(&sj)
	tools.CheckError(err)
}

//Commucate with master
func checkInfo(conn *net.TCPConn, kind string) {

	//say what do you want
	_, err := conn.Write([]byte(kind))
	tools.CheckError(err)

	conn.SetReadDeadline(time.Now().Add(tools.WAITFORDIA)) // set 2 minutes timeout

	decoder := gob.NewDecoder(conn)
	if kind == tools.CHECKSLAVES {
		var slaves node.SlaveMap
		err := decoder.Decode(&slaves)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Get slaves err! ")
			return
		}

		var devMap map[string]device.DeviceInfo = make(map[string]device.DeviceInfo)
		for _, sv := range slaves.Map {
			fmt.Println("ip: " + sv.IP)
			fmt.Println("devices:")
			for _, dev := range sv.DeviceStates {
				fmt.Println(dev.Info)
				fmt.Println(dev.State)
				_, ex := devMap[dev.Info.Id]
				if ex {
					fmt.Println("Error! Same device " + dev.Info.Id)
				} else {
					devMap[dev.Info.Id] = dev.Info
				}
			}

			//save devices infomation
			fpath := path.Join("devices.json")
			devf, err1 := os.Create(fpath)
			content, err2 := json.Marshal(devMap)
			if err1 != nil || err2 != nil {
				fmt.Println("save devices infomation err!")
				fmt.Println(err1)
				fmt.Println(err2)
				continue
			}
			_, err := devf.Write(content)
			if err != nil {
				fmt.Println("save devices infomation err!")
				fmt.Println(err)
				continue
			}
			devf.Sync()
			devf.Close()

			fmt.Println("tasks:")
			for _, ts := range sv.TaskStates {
				fmt.Println(ts)
			}
		}
	} else if kind == tools.CHECKJOBS {
		var jobs task.JobMap
		err := decoder.Decode(&jobs)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Get jobs err! ")
			return
		}
		for _, job := range jobs.Map {
			//fmt.Println(job)
			taskMap := job.TaskMap
			fmt.Println(job.JobId, ":")
			var countF int = 0
			var countW int = 0
			var countR int = 0
			var countE int = 0
			for _, ts := range taskMap {
				switch ts.State {
				case task.TASK_COMPLETE:
					countF++
				case task.TASK_RUN:
					countR++
					fmt.Println("run: ", ts.DeviceId)
				case task.TASK_WAIT:
					countW++
					fmt.Println("wait: ", ts.DeviceId)
				case task.TASK_FAIL:
					countE++
					fmt.Println("fail: ", ts.DeviceId)
				default:
				}
			}
			fmt.Printf("Fini: ", countF, " |Wait: ", countW, " |Run: ", countR, " |Fial: ", countE)
		}
	}

}
