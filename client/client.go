package main

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"nata/andevice"
	"nata/framework"
	"nata/node"
	"nata/task"
	"nata/tools"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
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
	js, err := simplejson.NewJson(content)
	tools.CheckError(err)
	framekind, err := js.Get("FrameKind").String()
	tools.CheckError(err)
	filterkind, err := js.Get("FilterKind").String()
	tools.CheckError(err)

	var sj task.SubJob

	switch framekind {
	case framework.FRAME_ROBOT:
		sj.FrameKind = framework.FRAME_ROBOT
		var rf framework.RobotFrame

		appPath, err1 := js.Get("Frame").Get("AppPath").String()
		testPath, err2 := js.Get("Frame").Get("TestPath").String()
		if err1 != nil || err2 != nil {
			fmt.Println("Robotium framework error in json file! File path of App and Test are needed!")
			return
		}
		rf.AppPath = appPath
		rf.TestPath = testPath
		sj.Frame = rf
	case framework.FRAME_MONKEY:
		sj.FrameKind = framework.FRAME_MONKEY
		var mf framework.MonkeyFrame

		appPath, err1 := js.Get("Frame").Get("AppPath").String()
		argu, err2 := js.Get("Frame").Get("Argu").String()
		pkg, err3 := js.Get("Frame").Get("PkgName").String()
		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Println("MonkeyFrame error in json file! AppPath, Argu and PkgName are needed!")
			return
		}
		mf.AppPath = appPath
		mf.Argu = argu
		mf.PkgName = pkg
		sj.Frame = mf
	case framework.FRAME_INSTALL:
		sj.FrameKind = framework.FRAME_INSTALL
		var inf framework.InstallFrame
		appPath, err1 := js.Get("Frame").Get("AppPath").String()
		pkg, err2 := js.Get("Frame").Get("PkgName").String()
		if err1 != nil || err2 != nil {
			fmt.Println("InstallFrame error in json file! File path of App and package are needed!")
			return
		}
		inf.AppPath = appPath
		inf.PkgName = pkg
		sj.Frame = inf
	default:
		fmt.Println("Unknow Framework!")
		return
	}

	switch filterkind {
	case framework.FILTER_SPECIFYDEVICES:
		sj.FilterKind = framework.FILTER_SPECIFYDEVICES
		var filter framework.SpecifyDevFilter
		idList, err1 := js.Get("Filter").Get("IdList").StringArray()
		repable, err2 := js.Get("Filter").Get("Replaceable").Bool()
		if err1 != nil || err2 != nil {
			fmt.Println("SpecifyDevices filter error in json file! IdList and Replaceable are needed!")
			return
		}
		filter.IdList = idList
		filter.Replaceable = repable
		sj.Filter = filter

	case framework.FILTER_SPECIFYATTR:
		sj.FilterKind = framework.FILTER_SPECIFYATTR
		var filter framework.SpecifyAttrFilter
		qt, err1 := js.Get("Filter").Get("Quantity").Int()
		at, err2 := js.Get("Filter").Get("Attr").String()
		vl, err3 := js.Get("Filter").Get("Value").String()
		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Println("SpecifyAttr filter error in json file! Quantity, Attr and Value  are needed!")
			return
		}
		filter.Attr = at
		filter.Quantity = qt
		filter.Value = vl
		sj.Filter = filter

	case framework.FILTER_COMPATIBILITY:
		sj.FilterKind = framework.FILTER_COMPATIBILITY
		var filter framework.CompatibilityFilter
		qt, err1 := js.Get("Filter").Get("Quantity").Int()
		dt, err2 := js.Get("Filter").Get("Dominate").String()
		if err1 != nil || err2 != nil {
			fmt.Println("Compatibility filter error in json file! Quantity, Dominate are needed!")
			return
		}
		filter.Dominate = dt
		filter.Quantity = qt
		sj.Filter = filter

	default:
		fmt.Println("Unknow Filter!")
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

		var devMap map[string]andevice.DeviceInfo = make(map[string]andevice.DeviceInfo)
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
