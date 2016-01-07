package task

import (
	"nata/framework"
	"time"
)

//Job
//job id, framework and device list, time
type Job struct {
	JobId      string
	JobInfo    SubJob
	TaskMap    map[string]Task
	StartTime  time.Time
	FinishTime time.Time
	LatestTime time.Time
	Finished   bool
}

type JobMap struct {
	Map map[string]Job
}

type SubJob struct {
	FrameKind  string
	Frame      framework.FrameStruct
	FilterKind string
	Filter     framework.FilterInterface
}

//Get the priority of this job
func (job Job) GetPriority() float64 {
	waitTime := job.GetWaitTime()
	var devTime float64 = 0
	for _, ts := range job.TaskMap {
		switch ts.State {
		case TASK_COMPLETE:
			fallthrough
		case TASK_FAIL:
			devTime += ts.FinishTime.Sub(ts.StartTime).Seconds()
		case TASK_RUN:
			devTime += time.Now().Sub(ts.StartTime).Seconds()
		}
	}
	if devTime != 0 {
		return waitTime / devTime
	}
	//no task of this job has run
	return -1
}

//Get wait time of this job
func (job Job) GetWaitTime() float64 {
	waitTime := time.Now().Sub(job.LatestTime).Seconds()
	return waitTime
}
