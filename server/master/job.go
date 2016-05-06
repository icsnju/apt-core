package master

import (
	"apsaras/comm"
	"apsaras/comm/comp"
	"apsaras/server/models"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type JobManager struct {
	jobMap  map[string]models.Job
	jobid   int
	jobLock *sync.Mutex
	idLock  *sync.Mutex
}

var jobManager *JobManager = &JobManager{make(map[string]models.Job), 0, new(sync.Mutex), new(sync.Mutex)}

func (m *JobManager) idGenerator() string {
	var id int64 = 0
	m.idLock.Lock()
	id = time.Now().Unix()
	m.idLock.Unlock()
	return strconv.FormatInt(id, 10)
}

func (m *JobManager) addJob(job models.Job) {
	m.jobLock.Lock()
	_, ex := m.jobMap[job.JobId]
	if ex {
		log.Println("Error! Job id repetitive: ", job.JobId)
	} else {
		m.jobMap[job.JobId] = job
	}
	m.jobLock.Unlock()
}

func (m *JobManager) deleteJob(id string) {
	m.jobLock.Lock()
	_, ex := m.jobMap[id]
	if ex {
		delete(m.jobMap, id)
	}
	m.jobLock.Unlock()
}

//handle sub job
func (m *JobManager) createJob(subjob models.SubJob) models.Job {

	mid := m.idGenerator()

	var job models.Job
	job.JobId = mid
	job.JobInfo = subjob
	job.StartTime = time.Now()
	job.LatestTime = time.Now()

	//get current device list
	devices := slaveManager.getDevices()

	//get this job
	devList := job.JobInfo.Filter.GetDeviceSet(devices)
	var taskMap map[string]comp.Task = make(map[string]comp.Task)
	for _, devId := range devList {
		var t comp.Task
		t.JobId = job.JobId
		t.DeviceId = devId
		t.TargetId = devId //勿忘初心
		t.State = comp.TASK_WAIT
		taskMap[devId] = t
	}
	job.TaskMap = taskMap

	return job
}

func (m *JobManager) updateJobs(taskinfo map[string]comp.Task) {
	m.jobLock.Lock()
	for _, t := range taskinfo {
		jid := t.JobId
		did := t.DeviceId
		job, ex := m.jobMap[jid]
		if !ex {
			//log.Println("Job not exist! ", jid)
			continue
		}
		_, ex = job.TaskMap[did]
		if !ex {
			log.Println("Device not exist! ", did)
			continue
		}
		job.TaskMap[did] = t
		m.jobMap[jid] = job
	}
	m.jobLock.Unlock()
}

//find the oldest job
func (m *JobManager) findOldJob(id string) string {

	var maxJobId string = "-1"
	oldTime := time.Now()

	m.jobLock.Lock()
	//find a old job
	for ke, jo := range m.jobMap {
		ts, ex := jo.TaskMap[id]
		if ex && ts.State == comp.TASK_WAIT {
			if jo.StartTime.Before(oldTime) {
				maxJobId = ke
				oldTime = jo.StartTime
			}
		}
	}
	m.jobLock.Unlock()
	return maxJobId
}

//find the best job
func (m *JobManager) findBestJob(id string) (comp.RunTask, bool) {
	var maxPriority float64 = -1
	var count1 int = 0
	var maxJob1 string = "-1"

	var maxWaitTime float64 = -1
	var count2 int = 0
	var maxJob2 string = "-1"

	bestId := "-1"
	//find a good job
	m.jobLock.Lock()
	for ke, jo := range m.jobMap {
		ts, ex := jo.TaskMap[id]
		if ex && ts.State == comp.TASK_WAIT {
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
		bestId = "-1"
	} else {
		rand.Seed(time.Now().UnixNano())
		lottery := rand.Intn(count1 + count2)
		if lottery < count1 {
			bestId = maxJob1
		} else {
			bestId = maxJob2
		}
	}

	var rt comp.RunTask
	var exit = false
	if bestId != "-1" {
		m.updateJobTaskState(bestId, id, comp.TASK_RUN)
		rt = m.createRuntask(bestId, id)
		exit = true
	}
	m.jobLock.Unlock()
	return rt, exit
}

//For func findBestJob only
func (m *JobManager) createRuntask(jobid, id string) comp.RunTask {
	var rt comp.RunTask
	rt.Frame = m.jobMap[jobid].JobInfo.Frame
	rt.TaskInfo = m.jobMap[jobid].TaskMap[id]
	return rt
}

//For func findBestJob only
func (m *JobManager) updateJobTaskState(jobId, taskId string, state int) {
	job, ex := m.jobMap[jobId]
	if ex {
		task := job.TaskMap[taskId]
		task.State = state
		task.StartTime = time.Now()
		job.TaskMap[taskId] = task
		job.LatestTime = time.Now()
		m.jobMap[jobId] = job
	}
}

func (m *JobManager) killJob(id string) {
	m.jobLock.Lock()
	job, ex := m.jobMap[id]
	if ex {
		isFail := false
		for tid, task := range job.TaskMap {
			if task.State != comp.TASK_COMPLETE && task.State != comp.TASK_FAIL {
				task.State = comp.TASK_FAIL
				task.FinishTime = time.Now()
				job.TaskMap[tid] = task
				isFail = true
			}
		}
		job.FinishTime = time.Now()
		var update map[string]interface{}
		//update job in db
		if isFail {
			update = bson.M{"$set": bson.M{models.JOB_STATUS: -1}}
		} else {
			update = bson.M{"$set": bson.M{models.JOB_STATUS: 100}}
		}
		models.UpdateJobSketchInDB(id, update)
		models.UpdateJobInDB(id, job)
		delete(m.jobMap, id)
	}
	m.jobLock.Unlock()
}

//Start find finished job cyclically
func (m *JobManager) updateJobInDB() {
	log.Println("Start update job state.")
	for {
		//update job status
		var finishedJobs []string = make([]string, 0)
		m.jobLock.Lock()
		for jid, job := range m.jobMap {

			all := 0
			pro := 0
			for _, ts := range job.TaskMap {
				all = all + 2
				if ts.State == comp.TASK_COMPLETE || ts.State == comp.TASK_FAIL {
					pro = pro + 2
				} else if ts.State == comp.TASK_RUN {
					pro = pro + 1
				}
			}
			rate := 100
			if all != 0 {
				rate = pro * 100 / all

			}
			//log.Println(all, " : ", pro)

			if all == pro {
				finishedJobs = append(finishedJobs, jid)
				job.FinishTime = time.Now()
			}
			update := bson.M{"$set": bson.M{models.JOB_STATUS: rate}}
			models.UpdateJobSketchInDB(jid, update)
			models.UpdateJobInDB(jid, job)
		}
		m.jobLock.Unlock()

		//delete finished job
		for _, id := range finishedJobs {
			log.Println("Delete finished job: ", id)
			m.deleteJob(id)
		}

		time.Sleep(comm.HEARTTIME)
	}
}

func (m *JobManager) contain(id string) bool {
	var ex bool
	m.jobLock.Lock()
	_, ex = m.jobMap[id]
	m.jobLock.Unlock()
	return ex
}
